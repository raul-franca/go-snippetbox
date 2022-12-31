package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/raul-franca/go-snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

// Insert method to add a usuario
func (m *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created, active) 
		VALUES(?, ?, ?, UTC_TIMESTAMP, TRUE)`
	_, err = m.DB.Exec(stmt, name, email, string(hash))
	if err != nil {
		// Se retornar um erro, usamos a função errors.As() para verificar
		//se o erro é do tipo *mysql.MySQLError. Se o fizer, o
		// o erro será atribuído à variável mySQLError. Podemos então verificar
		//se o erro está ou não relacionado à nossa chave users_uc_email por
		// verificando o conteúdo da string da mensagem. Se isso acontecer, nós retornamos
		//um erro ErrDuplicateEmail.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Authenticate method para verificar se existe um usuário com
// e-mail e a senha fornecidos. retornara o user ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// Recupera o id e a senha com hash associados ao e-mail fornecido. Se não
	// existe um e-mail correspondente ou o usuário não está ativo, retornamos o
	// Erro ErrInvalidCredentials.
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// Get method para buscar detalhes de um usuário específico com base em seu ID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
