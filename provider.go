package authentic

import "github.com/gobuffalo/buffalo"

//Provider allows developers to tell Authentic where to check user details
type Provider interface {

	//FindByID receives a userID and should check if the user exist,
	//this will be used by the AuthenticateMW for protected routes requests.
	FindByID(userID interface{}) (interface{}, error)

	FindByUsername(username string) (User, error)

	//SetUserDetails Allows App using Authentic to load details of the user
	//on every request after we've determined the user exists.
	SetUserDetails(user interface{}, c buffalo.Context) error
}
