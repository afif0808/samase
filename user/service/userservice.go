package userservice

import (
	"context"
	"errors"
	"fifentory/options"
	"log"
	"math/rand"
	samasemailservice "samase/samasemail/service"
	"samase/user"
	userrepo "samase/user/repository"
	"samase/useremail"
	useremailrepo "samase/useremail/repository"
	userpasswordrepo "samase/userpassword/repository"
	"time"
)

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
	gusfe userrepo.GetUserFetcherFunc,
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
		usfe := gusfe()
		uss, err := usfe.GetUsers(ctx, &opts)
		return len(uss) > 0, err
	}
}

func GetUserByEmail(gusfe userrepo.GetUserFetcherFunc) GetUserByEmailFunc {
	return func(ctx context.Context, email string) (*user.User, error) {
		usfe := gusfe()
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

func GetUserByID(gusfe userrepo.GetUserFetcherFunc) GetUserByIDFunc {
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
		usfe := gusfe()
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

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SendUserConfirmationEmail(
	sendEmail samasemailservice.SendEmailFunc,
	saveEmailConfirmationCode userrepo.SaveEmailConfirmationCodeFunc,
) SendUserConfirmationEmailFunc {
	return func(ctx context.Context, email string) error {
		code := randString(4)
		err := saveEmailConfirmationCode(ctx, code+"-"+email, 10800)
		if err != nil {
			return err
		}
		mailBody := `
			<h3>Terimakasih telah membuat akun member samase</h3>
			Masukan kode dibawah ke kolom kode di aplikasi <br>
			<h1>` + code + `</h1>
		`
		return sendEmail(ctx, []string{email}, "Konfirmasikan Email Anda", mailBody)
	}
}

func ConfirmUserEmail(
	checkEmailConfirmationCode userrepo.CheckEmailConfirmationCodeFunc,
	updateUserEmail useremailrepo.UpdateUserEmailsFunc,
) ConfirmUserEmailFunc {
	return func(ctx context.Context, email string, code string) error {
		exist, err := checkEmailConfirmationCode(ctx, code+"-"+email)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("Error : there's no such confirmation code or already expired")
		}
		fts := []options.Filter{
			options.Filter{
				Operator: "=",
				By:       "user_email.value",
				Value:    email,
			},
		}
		usem := useremail.UserEmail{
			Value:    email,
			Verified: true,
		}
		return updateUserEmail(ctx, usem, fts)
	}
}
