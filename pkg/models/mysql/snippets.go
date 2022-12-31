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

	return s, nil

}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() 
             ORDER BY created DESC LIMIT 10`
	// Use o Query() method no pool de conexão para executar nossa instrução SQL.
	// This returns a sql.Rows resultset containing the result of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Adiamos rows.Close() para garantir que o conjunto de resultados sql.Rows seja
	// sempre fechado corretamente antes que o método Latest() retorne. isso adia
	//instrução deve vir *depois* de você verificar se há um erro no Query()
	//Caso contrário, se Query() retornar um erro, você entrará em pânico
	// tentando fechar um conjunto de resultados nulo.
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		// Cria um ponteiro para um novo Snippet struct zerado.
		s := &models.Snippet{}
		// Use rows.Scan() para copiar os valores de cada campo na linha para o
		//devem ser ponteiros para o local onde você deseja copiar os dados e o
		//número de argumentos deve ser exatamente igual ao número de
		// colunas retornadas por sua instrução.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// Quando o rows.Next() loop terminou, chamamos rows.Err() para recuperar qualquer
	// erro que foi encontrado durante a iteração. é importante
	// chame isso - não assuma que uma iteração bem-sucedida foi concluída
	// sobre o conjunto de resultados.
	if err = rows.Err(); err != nil {
		return nil, err
	}

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
