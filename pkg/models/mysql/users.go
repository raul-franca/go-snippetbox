package mysql

import (
	"database/sql"
	"github.com/raul-franca/go-snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert method to add a usuario
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate method para verificar se existe um usuário com
// e-mail e a senha fornecidos. retornara o user ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get method para buscar detalhes de um usuário específico com base em seu ID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
