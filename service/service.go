package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
	"github.com/satori/go.uuid"
)

var dbi db

type db struct {
	database *sql.DB
}

//Token - the structure of a token response
type Token struct {
	IssuedAt time.Time `json:"issued_at"`
	Expires  time.Time `json:"expires"`
	ID       string    `json:"id"`
	Tenant   struct {
		Description interface{} `json:"description,omitempty"`
		Enabled     bool        `json:"enabled,omitempty"`
		ID          string      `json:"id"`
		Name        string      `json:"name"`
	} `json:"tenant"`
}

//Catalog - structure of the service catalog response
type Catalog struct {
	Endpoints []struct {
		AdminURL    string `json:"adminURL"`
		Region      string `json:"region"`
		InternalURL string `json:"internalURL"`
		ID          string `json:"id"`
		PublicURL   string `json:"publicURL"`
	} `json:"endpoints"`
	EndpointsLinks []interface{} `json:"endpoints_links"`
	Type           string        `json:"type"`
	Name           string        `json:"name"`
}

//User - structure of the user response
type User struct {
	RolesLinks []interface{} `json:"roles_links"`
	ID         int           `json:"id"`
	Roles      []Role        `json:"roles"`
	Name       string        `json:"name"`
}

//Role - role item
type Role struct {
	Name string `json:"name"`
}

//AuthRequest - incoming auth request from user.
type AuthRequest struct {
	Auth struct {
		PasswordCredentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"passwordCredentials"`
		TenantName string `json:"tenantName,omitempty"`
		Token      struct {
			ID string `json:"id,omitempty"`
		} `json:"token,omitempty"`
	} `json:"auth"`
}

//AuthResponse - the auth response structure
type AuthResponse struct {
	Access struct {
		Token    Token     `json:"token"`
		Catalog  []Catalog `json:"serviceCatalog,omitempty"`
		User     User      `json:"user"`
		Metadata struct {
			IsAdmin int      `json:"is_admin,omitempty"`
			Roles   []string `json:"roles,omitempty"`
		} `json:"metadata,omitempty"`
		Trust struct {
			ID            string `json:"id,omitempty"`
			TrusteeUserID string `json:"trustee_user_id,omitempty"`
			TrustorUserID string `json:"trustor_user_id,omitempty"`
			Impersonation bool   `json:"impersonation,omitempty"`
		} `json:"trust,omitempty"`
	} `json:"access"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/token", getTokenHandler)
	r.HandleFunc("/token/{tokenID}", validateTokenHandler)
	http.Handle("/", r)
}

//Serve - run service
func Serve() {
	http.ListenAndServe(":8080", nil)
}

//Service Handlers:
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "cool")
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("getToken")
	//pull apart incoming request
	decoder := json.NewDecoder(r.Body)
	var authRequest AuthRequest
	err := decoder.Decode(&authRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "json decode failure", http.StatusBadRequest)
		return
	}
	//validate user data
	dbi.connect()
	defer dbi.close()
	user, err := dbi.validateLogin(authRequest)
	if err != nil {
		http.Error(w, "invalid login", http.StatusUnauthorized)
	}
	//get the users token
	token, err := dbi.getUserToken(user.ID)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("%v", token)
	log.Printf("%v", authRequest)
	var authResponse AuthResponse
	authResponse.Access.Token = token
	authResponse.Access.User = user
	js, _ := json.Marshal(authResponse)
	w.Write([]byte(js))
}

func validateTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("validateToken")
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]
	log.Printf("tokenID: %v", tokenID)

	var validateResponse AuthResponse
	validateResponse.Access.User = User{Name: "userA", ID: 100}
	validateResponse.Access.Token = Token{
		IssuedAt: time.Now(),
		Expires:  time.Now().AddDate(0, 0, 1),
		ID:       "123456abcxzy"}
	validateResponse.Access.Token.Tenant.ID = "tenantA1"
	validateResponse.Access.Token.Tenant.Name = "tenantA1"
	validateResponse.Access.User.Roles = []Role{}

	js, _ := json.Marshal(validateResponse)
	w.Write([]byte(js))
}

//db interface methods
func (db *db) connect() {
	db.database, _ = sql.Open("sqlite3", "./auth.db")
}

func (db *db) close() {
	db.database.Close()
}

/* validateLogin - validates the user's login against the database
 */
func (db *db) validateLogin(ar AuthRequest) (User, error) {
	log.Println("validateLogin")
	username := ar.Auth.PasswordCredentials.Username
	password := ar.Auth.PasswordCredentials.Password
	var user User
	err := db.database.QueryRow("select id, name from users where name = ? and password = ?", username, password).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *db) validateToken(tokenID string) bool {
	var userID string
	err := db.database.QueryRow("SELCT user_id from tokens where token_id = ?", tokenID).Scan(&userID)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

/* getUserToken - get the token in the DB for the user, generate one if there isn't one
 */
func (db *db) getUserToken(userID int) (Token, error) {
	var token Token
	row := db.database.QueryRow("select token_id, expiration, created_at from tokens where user_id = ?", userID)
	log.Printf("%v", row)
	err := row.Scan(&token.ID, &token.Expires, &token.IssuedAt)
	if err == sql.ErrNoRows {
		token = db.newToken(userID)
	} else {
		if err != nil {
			return Token{}, err
		}
	}

	return token, nil
}

/* newToken: generates a new token and inserts it in the database for the user.
 */
func (db *db) newToken(userID int) Token {
	stmt, _ := db.database.Prepare("insert into tokens (token_id, user_id, created_at, expiration) values (?, ?, ?, ?)")
	defer stmt.Close()
	token := Token{
		ID:       uuid.NewV4().String(),
		IssuedAt: time.Now(),
		Expires:  time.Now().AddDate(0, 0, 1)}
	_, err := stmt.Exec(token.ID, userID, token.IssuedAt, token.Expires)
	if err != nil {
		return Token{}
	}
	return token
}

func (db *db) newAPIKey(userID int) (string, error) {
	key := uuid.NewV4()
	var apiKeyID int
	db.database.QueryRow("select id from apikeys where user_id = ?", userID).Scan(&apiKeyID)
	stmt, err := db.database.Prepare("delete from apikeys where user_id = ? and id = ?")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(userID, apiKeyID)
	if err != nil {
		return "", err
	}
	stmt, err = db.database.Prepare("insert into apikeys (key, user_id, created_by, created_at) values (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	res, err := stmt.Exec(key.String(), userID, userID, time.Now())
	log.Println(res.LastInsertId())
	return key.String(), nil
}
