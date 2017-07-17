package authentic

import "github.com/gobuffalo/buffalo"

//Provider allows developers to tell Authentic where to check user details
type Provider interface {

	//FindByID receives a userID and should check if the user exist,
	//this will be used by the AuthenticateMW for protected routes requests.
	FindByID(userID interface{}) (Authenticable, error)

	//FindByUsername will be called when authorizing a username by username/password
	//this allows applications to relate the username password with the form to the
	//desired field to look for the user.
	FindByUsername(username string) (Authenticable, error)

	//UserDetails Allows App using Authentic to load details of the user
	//on every request after we've determined the user exists, this function
	//will be called on every request by Authentic.
	UserDetails(user Authenticable, c buffalo.Context) error
}
