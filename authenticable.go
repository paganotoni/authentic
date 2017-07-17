package authentic

//Authenticable is the interface models would need to meet to be used for auth.
type Authenticable interface {
	GetEncryptedPassword() string
	GetID() interface{}
}
