package mysql

import (
	"database/sql"
	"errors"
	"github.com/raul-franca/go-snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires) 
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	//and the number of arguments must be exactly the same as the number of
	//columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// Se a query não retornar rows, então row.Scan() vai devolver um
		// sql.ErrNoRows error. é usado errors.Is() para checar qual foi o
		//error especificamente, models.ErrNoRecord error
		// instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	// Se tudo deu certo então retorne o objeto Snippet.
	return s, nil

}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() 
             ORDER BY created DESC LIMIT 10`
	// Use the Query() method on the connection pool to execute our SQL statement.
	// This returns a sql.Rows resultset containing the result of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	//statement should come *after* you check for an error from the Query()
	//Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()
	// Initialize an empty slice to hold the models.Snippets objects.
	snippets := []*models.Snippet{}
	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the // resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the // new Snippet object that we created. Again, the arguments to row.Scan() // must be pointers to the place you want to copy the data into, and the // number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	//error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

// Funções e variáveis importantes:
//
// DB.Exec() is used for statements which don’t return rows (like INSERT and DELETE).
// DB.Query() is used for SELECT queries which return multiple rows.
// DB.QueryRow() is used for SELECT queries which return a single row.
//
//  sql.Result interface returned by DB.Exec(). This provides two methods:
//		LastInsertId()  is not supported by PostgreSQL.
//		RowsAffected() — which returns the number of rows affected by the statement.
// rows.Scan() Converte automaticamente o output do SQL DB para os tipos nativos do Go
