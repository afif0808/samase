package useremailsqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/useremail"
	useremailrepo "samase/useremail/repository"
)

const (
	userEmailTable       = "user_email"
	userEmailFields      = "user_email.user_id,user_email.value,user_email.verified"
	createUserEmailQuery = "INSERT " + userEmailTable + " SET user_id = ? , value = ? , verified = ?"
	getUserEmailsQuery   = "SELECT " + userEmailFields + " FROM " + userEmailTable
	updateUserEmailQuery = "UPDATE " + userEmailTable + " SET user_email.value = ? , user_email.verified = ?"
)

func CreateUserEmail(conn *sql.DB) useremailrepo.CreateUserEmailFunc {
	return func(ctx context.Context, usem useremail.UserEmail) (useremail.UserEmail, error) {
		_, err := conn.ExecContext(ctx, createUserEmailQuery, usem.UserID, usem.Value, usem.Verified)
		if err != nil {
			log.Println(err)
			return usem, err
		}
		return usem, nil
	}
}

func GetUserEmails(conn *sql.DB) useremailrepo.GetUserEmailsFunc {
	return func(ctx context.Context, opts *options.Options) ([]useremail.UserEmail, error) {
		optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
		rows, err := conn.QueryContext(ctx, getUserEmailsQuery+" "+optionsQuery, optionsArgs...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		usems := []useremail.UserEmail{}
		for rows.Next() {
			usem := useremail.UserEmail{}
			err := rows.Scan(&usem.UserID, &usem.Value, &usem.Verified)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			usems = append(usems, usem)
		}
		defer rows.Close()
		return usems, nil
	}
}

func UpdateUserEmails(conn *sql.DB) useremailrepo.UpdateUserEmailsFunc {
	return func(ctx context.Context, usem useremail.UserEmail, fts []options.Filter) error {
		filtersQuery, filtersOptions := options.ParseFiltersToSQLQuery(fts)
		filtersOptions = append([]interface{}{usem.Value, usem.Verified}, filtersOptions...)
		query := updateUserEmailQuery + " " + filtersQuery
		_, err := conn.ExecContext(ctx, query, filtersOptions...)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
