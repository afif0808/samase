package userredisrepo

import (
	"context"
	"fmt"
	"log"
	userrepo "samase/user/repository"
	"strconv"

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

func GetUserIDByCode(conn redis.Conn) userrepo.GetUserIDByCodeFunc {
	return func(ctx context.Context, code string) (int64, error) {
		resp, err := conn.Do("GET", code)
		if err != nil {
			log.Println(err)
		}
		id, err := strconv.ParseInt((fmt.Sprint(resp)), 10, 64)
		if err != nil {
			return 0, err
		}
		return id, err
	}
}
func SaveUserIDByCode(conn redis.Conn) userrepo.SaveUserIDByCodeFunc {
	return func(ctx context.Context, code string, id int64) error {
		_, err := conn.Do("SET", code, id)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}

func RetrieveUserIDByCode(conn redis.Conn) userrepo.RetrieveUserIDByCodeFunc {
	return func(ctx context.Context, code string) (int64, error) {
		resp, err := conn.Do("GET", code)
		if err != nil {
			log.Println(err)
			return 0, nil
		}
		id, err := strconv.ParseInt(fmt.Sprintf("%s", resp), 10, 64)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		_, err = conn.Do("DEL", code)

		if err != nil {
			log.Println(err)
		}
		return id, err
	}
}
