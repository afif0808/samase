package authenticationservice

type LogoutFunc func(token string) error
type IsLoggedOutFunc func(token string) (bool, error)
