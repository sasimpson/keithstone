package service

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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
	user, err := dbi.validateLogin(ar)
	assert.Equal(user.Name, ar.Auth.PasswordCredentials.Username)
	assert.Equal(user.ID, 1)
	mock.ExpectQuery("select id, name from users").
		WithArgs("foo", "bar").
		WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	user, err = dbi.validateLogin(ar)
	if err != nil {
		t.Logf("%v", err)
	}
}

func TestNewToken(t *testing.T) {
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbi.close()
	testToken := Token{
		ID:       uuid.NewV4().String(),
		IssuedAt: time.Now(),
		Expires:  time.Now().AddDate(0, 0, 1)}
	mock.ExpectPrepare("insert into tokens (token_id, user_id, created_at, expiration)").
		ExpectExec().
		WithArgs(testToken.ID, 1, testToken.IssuedAt, testToken.Expires)
	token := dbi.newToken(1)
	assert.Equal(Token{}.ID, token.ID)
}
