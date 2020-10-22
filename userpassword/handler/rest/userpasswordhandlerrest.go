package userpasswordresthandler

import (
	"database/sql"
	"log"
	"net/http"
	"samase/userpassword"
	userpasswordrepo "samase/userpassword/repository"
	userpasswordsqlrepo "samase/userpassword/repository/sql"
	userpasswordservice "samase/userpassword/service"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

func InjectUserPasswordRESTHandler(conn *sql.DB, ee *echo.Echo) {
	getUserPasswordByID := userpasswordservice.GetUserPasswordByUserID(userpasswordsqlrepo.GetUserPasswords(conn))
	createUserPassword := userpasswordsqlrepo.CreateUserPassword(conn)
	updateUserPassword := userpasswordservice.UpdateUserPassword(userpasswordsqlrepo.UpdateUserPassword(conn))
	ee.POST("/users/password/:id", SetUserPassword(updateUserPassword, createUserPassword, getUserPasswordByID))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func SetUserPassword(
	updatePassword userpasswordservice.UpdateUserPasswordFunc,
	createUserPassword userpasswordrepo.CreateUserPasswordFunc,
	getUserPasswordByID userpasswordservice.GetUserPasswordByUserIDFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Password string `json:"password"`
		}
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		log.Println(post)

		uspa, err := getUserPasswordByID(ctx, id)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		passwordHash, err := hashPassword(post.Password)

		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		if uspa == nil {
			uspa = &userpassword.UserPassword{
				Hash:   passwordHash,
				UserID: id,
			}
			*uspa, err = createUserPassword(ctx, *uspa)
			if err != nil {
				return ectx.JSON(http.StatusInternalServerError, nil)
			}
		} else {
			uspa.Hash = passwordHash
			updatePassword(ctx, *uspa)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}
