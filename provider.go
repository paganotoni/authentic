package authentic

//Provider allows developers to tell Authentic where to check user details
type Provider interface {

	//UserExists receives a userID and should check if the user exist,
	//this will be used by the AuthenticateMW for protected routes requests.
	UserExists(userID interface{}) (bool, error)

	//AuthenticateUser will be called to check if the user/password combination
	// is correct, this should return the userID if valid, otherwise should return
	// nil and an error.
	AuthenticateUser(username, password string) (userID interface{}, err error)
}
