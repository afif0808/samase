package jsonwebtokenredisrepo

import (
	"log"
	jsonwebtokenrepo "samase/jsonwebtoken/repository"

	"github.com/gomodule/redigo/redis"
)

//IsJWTBlackListed return if a jwt is blacklisted
func IsJWTBlackListed(conn redis.Conn) jsonwebtokenrepo.IsJWTBlackListedFunc {
	return func(token string) (bool, error) {
		rep, err := conn.Do("GET", "BlackListed"+token)
		if err != nil {
			log.Println(err)
			return false, err
		}
		return rep != nil, nil
	}
}

//BlackListJWT black list jwt
func BlackListJWT(conn redis.Conn) jsonwebtokenrepo.BlackListJWTFunc {
	return func(token string) error {
		_, err := conn.Do("SET", "BlackListed"+token, "LoggedOut")
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
