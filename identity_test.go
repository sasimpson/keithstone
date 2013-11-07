package identity_test

import (
    "os"
    // "fmt"
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/sasimpson/identity"
)

func TestCredentialsFromEnvironment (t *testing.T) {
    os.Setenv("TEST_USER", "test_user")
    os.Setenv("TEST_KEY", "test_key")
    os.Setenv("TEST_URL", "test_url")

    i := identity.Identity{}
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

func TestCredentialsFromStrings (t *testing.T) {
    i := identity.Identity{}
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

/*
const JSONAuthReply = `{
    "access": {
        "serviceCatalog": [
            {
                "endpoints": [
                   {
                        "publicURL": "https://ord.servers.api.rackspacecloud.com/v2/12345",
                        "region": "ORD",
                        "tenantId": "12345",
                        "versionId": "2",
                        "versionInfo": "https://ord.servers.api.rackspacecloud.com/v2",
                        "versionList": "https://ord.servers.api.rackspacecloud.com/"
                    },
                    {
                        "publicURL": "https://dfw.servers.api.rackspacecloud.com/v2/12345",
                        "region": "DFW",
                        "tenantId": "12345",
                        "versionId": "2",
                        "versionInfo": "https://dfw.servers.api.rackspacecloud.com/v2",
                        "versionList": "https://dfw.servers.api.rackspacecloud.com/"
                    }
                ],
                "name": "cloudServersOpenStack",
                "type": "compute"
            },
            {
                "endpoints": [
                    {
                        "publicURL": "https://ord.databases.api.rackspacecloud.com/v1.0/12345",
                        "region": "ORD",
                        "tenantId": "12345"
                    },
                    {
                        "publicURL": "https://dfw.databases.api.rackspacecloud.com/v1.0/12345",
                        "region": "DFW",
                        "tenantId": "12345"
                    }
                ],
                "name": "cloudDatabases",
                "type": "rax:database"
            },
            {
                "endpoints": [
                    {
                        "publicURL": "https://ord.loadbalancers.api.rackspacecloud.com/v1.0/12345",
                        "region": "ORD",
                        "tenantId": "645990"
                    },
                    {
                        "publicURL": "https://dfw.loadbalancers.api.rackspacecloud.com/v1.0/12345",
                        "region": "DFW",
                        "tenantId": "12345"
                    }
                ],
                "name": "cloudLoadBalancers",
                "type": "rax:load-balancer"
            },
            {
                "endpoints": [
                    {
                        "publicURL": "https://cdn1.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc ",
                        "region": "DFW",
                        "tenantId": "MossoCloudFS_aaaa-bbbb-cccc "
                    },
                    {
                        "publicURL": "https://cdn2.clouddrive.com/v1/MossoCloudFS_aaaa-bbbb-cccc ",
                        "region": "ORD",
                        "tenantId": "MossoCloudFS_aaaa-bbbb-cccc "
                    }
                ],
                "name": "cloudFilesCDN",
                "type": "rax:object-cdn"
            },
            {
                "endpoints": [
                    {
                        "publicURL": "https://dns.api.rackspacecloud.com/v1.0/12345",
                        "tenantId": "12345"
                    }
                ],
                "name": "cloudDNS",
                "type": "rax:dns"
            },
            {
                "endpoints": [
                    {
                        "publicURL": "https://servers.api.rackspacecloud.com/v1.0/12345",
                        "tenantId": "12345",
                        "versionId": "1.0",
                        "versionInfo": "https://servers.api.rackspacecloud.com/v1.0",
                        "versionList": "https://servers.api.rackspacecloud.com/"
                    }
                ],
                "name": "cloudServers",
                "type": "compute"
            },
            {
                "endpoints": [
                    {
                        "publicURL": "https://monitoring.api.rackspacecloud.com/v1.0/12345",
                        "tenantId": "12345"
                    }
                ],
                "name": "cloudMonitoring",
                "type": "rax:monitor"
            },
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
*/

func TestAuthenticate(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        const serverStatus = 404
        const serverResponse = "Not Found"
        w.Header().Set("X-Foo", "w00t")
        w.WriteHeader(serverStatus)
        w.Write([]byte(serverResponse))
    }))
    defer server.Close()
    i := identity.Identity{}
    i.CredentialsFromStrings("test_user", "test_key", server.URL)
    err := i.Authenticate()
    t.Error(err)
}

/*
func ServeAuthRequest(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(*r.URL)
        h.ServeHTTP(w, r)
    })
}
*/

/*
func TestAuthenticate (t *testing.T) {
    i := identity.Identity{}
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello world\n") }))
    defer ts.Close()
    i.CredentialsFromStrings("test_user", "test_key", ts.URL)
    err := i.Authenticate()
    t.Error(err)
}
*/
