package userpasswordsqlrepo

import (
	"context"
	"database/sql"
	"log"
	"samase/userpassword"
	userpasswordrepo "samase/userpassword/repository"
)

const (
	userPasswordTable       = "user_password"
	createUserPasswordQuery = "INSERT " + userPasswordTable + " SET user_id = ? , hash = ?"
)

func CreateUserPassword(conn *sql.DB) userpasswordrepo.CreateUserPasswordFunc {
	return func(ctx context.Context, uspa userpassword.UserPassword) (userpassword.UserPassword, error) {
		_, err := conn.ExecContext(ctx, createUserPasswordQuery, uspa.UserID, uspa.Hash)
		if err != nil {
			log.Println(err)
		}
		return uspa, err
	}
}
