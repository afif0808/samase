package useremailresthandler

import (
	"database/sql"
	"log"
	"net/http"
	useremailsqlrepo "samase/useremail/repository/sql"
	useremailservice "samase/useremail/service"

	"github.com/labstack/echo"
)

func InjectUserEmailRESTHandler(conn *sql.DB, ee *echo.Echo) {
	getUserEmails := useremailsqlrepo.GetUserEmails(conn)
	doesEmailExist := useremailservice.DoesEmailExist(getUserEmails)
	ee.GET("/users/email/exists/:email", DoesEmailExist(doesEmailExist))
}

func DoesEmailExist(
	doesEmailExist useremailservice.DoesEmailExistFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		email := ectx.Param("email")
		exist, err := doesEmailExist(ctx, email)
		if err != nil {
			log.Println(err)
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		var resp struct {
			Exist bool `json:"exist"`
		}
		resp.Exist = exist
		return ectx.JSON(http.StatusOK, resp)
	}
}
