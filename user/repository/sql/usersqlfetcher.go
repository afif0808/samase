package usersqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/user"
	userrepo "samase/user/repository"
	"samase/useremail"
	"samase/userpassword"
	"samase/userpoint"
	"samase/voucher"
)

type receiver struct {
	User         *user.User
	UserEmail    *useremail.UserEmail
	UserPassword *userpassword.UserPassword
	UserPoint    *userpoint.UserPoint
	Vouchers     *[]voucher.Voucher
}

type UserSQLFetcher struct {
	joins    string
	scanDest []interface{}
	fields   string
	Receiver *receiver
	conn     *sql.DB
}

func GetUserSQLFetcher(conn *sql.DB) userrepo.GetUserFetcherFunc {
	return func() userrepo.UserFetcher {
		ussf := NewUserSQLFetcher(conn)
		return &ussf
	}
}

func NewUserSQLFetcher(conn *sql.DB) UserSQLFetcher {
	ussf := UserSQLFetcher{
		Receiver: &receiver{User: &user.User{}},
		conn:     conn,
	}
	return ussf
}

func (ussf *UserSQLFetcher) GetUsers(ctx context.Context, opts *options.Options) ([]user.User, error) {
	ussf.fields += "user.id,user.name,user.fullname"

	ussf.scanDest = append(
		ussf.scanDest,
		&ussf.Receiver.User.ID,
		&ussf.Receiver.User.Name,
		&ussf.Receiver.User.Fullname,
	)

	defer func() {
		ussf.fields = ""
		ussf.joins = ""
		ussf.scanDest = []interface{}{}
	}()

	optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
	query := "SELECT " + ussf.fields + " FROM " + userTable + " " + ussf.joins + " " + optionsQuery
	rows, err := ussf.conn.QueryContext(ctx, query, optionsArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	uss := []user.User{}
	for rows.Next() {
		err := rows.Scan(ussf.scanDest...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		re := ussf.Receiver
		us := user.User{
			ID:       re.User.ID,
			Name:     re.User.Name,
			Fullname: re.User.Fullname,
		}

		if ussf.Receiver.UserPassword != nil {
			us.Password = &userpassword.UserPassword{
				Hash: ussf.Receiver.UserPassword.Hash,
			}
		}

		if ussf.Receiver.UserEmail != nil {
			us.Email = &useremail.UserEmail{
				Value:    ussf.Receiver.UserEmail.Value,
				Verified: ussf.Receiver.UserEmail.Verified,
			}
		}

		uss = append(uss, us)
	}
	return uss, nil
}

func (ussf *UserSQLFetcher) WithEmail() {
	ussf.Receiver.UserEmail = &useremail.UserEmail{}
	ussf.fields += "user_email.value,user_email.verified,"
	ussf.joins += " INNER JOIN user_email ON user.id = user_email.user_id "
	ussf.scanDest = append(
		ussf.scanDest,
		&ussf.Receiver.UserEmail.Value,
		&ussf.Receiver.UserEmail.Verified,
	)
}
func (ussf *UserSQLFetcher) WithPassword() {
	ussf.Receiver.UserPassword = &userpassword.UserPassword{}
	ussf.fields += "user_password.hash,"
	ussf.joins += "INNER JOIN user_password ON user.id = user_password.user_id "
	ussf.scanDest = append(
		ussf.scanDest,
		&ussf.Receiver.UserPassword.Hash,
	)
}

func (ussf *UserSQLFetcher) WithPoint() {
	ussf.Receiver.UserPoint = &userpoint.UserPoint{}
	ussf.fields += "user_point.value,"
	ussf.joins += "INNER JOIN user_point ON user.id = user_point.user_id"
	ussf.scanDest = append(
		ussf.scanDest,
		&ussf.Receiver.UserPoint.Value,
	)
}

func (ussf *UserSQLFetcher) WithVouchers() {

}

func (ussf *UserSQLFetcher) AddJoins(joins string) {
	ussf.joins += joins
}
func (ussf *UserSQLFetcher) AddFields(fields string) {
	ussf.fields += fields
}
func (ussf *UserSQLFetcher) AddScanDest(dest []interface{}) {
	ussf.scanDest = append(ussf.scanDest, dest...)
}
