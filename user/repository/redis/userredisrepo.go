package userredisrepo

import (
	"context"
	"log"
	userrepo "samase/user/repository"

	"github.com/gomodule/redigo/redis"
)

func SaveEmailConfirmationCode(conn redis.Conn) userrepo.SaveEmailConfirmationCodeFunc {
	return func(ctx context.Context, code string, duration int) error {
		_, err := conn.Do("SET", code, true)
		if err != nil {
			log.Println(err)
		}
		_, err = conn.Do("EXPIRE", code, duration)
		if err != nil {
			log.Println("here!", err)
		}
		return err
	}
}

func CheckEmailConfirmationCode(conn redis.Conn) userrepo.CheckEmailConfirmationCodeFunc {
	return func(ctx context.Context, code string) (bool, error) {
		rep, err := conn.Do("GET", code)
		if err != nil {
			log.Println(err)
			return false, err
		}
		return rep != nil, err
	}
}

func SavePasswordRecoveryCode(conn redis.Conn) userrepo.SavePasswordRecoveryCodeFunc {
	return func(ctx context.Context, code string, duration int) error {
		_, err := conn.Do("SET", code, true)
		if err != nil {
			log.Println(err)
		}
		_, err = conn.Do("EXPIRE", code, duration)
		if err != nil {
			log.Println("here!", err)
		}
		return err
	}
}

func CheckPasswordRecoveryCode(conn redis.Conn) userrepo.CheckPasswordRecoveryCodeFunc {
	return func(ctx context.Context, code string) (bool, error) {
		rep, err := conn.Do("GET", code)
		if err != nil {
			log.Println(err)
			return false, err
		}
		return rep != nil, err
	}
}
func RemovePasswordRecoveryCode(conn redis.Conn) userrepo.RemovePasswordRecoveryCodeFunc {
	return func(ctx context.Context, code string) error {
		_, err := conn.Do("DEL", code)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
