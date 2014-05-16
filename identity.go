package identity

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

type Identity struct {
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
        Token struct {
            Expires time.Time
            ID      string
            Tenant struct {
                ID string
                Name string
            }
        }
        User struct {
            ID            string
            Name          string
            DefaultRegion string
            Roles         []struct {
                Description string
                ID          string
                Name        string
                TenantID    string
            }
        }
        ServiceCatalog []struct {
            Name      string
            Type      string
            Endpoints []Endpoint
        }
    }
}

type Endpoint struct {
    InternalURL string
    PublicURL   string
    Region      string
    TenantID    string
}

func (i *Identity) Authenticate() error {
    auth := v2AuthRequest{}
    auth.Auth.ApiKeyCredentials.UserName = i.User
    auth.Auth.ApiKeyCredentials.ApiKey = i.Key
    v2auth, error := json.Marshal(auth)
    if error != nil {
        return error
    }
    fmt.Printf("%s\n\n", v2auth)
    client := &http.Client{}
    url := []string{i.AuthURL, "/v2.0/tokens"}
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
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Printf("%s", body)
        return fmt.Errorf("%d", status)
    }
    defer resp.Body.Close()
    body, error := ioutil.ReadAll(resp.Body)
    authResponse := v2AuthResponse{}
    error = json.Unmarshal(body, &authResponse)
    if error != nil {
        return error
    }
    i.AuthData = authResponse
    return nil
}

func (i *Identity) GetService(service string, region string) Endpoint {
    if region == "" {
        region = i.AuthData.Access.User.DefaultRegion
    }
    serviceCatalog := i.AuthData.Access.ServiceCatalog
    for _, sc := range serviceCatalog {
        if sc.Type == strings.ToLower(service) {
            for _, ep := range sc.Endpoints {
                if ep.Region == strings.ToUpper(region) {
                    return ep
                }
            }
        }
    }
    return Endpoint{}
}

func (i *Identity) GetToken() string {
    return i.AuthData.Access.Token.ID
}

func (i *Identity) GetExpires() time.Time {
    return i.AuthData.Access.Token.Expires
}

func (i *Identity) CredentialsFromEnvironment(user_env string, key_env string, url_env string) {
    i.User = os.Getenv(user_env)
    i.Key = os.Getenv(key_env)
    i.AuthURL = os.Getenv(url_env)
}

func (i *Identity) CredentialsFromStrings(user string, key string, url string) {
    i.User = user
    i.Key = key
    i.AuthURL = url
}
