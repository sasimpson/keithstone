package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
	"github.com/satori/go.uuid"
)

var (
	dbi       db
	appConfig config
)

type db struct {
	database *sql.DB
}

type config struct {
	LazyAuth bool
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
		// TenantName string `json:"tenantName,omitempty"`
		// Token      struct {
		// 	ID string `json:"id,omitempty"`
		// } `json:"token,omitempty"`
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

	s := r.PathPrefix("/v2").Subrouter()
	s.HandleFunc("/token", getTokenHandler).Methods("POST")
	s.HandleFunc("/token/{tokenID}", validateTokenHandler).Methods("POST")
	http.Handle("/", r)
}

//Serve - run service
func Serve() {
	appConfig.LazyAuth = true
	http.ListenAndServe(":8080", nil)
}

//Service Handlers:
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "cool")
}

func getTokenHandler(w http.ResponseWriter, r *http.Request) {
	//pull apart incoming request
	decoder := json.NewDecoder(r.Body)
	var authRequest AuthRequest
	err := decoder.Decode(&authRequest)
	log.Printf("getTokenHandler authRequest: %v", authRequest)
	if err != nil {
		http.Error(w, "json decode failure", http.StatusBadRequest)
		return
	}
	//validate user data
	dbi.connect()
	defer dbi.close()
	user, err := dbi.validateLogin(authRequest)
	if err != nil {
		http.Error(w, "invalid login", http.StatusUnauthorized)
		return
	}
	//get the users token
	token, err := dbi.getUserToken(user)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("getTokenHandler token: %v", token)
	log.Printf("getTokenHandler user: %v", user)
	var authResponse AuthResponse
	authResponse.Access.Token = token
	authResponse.Access.User = user
	js, _ := json.Marshal(authResponse)
	w.Write([]byte(js))
}

func validateTokenHandler(w http.ResponseWriter, r *http.Request) {
	userTokenID := r.Header.Get("X-Auth-Token")
	if validateToken(userTokenID) {
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
		return
	}
	http.Error(w, "invalid login", http.StatusForbidden)
	return
}

func validateToken(token string) bool {
	if appConfig.LazyAuth == true && checkLazyAuth(token) {
		return true
	}
	//TODO add in token validation stack
	return false
}

func checkLazyAuth(token string) bool {
	if strings.HasSuffix(token, "a0") || strings.HasSuffix(token, "sm00th") {
		return true
	}
	return false
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
	username := ar.Auth.PasswordCredentials.Username
	password := ar.Auth.PasswordCredentials.Password
	log.Printf("validateLogin username: %v", username)
	log.Printf("validateLogin password: %v", password)

	var userID int
	err := db.database.QueryRow("select id from users where name = ? and password = ?", username, password).Scan(&userID)
	log.Printf("validateLogin userID: %v", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, err
		}
		log.Panic(err)
	}
	user, err := db.User(userID)
	log.Printf("validateLogin user: %v", user)
	if err != nil {
		log.Panic(err)
	}
	return user, nil
}

func (db *db) validateToken(tokenID string) bool {
	var userID string
	err := db.database.QueryRow("select user_id from tokens where token_id = ?", tokenID).Scan(&userID)
	if err != nil {
		return false
	}
	return true
}

/* getUserToken - get the token in the DB for the user, generate one if there isn't one
 */
func (db *db) getUserToken(user User) (Token, error) {
	var token Token
	tokens, err := db.UserTokens(user)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("getUserToken tokens: %d %v", len(tokens), tokens)
	if len(tokens) == 0 {
		log.Printf("getUserToken user: %v", user)
		token, err = db.newToken(user)
		log.Printf("getUserToken token: %v", token)
		if err != nil {
			log.Panic(err)
		}
	} else {
		token = tokens[0]
	}
	return token, nil
}

func (db *db) newAPIKey(userID int) (string, error) {
	key := uuid.NewV4()
	var keyCount int
	err := db.database.QueryRow("select count(id) from apikeys where user_id = ?", userID).Scan(&keyCount)
	if keyCount > 0 {
		var apiKeyID int
		err := db.database.QueryRow("select id from apikeys where user_id = ?", userID).Scan(&apiKeyID)
		stmt, err := db.database.Prepare("delete from apikeys where user_id = ? and id = ?")
		defer stmt.Close()
		if err != nil {
			return "", err
		}
		_, err = stmt.Exec(userID, apiKeyID)
		if err != nil {
			return "", err
		}
	}
	stmt, err := db.database.Prepare("insert into apikeys (key, user_id, created_by, created_at) values (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(key.String(), userID, userID, time.Now())
	return key.String(), nil
}

/* newToken: generates a new token and inserts it in the database for the user.
 */
func (db *db) newToken(user User) (Token, error) {
	token := Token{
		ID:       uuid.NewV4().String(),
		IssuedAt: time.Now(),
		Expires:  time.Now().AddDate(0, 0, 1)}
	db.SetToken(token, user)
	return token, nil
}

func (db *db) Token(id int) (Token, error) {
	var token Token
	err := db.database.QueryRow("select token_id, created_at, expiration from tokens where id = ?", id).Scan(&token.ID, &token.IssuedAt, &token.Expires)
	return token, err
}

func (db *db) SetToken(token Token, user User) error {
	stmt, err := db.database.Prepare("insert into tokens (token_id, user_id, created_at, expiration) values (?,?,?,?)")
	_, err = stmt.Exec(token.ID, user.ID, token.IssuedAt, token.Expires)
	stmt.Close()
	return err
}

func (db *db) UserForToken(token string) (User, error) {
	var userID int
	err := db.database.QueryRow("select user_id from tokens where token_id = ?", token).Scan(&userID)
	if err != nil {
		return User{}, err
	}
	user, err := db.User(userID)
	if err != nil {
		return User{}, err
	}
	log.Printf("db.UserForToken: %v", user.ID)
	return user, nil
}

func (db *db) User(id int) (User, error) {
	var user User
	err := db.database.QueryRow("select id, name from users where id = ?", id).Scan(&user.ID, &user.Name)
	log.Printf("db.User user: %v", user)
	return user, err
}

func (db *db) UserTokens(user User) ([]Token, error) {
	var tokens []Token
	rows, err := db.database.Query("select id from tokens where user_id = ? order by expiration", user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Token{}, nil
		}
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var tokenID int
		// log.Printf("db.UserTokens tokenID: %v", tokenID)
		if err := rows.Scan(&tokenID); err != nil {
			log.Fatal(err)
		}
		token, _ := db.Token(tokenID)
		tokens = append(tokens, token)
	}
	return tokens, nil
}
