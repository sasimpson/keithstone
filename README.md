# keithstone
openstack identity mocking and bindings for Go

tired of the insane list of dependencies you need to install in order to test your openstack stuff?  try keith stone!  he's smooth.  only one binary to test v2 and soon v3 identity.

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
