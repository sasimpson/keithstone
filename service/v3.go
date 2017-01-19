package service

import (
	"encoding/json"
	"net/http"
	"time"
)

type v3Token struct {
	Methods   []string  `json:"methods"`
	ExpiresAt time.Time `json:"expires_at"`
	Extras    struct {
	} `json:"extras"`
	User struct {
		Domain struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"domain"`
		ID                string      `json:"id"`
		Name              string      `json:"name"`
		PasswordExpiresAt interface{} `json:"password_expires_at"`
	} `json:"user"`
	Projects []v3Project `json:"projects"`
	AuditIds []string    `json:"audit_ids"`
	IssuedAt time.Time   `json:"issued_at"`
}

type v3Project struct {
	IsDomain    bool        `json:"is_domain"`
	Description interface{} `json:"description"`
	DomainID    string      `json:"domain_id"`
	Enabled     bool        `json:"enabled"`
	ID          string      `json:"id"`
	Links       struct {
		Self string `json:"self"`
	} `json:"links"`
	Name     string      `json:"name"`
	ParentID interface{} `json:"parent_id"`
}

type v3TokenResponse struct {
	Token v3Token `json:"token"`
}

type v3AuthRequest struct {
	Auth struct {
		Identity struct {
			Methods  []string `json:"methods"`
			Password struct {
				User struct {
					Name   string `json:"name"`
					Domain struct {
						Name string `json:"name"`
					} `json:"domain"`
					Password string `json:"password"`
				} `json:"user"`
			} `json:"password"`
		} `json:"identity"`
	} `json:"auth"`
}

type v3AuthResponse struct {
	Token v3Token `json:"token"`
}

func v3ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {

	adminTokenID := r.Header.Get("X-Auth-Token")
	userTokenID := r.Header.Get("X-Subject-Token")

	//validate admin token
	err := v3ValidateAdminToken(adminTokenID)
	if err != nil {
		if authError, ok := err.(*AuthError); ok {
			http.Error(w, authError.Error(), authError.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//validate user token
	err = v3ValidateToken(userTokenID)
	if err != nil {
		if authError, ok := err.(*AuthError); ok {
			http.Error(w, authError.Error(), authError.Status)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp v3TokenResponse
	var project v3Project

	project.ID = getProjectID(userTokenID)
	project.Enabled = true

	resp.Token.Methods = []string{"token"}
	resp.Token.ExpiresAt = time.Now().AddDate(0, 0, 1)
	resp.Token.IssuedAt = time.Now()
	resp.Token.Projects = []v3Project{project}
	resp.Token.User.Domain.ID = "default"
	resp.Token.User.Domain.Name = "Default"
	resp.Token.User.ID = getUserID(userTokenID)

	js, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, "json marshal error", http.StatusInternalServerError)
	}
	w.Header().Set("X-Subject-Token", userTokenID)
	w.Write([]byte(js))
	return
}

func v3GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// var authRequest v3AuthRequest
	// err := decoder.Decode(&authRequest)
	// if err != nil {
	// 	http.Error(w, "json decode failure", http.StatusBadRequest)
	// 	return
	// }
	// adminTokenID := r.Header.Get("X-Auth-Token")
	// userTokenID := r.Header.Get("X-Subject-Token")

	// // _, authError := v3GetToken(authRequest, adminTokenID, userTokenID)

	// if authError.Status/100 != 2 {
	// 	http.Error(w, authError.Error(), authError.Status)
	// 	return
	// }

}

func v3ValidateToken(token string) error {
	if validateToken(token) {
		return nil
	}
	return &AuthError{msg: "invalid", Status: http.StatusUnauthorized}
}

func v3ValidateAdminToken(token string) error {
	if validateAdminToken(token) {
		return nil
	}
	return &AuthError{msg: "invalid admin login", Status: http.StatusUnauthorized}
}
