package keithstone_test

import (
	"github.com/sasimpson/keithstone"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func MockServer(status int, response string) (server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(response))
	}))
	return
}

func TestCredentialsFromEnvironment(t *testing.T) {
	os.Setenv("TEST_USER", "test_user")
	os.Setenv("TEST_KEY", "test_key")
	os.Setenv("TEST_URL", "test_url")

	i := keithstone.KeithStone{}
	i.CredentialsFromEnvironment("TEST_USER", "TEST_KEY", "TEST_URL")

	if i.User != "test_user" {
		t.Fatal("user env getter failed")
	}
	if i.Key != "test_key" {
		t.Fatal("key env getter failed")
	}
	if i.AuthURL != "test_url" {
		t.Fatal("url env getter failed")
	}
}

func TestCredentialsFromStrings(t *testing.T) {
	i := keithstone.KeithStone{}
	i.CredentialsFromStrings("test_user", "test_key", "test_url")

	if i.User != "test_user" {
		t.Fatal("user string failed")
	}
	if i.Key != "test_key" {
		t.Fatal("key string failed")
	}
	if i.AuthURL != "test_url" {
		t.Fatal("url string failed")
	}
}

func TestAuthenticateValid(t *testing.T) {
	server := MockServer(200, JSON200AuthReply)
	defer server.Close()
	i := keithstone.KeithStone{}
	i.CredentialsFromStrings("test_user", "test_key", server.URL)
	err := i.Authenticate()
	if err != nil {
		t.Error(err)
		t.Fatal("error returned from Authenicate method")
	}
}

func TestAuthenticateNotFound(t *testing.T) {
	server := MockServer(404, "Not Found")
	defer server.Close()
	i := keithstone.KeithStone{}
	i.CredentialsFromStrings("test_user", "test_key", server.URL)
	if err := i.Authenticate(); err == nil {
		t.Fatal("should have received error")
	}
}

func TestGetService(t *testing.T) {
    server := MockServer(200, JSON200AuthReply)
    defer server.Close()
    i := keithstone.KeithStone{}
    i.CredentialsFromStrings("test_user", "test_key", server.URL)
    i.Authenticate()
    if endpoints := i.GetService("object-store", "DFW"); endpoints.PublicURL != "https://storage101.dfw1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc " {
        t.Fatal("endpoints don't match for GetService DFW")
    }
    if endpoints := i.GetService("object-store", "ORD"); endpoints.PublicURL != "https://storage101.ord1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc " {
        t.Fatal("endpoints don't match for GetService ORD")
    }
}

//{"itemNotFound":{"code":404,"message":"Resource Not Found"}}
//{"badRequest":{"code":400,"message":"JSON Parsing error"}}

const JSON200AuthReply = `{
    "access": {
        "serviceCatalog": [
            {
                "endpoints": [
                    {
                        "internalURL": "https://snet-storage101.dfw1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc ",
                        "publicURL": "https://storage101.dfw1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc ",
                        "region": "DFW",
                        "tenantId": "MossoCloudFS_aaaa-bbbb-cccc"
                    },
                    {
                        "internalURL": "https://snet-storage101.ord1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc ",
                        "publicURL": "https://storage101.ord1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc ",
                        "region": "ORD",
                        "tenantId": "MossoCloudFS_aaaa-bbbb-cccc"
                    }
                ],
                "name": "cloudFiles",
                "type": "object-store"
            }
        ],
        "token": {
            "expires": "2012-04-13T13:15:00.000-05:00",
            "id": "aaaaa-bbbbb-ccccc-dddd"
        },
        "user": {
        "RAX-AUTH:defaultRegion": "DFW",
            "id": "161418",
            "name": "demoauthor",
            "roles": [
                {
                    "description": "User Admin Role.",
                    "id": "3",
                    "name": "identity:user-admin"
                }
            ]
        }
    }
}`
