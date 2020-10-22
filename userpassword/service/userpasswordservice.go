package userpasswordservice

import (
	"context"
	"fifentory/options"
	"samase/userpassword"
	userpasswordrepo "samase/userpassword/repository"
)

func UpdateUserPassword(updateUserPassword userpasswordrepo.UpdateUserPasswordFunc) UpdateUserPasswordFunc {
	return func(ctx context.Context, uspa userpassword.UserPassword) error {
		fts := []options.Filter{
			options.Filter{
				By:       "user_password.user_id",
				Operator: "=",
				Value:    uspa.UserID,
			},
		}
		return updateUserPassword(ctx, uspa, fts)
	}
}
func GetUserPasswordByUserID(
	getUserPasswords userpasswordrepo.GetUserPasswordsFunc,
) GetUserPasswordByUserIDFunc {
	return func(ctx context.Context, id int64) (*userpassword.UserPassword, error) {
		fts := []options.Filter{
			options.Filter{
				By:       "user_password.user_id",
				Operator: "=",
				Value:    id,
			},
		}
		opts := options.Options{Filters: fts}
		uspas, err := getUserPasswords(ctx, &opts)
		if err != nil {
			return nil, err
		}
		if len(uspas) < 1 {
			return nil, nil
		}
		uspa := uspas[0]
		return &uspa, nil
	}
}
