package userresthandler

import (
	"context"
	"database/sql"
	"net/http"
	samasemodels "samase/models"
	"samase/user"
	userrepo "samase/user/repository"
	usersqlrepo "samase/user/repository/sql"
	"samase/useremail"
	useremailrepo "samase/useremail/repository"
	useremailsqlrepo "samase/useremail/repository/sql"
	"samase/userpassword"
	userpasswordrepo "samase/userpassword/repository"
	userpasswordsqlrepo "samase/userpassword/repository/sql"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func InjectUserRESTHandler(conn *sql.DB, ee *echo.Echo) {
	createUser := usersqlrepo.CreateUser(conn)
	createUserEmail := useremailsqlrepo.CreateUserEmail(conn)
	createUserPassword := userpasswordsqlrepo.CreateUserPassword(conn)

	ee.POST("/users", CreateUser(createUser, createUserEmail, createUserPassword))
}

func CreateUser(
	createUser userrepo.CreateUserFunc,
	createUserEmail useremailrepo.CreateUserEmailFunc,
	createUserPassword userpasswordrepo.CreateUserPasswordFunc,
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
			Name          string `json:"name"`
			Fullname      string `json:"fullname"`
			Password      string `json:"password"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, samasemodels.RESTResponse{Message: "Error : invalid or empty body"})
		}
		us := user.User{
			Name:     post.Name,
			Fullname: post.Fullname,
		}
		us, err = createUser(ctx, us)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		passwordHash, err := hashPassword(post.Password)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		uspa := userpassword.UserPassword{
			UserID: us.ID,
			Hash:   passwordHash,
		}
		uspa, err = createUserPassword(ctx, uspa)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		usem := useremail.UserEmail{
			UserID: us.ID,
			Value:  post.Email,
		}
		usem, err = createUserEmail(ctx, usem)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		us.Email = &usem
		return ectx.JSON(http.StatusOK, us)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
