package userredisrepo

import (
	"context"
	"log"
	userrepo "samase/user/repository"
	"time"

	"github.com/gomodule/redigo/redis"
)

func SaveEmailConfirmationCode(conn redis.Conn) userrepo.SaveEmailConfirmationCodeFunc {
	return func(ctx context.Context, code string, expireTime time.Duration) error {
		_, err := conn.Do("SET", code, true)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
