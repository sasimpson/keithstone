package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

/*
type TestDB struct {
}

func (db *TestDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return sql.Row{}
}
func (db *TestDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New("error")
}
func (db *TestDB) Prepare(query string) (*sql.Stmt, error) {
	return nil, errors.New("error")
}
*/

func TestValidateLogin(t *testing.T) {
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbi.database.Close()
	mock.ExpectQuery("select id, name from users").
		WithArgs("foo", "bar").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "foo"))
	var ar AuthRequest
	ar.Auth.PasswordCredentials.Username = "foo"
	ar.Auth.PasswordCredentials.Password = "bar"
	testUser, err := dbi.validateLogin(ar)
	assert.Equal(testUser.Name, ar.Auth.PasswordCredentials.Username)
	assert.Equal(testUser.ID, 1)
}
