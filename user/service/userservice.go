package userservice

import (
	"context"
	"fifentory/options"
	"log"
	"samase/user"
	userrepo "samase/user/repository"
	useremailrepo "samase/useremail/repository"
	userpasswordrepo "samase/userpassword/repository"
)

type DoesNameExistFunc func(ctx context.Context, name string) (bool, error)

type GetUserByEmailFunc func(ctx context.Context, email string) (*user.User, error)
type GetUserByIDFunc func(ctx context.Context, id int64) (*user.User, error)

type GetAllUsersFunc func(ctx context.Context) ([]user.User, error)


type CreateUserFunc func(ctx context.Context, us user.User) (user.User, error)

type UpdateUserFunc func(ctx context.Context, us user.User) error

func CreateUser(
	createUser userrepo.CreateUserFunc,
	createUserEmail useremailrepo.CreateUserEmailFunc,
	createUserPassword userpasswordrepo.CreateUserPasswordFunc,
) CreateUserFunc {
	return func(ctx context.Context, us user.User) (user.User, error) {
		us, err := createUser(ctx, us)
		if err != nil {
			return us, err
		}
		if us.Email != nil {
			us.Email.UserID = us.ID
			_, err = createUserEmail(ctx, *us.Email)
			if err != nil {
				return us, err
			}
		}
		if us.Password != nil {
			us.Password.UserID = us.ID
			_, err = createUserPassword(ctx, *us.Password)
			if err != nil {
				return us, err
			}
		}
		return us, nil
	}
}

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

func GetUserByEmail(usfe userrepo.UserFetcher) GetUserByEmailFunc {
	return func(ctx context.Context, email string) (*user.User, error) {
		usfe.WithEmail()
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "user_email.value",
					Value:    email,
					Operator: "=",
				},
			},
		}
		uss, err := usfe.GetUsers(ctx, &opts)
		if err != nil {
			return nil, err
		}
		if len(uss) <= 0 {
			return nil, nil
		}
		return &uss[0], nil
	}
}

func GetUserByID(usfe userrepo.UserFetcher) GetUserByIDFunc {
	return func(ctx context.Context, id int64) (*user.User, error) {
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "user.id",
					Value:    id,
					Operator: "=",
				},
			},
		}
		usfe.WithEmail()
		uss, err := usfe.GetUsers(ctx, &opts)
		if err != nil {
			return nil, err
		}
		if len(uss) <= 0 {
			return nil, nil
		}
		return &uss[0], nil
	}
}

func UpdateUser(
	updateUser userrepo.UpdateUsersFunc,
	updateUserEmail useremailrepo.UpdateUserEmailsFunc,
) UpdateUserFunc {
	return func(ctx context.Context, us user.User) error {
		log.Println(us)
		if us.Email != nil {
			usemfts := []options.Filter{
				options.Filter{
					By:       "user_email.user_id",
					Operator: "=",
					Value:    us.Email.UserID,
				},
			}
			err := updateUserEmail(ctx, *us.Email, usemfts)
			if err != nil {
				return err
			}
		}
		usfts := []options.Filter{
			options.Filter{
				By:       "user.id",
				Operator: "=",
				Value:    us.ID,
			},
		}
		return updateUser(ctx, us, usfts)
	}
}
