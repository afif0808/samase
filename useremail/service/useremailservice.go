package useremailservice

import (
	"context"
	"fifentory/options"
	useremailrepo "samase/useremail/repository"
)

type DoesEmailExistFunc func(ctx context.Context, email string) (bool, error)

func DoesEmailExist(
	getUserEmails useremailrepo.GetUserEmailsFunc,
) DoesEmailExistFunc {
	return func(ctx context.Context, email string) (bool, error) {
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "user_email.value",
					Operator: "=",
					Value:    email,
				},
			},
		}
		usems, err := getUserEmails(ctx, &opts)

		return len(usems) > 0, err
	}
}
