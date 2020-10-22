package userresthandler

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	samasemodels "samase/models"
	samasemailservice "samase/samasemail/service"
	"samase/user"
	usersqlrepo "samase/user/repository/sql"
	userservice "samase/user/service"
	"samase/useremail"
	useremailsqlrepo "samase/useremail/repository/sql"
	"samase/userpassword"
	userpasswordsqlrepo "samase/userpassword/repository/sql"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func InjectUserRESTHandler(conn *sql.DB, ee *echo.Echo) {
	createUser := usersqlrepo.CreateUser(conn)
	createUserEmail := useremailsqlrepo.CreateUserEmail(conn)
	createUserPassword := userpasswordsqlrepo.CreateUserPassword(conn)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", mailer.FormatAddress("afifsamase@gmail.com", "SamaseApp"))
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "afifsamase@gmail.com", "samaseafif87")
	sendEmail := samasemailservice.SendEmail(dialer, mailer)

	ee.POST("/users", CreateUser(userservice.CreateUser(createUser, createUserEmail, createUserPassword), sendEmail))

	ussf := usersqlrepo.NewUserSQLFetcher(conn)
	doesNameExist := userservice.DoesNameExist(&ussf)

	ee.GET("/users/name/exists/:name", DoesNameExist(doesNameExist))
	updateUser := userservice.UpdateUser(usersqlrepo.UpdateUsers(conn), useremailsqlrepo.UpdateUserEmails(conn))
	ee.POST("/users/:id", UpdateUser(updateUser))
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

func CreateUser(
	createUser userservice.CreateUserFunc,
	sendEmail samasemailservice.SendEmailFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {

		ctx := ectx.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}
		if _, isWithDeadline := ctx.Deadline(); !isWithDeadline {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Second*5)
			defer cancel()
		}
		var post struct {
			Fullname      string `json:"fullname"`
			Password      string `json:"password"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
		}
		err := ectx.Bind(&post)
		log.Println(post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, samasemodels.RESTResponse{Message: "Error : invalid or empty body"})
		}

		name := strings.ToLower(post.Fullname)
		name = strings.ReplaceAll(name, " ", "")

		passwordHash, err := hashPassword(post.Password)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		us := user.User{
			Name:     name,
			Fullname: post.Fullname,
			Email: &useremail.UserEmail{
				Value: post.Email,
			},
			Password: &userpassword.UserPassword{
				Hash: passwordHash,
			},
		}
		us, err = createUser(ctx, us)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		code := "<h1>" + randString(4) + "</h1>"
		go func() {
			err = sendEmail(ctx, []string{post.Email}, "Konfirmasi email anda", code)
			log.Println(err)
		}()

		return ectx.JSON(http.StatusOK, us)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func DoesNameExist(
	doesNameExist userservice.DoesNameExistFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		name := ectx.Param("name")
		if name == "" {
			return ectx.JSON(http.StatusBadRequest, samasemodels.RESTResponse{Message: "Error : name cannot be empty"})
		}
		exist, err := doesNameExist(ctx, name)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		var resp struct {
			Exist bool `json:"exist"`
		}
		resp.Exist = exist
		return ectx.JSON(http.StatusOK, resp)
	}
}

func UpdateUser(updateUser userservice.UpdateUserFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Fullname string `json:"fullname"`
			Email    string `json:"email"`
		}
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		us := user.User{ID: id, Fullname: post.Fullname, Email: &useremail.UserEmail{Value: post.Email, UserID: id}}
		err = updateUser(ctx, us)
		if err != nil {
			log.Println(err)
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}
