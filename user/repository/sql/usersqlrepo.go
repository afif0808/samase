package usersqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/user"
	userrepo "samase/user/repository"
)

const (
	userTable        = "user"
	createUserQuery  = "INSERT " + userTable + " SET name = ? , fullname = ?"
	updateUsersQuery = "UPDATE " + userTable + " SET name = ? , fullname = ?"
	deleteUsersQuery = "DELETE FROM " + userTable
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

func UpdateUsers(conn *sql.DB) userrepo.UpdateUsersFunc {
	return func(ctx context.Context, us user.User, fts []options.Filter) error {
		filtersQuery, filtersOptions := options.ParseFiltersToSQLQuery(fts)
		filtersOptions = append([]interface{}{us.Name, us.Fullname}, filtersOptions...)
		query := updateUsersQuery + " " + filtersQuery
		_, err := conn.ExecContext(ctx, query, filtersOptions...)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}

func DeleteUsers(conn *sql.DB) userrepo.DeleteUsersFunc {
	return func(ctx context.Context, fts []options.Filter) error {
		_, err := conn.ExecContext(ctx, deleteUsersQuery)
		if err != nil {
			log.Println(err)
		}
		return nil
	}
}
