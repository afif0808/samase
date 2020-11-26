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
	"samase/userpassword"
	userpasswordrepo "samase/userpassword/repository"
	userpasswordservice "samase/userpassword/service"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func SendPasswordRecoveryCode(
	sendEmail samasemailservice.SendEmailFunc,
	savePasswordRecovery userrepo.SavePasswordRecoveryCodeFunc,
) SendPasswordRecoveryCodeFunc {
	return func(ctx context.Context, email string) error {
		code := randString(4)
		err := savePasswordRecovery(ctx, "PR-"+code+"-"+email, 10800)
		if err != nil {
			return err
		}
		mailBody := `
			<h3>Berikut merupakan kode untuk mengatur ulang password anda</h3>
			Masukan kode dibawah ke kolom kode di aplikasi <br>
			<h1>` + code + `</h1>
			Jangan beritahukan kode ini kepada siapa pun
		`
		return sendEmail(ctx, []string{email}, "Atur ulang kata sandi", mailBody)
	}
}

func ConfirmPasswordRecoveryCode(
	check userrepo.CheckPasswordRecoveryCodeFunc,
	removeCode userrepo.RemovePasswordRecoveryCodeFunc,
) ConfirmPasswordRecoveryCodeFunc {
	return func(ctx context.Context, email, code string) error {
		code = "PR-" + code + "-" + email
		exist, err := check(ctx, code)
		if err != nil {
			return errors.New("Error : there was a problem checking the code")
		}
		if !exist {
			return errors.New("Error : the given code not found , either it's already expired or never existed in the first place")
		}
		err = removeCode(ctx, code)
		if err != nil {
			return errors.New("Error : failed to remove the recovery code")
		}
		return nil
	}
}

func SendAccountPasswordRecoveryLink(
	saveUserIDByCode userrepo.SaveUserIDByCodeFunc,
	getUserByEmail GetUserByEmailFunc,
	sendEmail samasemailservice.SendEmailFunc,
	baseURL string,
) SendAccountPasswordRecoveryLinkFunc {
	return func(ctx context.Context, email string) error {
		us, err := getUserByEmail(ctx, email)
		if err != nil {
			return nil
		}

		code := randString(8)

		err = saveUserIDByCode(ctx, code, us.ID)

		if err != nil {
			return nil
		}
		link := baseURL + "?code=" + code
		body := `
			<h4>Klik link dibawah ini untuk mengatur ulang kata sandi anda</h4>
			` + link + `<br>
			Jangan sebarkan link ini kepada siapa pun
		`
		err = sendEmail(ctx, []string{email}, "Atur ulang kata sandi anda", body)
		return err
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RecoverUserPassword(
	retrieveUserIDByCode userrepo.RetrieveUserIDByCodeFunc,
	updateUsserPassword userpasswordservice.UpdateUserPasswordFunc,
) RecoverUserPasswordFunc {
	return func(ctx context.Context, code, password string) error {
		id, err := retrieveUserIDByCode(ctx, code)
		log.Println(id, err)
		if err != nil {
			return err
		}

		passwordHash, err := hashPassword(password)

		if err != nil {
			return err
		}

		uspa := userpassword.UserPassword{
			UserID: id,
			Value:  passwordHash,
		}
		return updateUsserPassword(ctx, uspa)
	}
}
