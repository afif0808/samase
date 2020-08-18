package jsonwebtokenrepo

type IsJWTBlackListedFunc func(token string) (bool, error)
type BlackListJWTFunc func(token string) error
