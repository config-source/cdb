package auth

type User struct {
	ID int `db:"id"`

	Email string `db:"email"`
	// Completely ignored in JSON to avoid any accidental leakages. Incoming
	// requests have a dedicated write-only type that's stored in the api
	// package.
	Password string `db:"password" json:"-"`

	Roles []string `json:"roles"`
}

type TokenSet struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
}

type AuthenticationGateway interface {
	Login(username, password string) (User, TokenSet, error)
}

type AuthorizationGateway interface {
	// TODO: should this just return an error or just a bool?
	HasPermission(roles []string, permission string) (bool, error)
}
