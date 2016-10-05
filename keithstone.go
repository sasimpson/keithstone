package keithstone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type KeithStone struct {
	AuthData v2AuthResponse
	User     string
	Key      string
	AuthURL  string
	Token    string
	Expires  time.Time
}

/* need to figure out the marshalling for the struct pointer, this works
   right now for api requests only
*/
type v2AuthRequest struct {
	Auth struct {
		ApiKeyCredentials struct {
			UserName string `json:"username,omitempty"`
			ApiKey   string `json:"apiKey,omitempty"`
		} `json:"RAX-KSKEY:apiKeyCredentials,omitempty"`
		PasswordCredentials *struct {
			UserName string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"passwordCredentials,omitempty"`
		Tenant   string `json:"tenantName,omitempty"`
		TenantId string `json:"tenantId,omitempty"`
	} `json:"auth"`
}

type v2AuthResponse struct {
	Access struct {
		ServiceCatalog []struct {
			Endpoints []Endpoint `json:"endpoints"`
			Name      string     `json:"name"`
			Type      string     `json:"type"`
		} `json:"serviceCatalog"`
		Token struct {
			Expires time.Time `json:"expires"`
			ID      string    `json:"id"`
			Tenant  struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
		} `json:"token"`
		User struct {
			DefaultRegion string `json:"RAX-AUTH:defaultRegion"`
			ID            string `json:"id"`
			Name          string `json:"name"`
			Roles         []struct {
				Description string `json:"description"`
				ID          string `json:"id"`
				Name        string `json:"name"`
			} `json:"roles"`
		} `json:"user"`
	} `json:"access"`
}

type Endpoint struct {
	InternalURL string
	PublicURL   string
	Region      string
	TenantID    string
}

func (ks *KeithStone) Authenticate() error {
	auth := v2AuthRequest{}
	auth.Auth.ApiKeyCredentials.UserName = ks.User
	auth.Auth.ApiKeyCredentials.ApiKey = ks.Key
	v2auth, error := json.Marshal(auth)
	if error != nil {
		return error
	}
	// fmt.Printf("auth request: %s\n", v2auth)
	client := &http.Client{}
	url := []string{ks.AuthURL, "/v2.0/tokens"}
	req, error := http.NewRequest("GET", strings.Join(url, ""), bytes.NewBuffer(v2auth))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	if error != nil {
		return error
	}
	resp, error := client.Do(req)
	if error != nil {
		return error
	}
	status := resp.StatusCode
	if status/100 != 2 {
		/*  our identity is stupid.  it returns stupid json with the main key
		being dependent on the error instead of a common key like 'code' or
		'error' to pull from.  getting the "message" will mean pulling all
		those apart.  until that is done, this is just going to return the
		status code from the server.
		*/
		/*
			   body, error := ioutil.ReadAll(resp.Body)
			   var f interface{}
			   error = json.Unmarshal(body, &f)
			   // fmt.Printf("%s", body)
			   if error != nil {
				   return error
			   }
			   m := f.(map[string]interface{})
			   for k,v := range m {
				   fmt.Printf("%s=>%s\n", k,v)
				   switch k {
				   case "itemNotFound":
					   return fmt.Errorf("%d %s", m[k].code, m[k].message)
				   }
			   }*/
		// body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println("Auth Body: ", body)
		return fmt.Errorf("%d", status)
	}
	defer resp.Body.Close()
	body, error := ioutil.ReadAll(resp.Body)
	// fmt.Printf("auth body: %s\n", body)
	authResponse := v2AuthResponse{}
	error = json.Unmarshal(body, &authResponse)
	if error != nil {
		return error
	}
	ks.AuthData = authResponse
	return nil
}

func (ks *KeithStone) GetService(service string, region string) Endpoint {
	if region == "" {
		region = ks.AuthData.Access.User.DefaultRegion
	}
	serviceCatalog := ks.AuthData.Access.ServiceCatalog
	// fmt.Println("serviceCatalog: ", serviceCatalog)
	// fmt.Println("access: ", ks.AuthData.Access)
	for _, sc := range serviceCatalog {
		// fmt.Println("sc: ", sc)
		if sc.Type == strings.ToLower(service) {
			for _, ep := range sc.Endpoints {
				// fmt.Printf("ep: %s ", ep)
				if ep.Region == strings.ToUpper(region) {
					return ep
				}
			}
		}
	}
	// fmt.Println("didn't find endpoint")
	return Endpoint{}
}

func (ks *KeithStone) GetToken() string {
	return ks.AuthData.Access.Token.ID
}

func (ks *KeithStone) GetExpires() time.Time {
	return ks.AuthData.Access.Token.Expires
}

func (ks *KeithStone) CredentialsFromEnvironment(user_env string, key_env string, url_env string) {
	ks.User = os.Getenv(user_env)
	ks.Key = os.Getenv(key_env)
	ks.AuthURL = os.Getenv(url_env)
}

func (ks *KeithStone) CredentialsFromStrings(user string, key string, url string) {
	ks.User = user
	ks.Key = key
	ks.AuthURL = url
}
