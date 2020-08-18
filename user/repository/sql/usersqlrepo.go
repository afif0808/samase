package usersqlrepo

import (
	"context"
	"database/sql"
	"log"
	"samase/user"
	userrepo "samase/user/repository"
)

const (
	userTable       = "user"
	createUserQuery = "INSERT " + userTable + " SET name = ? , fullname = ?"
)

func CreateUser(conn *sql.DB) userrepo.CreateUserFunc {
	return func(ctx context.Context, us user.User) (user.User, error) {
		res, err := conn.ExecContext(ctx, createUserQuery, us.Name, us.Fullname)
		if err != nil {
			log.Println(err)
			return us, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			return us, err
		}
		us.ID = id
		return us, nil
	}
}
