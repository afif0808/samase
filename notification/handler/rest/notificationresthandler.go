package notificationresthandler

import (
	"database/sql"
	"net/http"
	notificationsqlrepo "samase/notification/repository/sql"
	notificationservice "samase/notification/service"
	usersqlrepo "samase/user/repository/sql"
	"strconv"

	"github.com/labstack/echo"
)

func InjectNotificationRESTHandler(conn *sql.DB, ee *echo.Echo) {
	notfsf := notificationsqlrepo.NewNotificationSQLFetcher(conn)
	getNotificationsByUserID := notificationservice.GetNotificationsByUserID(&notfsf)
	ee.GET("/users/:id/notifications", GetNotificationsByUserID(getNotificationsByUserID))

	ussf := usersqlrepo.NewUserSQLFetcher(conn)
	createNotification := notificationsqlrepo.CreateNotification(conn)
	createNotificationForAllUsers := notificationservice.CreateNotificationForAllUsers(&ussf, createNotification)

	ee.POST("/notifications", CreateNotificationForAllUsers(createNotificationForAllUsers))
}

func GetNotificationsByUserID(getNotificationsByUserID notificationservice.GetNotificationsByUserIDFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		notfs, err := getNotificationsByUserID(ctx, id)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, notfs)
	}
}

func createNotificationForAllUsers() {

}

func CreateNotificationForAllUsers(
	createNotificationForALlUsers notificationservice.CreateNotificationForAllUsersFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Title   string `json:"title"`
			Message string `json:"message"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = createNotificationForALlUsers(ctx, post.Title, post.Message)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusCreated, nil)
	}
}
