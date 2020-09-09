package userservice

import (
	"context"
	"fifentory/options"
	userrepo "samase/user/repository"
)

type DoesNameExistFunc func(ctx context.Context, name string) (bool, error)

func DoesNameExist(
	userFetcher userrepo.UserFetcher,
) DoesNameExistFunc {
	return func(ctx context.Context, name string) (bool, error) {
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "user.name",
					Value:    name,
					Operator: "=",
				},
			},
		}
		uss, err := userFetcher.GetUsers(ctx, &opts)
		return len(uss) > 0, err
	}
}
