package userpasswordsqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/userpassword"
	userpasswordrepo "samase/userpassword/repository"
)

const (
	userPasswordTable       = "user_password"
	fields                  = "user_id,hash"
	createUserPasswordQuery = "INSERT " + userPasswordTable + " SET user_id = ? , hash = ?"
	updateUserPasswordQuery = "UPDATE " + userPasswordTable + " SET hash = ?  "
	getUserPasswordsQuery   = "SELECT " + fields + " FROM " + userPasswordTable
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

func UpdateUserPassword(conn *sql.DB) userpasswordrepo.UpdateUserPasswordFunc {
	return func(ctx context.Context, uspa userpassword.UserPassword, fts []options.Filter) error {
		filtersQuery, filterOptions := options.ParseFiltersToSQLQuery(fts)
		filterOptions = append([]interface{}{uspa.Hash}, filterOptions...)
		query := updateUserPasswordQuery + " " + filtersQuery
		_, err := conn.ExecContext(ctx, query, filterOptions...)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}

func GetUserPasswords(conn *sql.DB) userpasswordrepo.GetUserPasswordsFunc {
	return func(ctx context.Context, opts *options.Options) ([]userpassword.UserPassword, error) {
		optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
		query := getUserPasswordsQuery + " " + optionsQuery
		rows, err := conn.QueryContext(ctx, query, optionsArgs...)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		uspas := []userpassword.UserPassword{}
		for rows.Next() {
			uspa := userpassword.UserPassword{}
			rows.Scan(&uspa.UserID, &uspa.Hash)
			uspas = append(uspas, uspa)
		}
		return uspas, nil
	}
}
