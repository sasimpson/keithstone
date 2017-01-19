# Keith Stone
![keith stone](http://payload.cargocollective.com/1/3/101945/1323574/Daily_OOH.jpg)

Tired of the insane list of dependencies you need to install in order to test your openstack stuff?  try Keith Stone!  He's always smooth.  Only one binary to make your Keystone testing as smooth as Keith Stone.

## use

tokens look like this:

    project id         config bitmask
    |                  |
    vvvvvvvv        vvvv
    abcdefghstuvwxyz0001
            ^^^^^^^^
                   |
                userid

the *abcdefgh* part is the project id you wish to have passed back and the same for the *stuvwxyz* component for the user id

the config bitmask is like so:

| config bitmask | configuration |
| --- | --- |
|0000|invalid token|
|0001| valid token, not as admin|
|0002| invalid token, as admin |
|0003| valid token as admin |


## validating a token

v3 token validation requires and admin token and a user token.  they are sent as headers, admin as X-Auth-Token and user as X-Subject-Token.

valid token, valid auth user and token (i like using httpie):
    
    http -j :8080/v3/tokens x-auth-token:adminprjadminusr0003 x-subject-token:projidxxuseridxx0001

returns:

    HTTP/1.1 200 OK
    Content-Length: 410
    Content-Type: text/plain; charset=utf-8
    Date: Thu, 19 Jan 2017 03:42:32 GMT
    X-Subject-Token: projidxxuseridxx0001

    {
        "token": {
            "audit_ids": null,
            "expires_at": "2017-01-19T21:42:32.051936868-06:00",
            "extras": {},
            "issued_at": "2017-01-18T21:42:32.05193867-06:00",
            "methods": [
                "token"
            ],
            "projects": [
                {
                    "description": null,
                    "domain_id": "",
                    "enabled": true,
                    "id": "projidxx",
                    "is_domain": false,
                    "links": {
                        "self": ""
                    },
                    "name": "",
                    "parent_id": null
                }
            ],
            "user": {
                "domain": {
                    "id": "default",
                    "name": "Default"
                },
                "id": "useridxx",
                "name": "",
                "password_expires_at": null
            }
        }
    }

setting the user token bitmask to 0000 will return a 401, as will setting the admin bitmask to 0002