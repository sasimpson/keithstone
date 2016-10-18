package service

import (
	"testing"
	"time"

	"database/sql"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var testUser User
var testToken Token

func init() {
	testUser = User{
		ID:   1,
		Name: "foo"}
	testToken = Token{
		ID:       uuid.NewV4().String(),
		IssuedAt: time.Now(),
		Expires:  time.Now().AddDate(0, 0, 1)}
}

func TestGetUserTokenExistingToken(t *testing.T) {
	//setup
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	assert.Nil(err, "opening database connection should work")
	defer dbi.close()

	mock.ExpectQuery("select id from tokens").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("select token_id, created_at, expiration from tokens").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"token_id", "created_at", "expiration"}).AddRow(testToken.ID, testToken.IssuedAt, testToken.Expires))

	token, err := dbi.getUserToken(testUser)
	assert.Equal(testToken.ID, token.ID)
	err = mock.ExpectationsWereMet()
	assert.Empty(err)
}

func TestGetUserTokenNoTokens(t *testing.T) {
	//setup
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	assert.Nil(err, "opening database connection should work")
	defer dbi.close()

	mock.ExpectQuery("select id from tokens").WithArgs(1).WillReturnError(sql.ErrNoRows)
	mock.ExpectPrepare("insert into tokens").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = dbi.getUserToken(testUser)
	err = mock.ExpectationsWereMet()
	assert.Empty(err)
}

func TestValidateLoginValid(t *testing.T) {
	//setup
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	assert.Nil(err, "opening database connection should work")
	defer dbi.database.Close()
	//test valid login for foo user
	mock.ExpectQuery("select id from users").
		WithArgs("foo", "bar").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("select id, name from users").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "foo"))

	var ar AuthRequest
	ar.Auth.PasswordCredentials.Username = "foo"
	ar.Auth.PasswordCredentials.Password = "bar"
	user, err := dbi.validateLogin(ar)
	err = mock.ExpectationsWereMet()
	assert.Empty(err)
	assert.Equal(user.Name, ar.Auth.PasswordCredentials.Username)
	assert.Equal(user.ID, 1)
}

func TestValidateLoginNotValid(t *testing.T) {
	//setup
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	assert.Nil(err, "opening database connection should work")
	defer dbi.database.Close()
	mock.ExpectQuery("select id from users").
		WithArgs("foo", "bar").
		WillReturnError(sql.ErrNoRows)
	var ar AuthRequest
	ar.Auth.PasswordCredentials.Username = "foo"
	ar.Auth.PasswordCredentials.Password = "bar"
	_, err = dbi.validateLogin(ar)
	assert.Error(err, "sql: no rows in result set")
	err = mock.ExpectationsWereMet()
	assert.Empty(err)
}

func TestNewToken(t *testing.T) {
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	assert.Nil(err, "opening database connection should work")
	defer dbi.close()

	// test getting a new token generated and inserted into the db
	mock.ExpectPrepare("insert into tokens").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))
	token, err := dbi.newToken(testUser)
	assert.Nil(err)
	assert.NotEmpty(token)
	err = mock.ExpectationsWereMet()
	assert.Empty(err)
}

/*
func TestNewAPIKey(t *testing.T) {
	t.Skip()
	dba, mock, err := sqlmock.New()
	assert := assert.New(t)
	dbi.database = dba
	assert.Nil(err, "opening database connection should work")
	defer dbi.close()
	mock.ExpectQuery("select count(id) from apikeys where user_id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).AddRow(1))
	_, err = dbi.newAPIKey(1)
	err = mock.ExpectationsWereMet()
}
*/
