package service

import "strconv"

// AuthError - Used to pass auth errors that include the http status code
type AuthError struct {
	msg    string
	Status int
}

func (e *AuthError) Error() string { return e.msg }

func validateToken(token string) bool {
	if len(token) < 5 {
		return false
	}
	settings, err := strconv.ParseInt(token[len(token)-MASKLEN:], 16, 64)
	if err != nil {
		return false
	}
	return settings&ISVALID != 0
}

func validateAdminToken(token string) bool {
	if validateToken(token) {
		settings, err := strconv.ParseInt(token[len(token)-MASKLEN:], 16, 64)
		if err != nil {
			return false
		}
		return settings&ISADMIN != 0
	}
	return false
}

func getProjectID(token string) string {
	//assume token is long enough because it has been validated already
	return token[0:PROJECTIDLEN]
}

func getUserID(token string) string {
	return token[PROJECTIDLEN : PROJECTIDLEN+USERIDLEN]
}
