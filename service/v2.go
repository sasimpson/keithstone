package service

import (
	"net/http"
	"time"
)

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

func v2getTokenHandler(w http.ResponseWriter, r *http.Request) {

}

func v2validateTokenHandler(w http.ResponseWriter, r *http.Request) {

}
