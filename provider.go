package authentic

import "github.com/gobuffalo/buffalo"

//Provider allows developers to tell Authentic where to check user details
type Provider interface {

	//FindUser receives a userID and should check if the user exist,
	//this will be used by the AuthenticateMW for protected routes requests.
	FindUser(userID interface{}) (interface{}, error)

	//AuthenticateUser will be called to check if the user/password combination
	// is correct, this should return the userID if valid, otherwise should return
	// nil and an error.
	AuthenticateUser(username, password string) (userID interface{}, err error)

	//SetUserDetails Allows App using Authentic to load details of the user
	//on every request after we've determined the user exists.
	SetUserDetails(user interface{}, c buffalo.Context) error
}
