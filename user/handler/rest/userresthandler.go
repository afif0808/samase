package userresthandler

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	samasemodels "samase/models"
	notificationsqlrepo "samase/notification/repository/sql"
	notificationservice "samase/notification/service"
	samasemailservice "samase/samasemail/service"
	"samase/user"
	userredisrepo "samase/user/repository/redis"
	usersqlrepo "samase/user/repository/sql"
	userservice "samase/user/service"
	"samase/useremail"
	useremailsqlrepo "samase/useremail/repository/sql"
	"samase/userpassword"
	userpasswordsqlrepo "samase/userpassword/repository/sql"
	userpasswordservice "samase/userpassword/service"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   50,
		MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ":6379")
			// Connection error handling
			if err != nil {
				log.Printf("ERROR: fail initializing the redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func InjectUserRESTHandler(conn *sql.DB, ee *echo.Echo) {
	createUser := usersqlrepo.CreateUser(conn)
	createUserEmail := useremailsqlrepo.CreateUserEmail(conn)
	createUserPassword := userpasswordsqlrepo.CreateUserPassword(conn)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", mailer.FormatAddress("samaseapp@gmail.com", "SamaseApp"))
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "samaseapp@gmail.com", "samaseapp0808")
	sendEmail := samasemailservice.SendEmail(dialer, mailer)

	rp := redisPool()
	rc, _ := rp.Dial()

	saveEmailConfirmationCode := userredisrepo.SaveEmailConfirmationCode(rc)
	sendUserConfirmationEmail := userservice.SendUserConfirmationEmail(sendEmail, saveEmailConfirmationCode)
	sendWelcomeNotification := notificationservice.SendWelcomeNotification(notificationsqlrepo.CreateNotification(conn))
	ee.POST("/users", CreateUser(userservice.CreateUser(createUser, createUserEmail, createUserPassword), sendWelcomeNotification))
	gussf := usersqlrepo.GetUserSQLFetcher(conn)
	// ussf := usersqlrepo.NewUserSQLFetcher(conn)
	doesNameExist := userservice.DoesNameExist(gussf)

	ee.GET("/users/name/exists/:name", DoesNameExist(doesNameExist))
	updateUser := userservice.UpdateUser(usersqlrepo.UpdateUsers(conn), useremailsqlrepo.UpdateUserEmails(conn))
	ee.POST("/users/:id", UpdateUser(updateUser))
	confirmUserEmail := userservice.ConfirmUserEmail(
		userredisrepo.CheckEmailConfirmationCode(rc),
		useremailsqlrepo.UpdateUserEmails(conn),
	)
	ee.POST("/users/email/confirm", ConfirmUserEmail(confirmUserEmail))
	ee.POST("/users/email/confirmationcode", SendUserEmailConfirmation(sendUserConfirmationEmail))

	savePasswordRecoveryCode := userredisrepo.SavePasswordRecoveryCode(rc)
	sendPasswordRecoveryCode := userservice.SendPasswordRecoveryCode(sendEmail, savePasswordRecoveryCode)

	ee.POST("/users/password/recoverycode", SendPasswordRecoveryCode(sendPasswordRecoveryCode))

	confirmPasswordRecovery := userservice.ConfirmPasswordRecoveryCode(userredisrepo.CheckPasswordRecoveryCode(rc), userredisrepo.RemovePasswordRecoveryCode(rc))

	getUserByEmail := userservice.GetUserByEmail(gussf)

	ee.POST("/users/password/confirmrecoverycode", ConfirmPasswordRecoveryCode(confirmPasswordRecovery, getUserByEmail))

	saveUserIDByCode := userredisrepo.SaveUserIDByCode(rc)

	sendAccountPasswordRecoveryLink := userservice.SendAccountPasswordRecoveryLink(
		saveUserIDByCode,
		getUserByEmail,
		sendEmail,
		"https://nother.samasecentro.com",
	)

	ee.POST("/users/password/recoverylink", SendAccountPasswordRecoveryLink(sendAccountPasswordRecoveryLink))

	recoverAccountPassword := userservice.RecoverUserPassword(userredisrepo.RetrieveUserIDByCode(rc), userpasswordservice.UpdateUserPassword(userpasswordsqlrepo.UpdateUserPassword(conn)))

	ee.POST("/users/password/recover", RecoverAccountPassword(recoverAccountPassword))
	unregisterUserWebSocket := userservice.UnregisterUserWebSocket()
	registerUserWebSocket := userservice.RegisterUserWebSocket()
	ee.GET("/userws", UserWebSocket(registerUserWebSocket, unregisterUserWebSocket))

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
	sendWelcomeNotification notificationservice.SendWelcomeNotificationFunc,
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
		name = strings.Replace(name, " ", "", -1)
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
		err = sendWelcomeNotification(ctx, us.ID)
		if err != nil {

			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		// go func() {
		// 	err = sendUserConfirmationEmail(ctx, post.Email)
		// 	log.Println(err)
		// }()
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
			Fullname      string `json:"fullname"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
			Name          string `json:"name"`
		}
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		us := user.User{ID: id, Fullname: post.Fullname, Name: post.Name, Email: &useremail.UserEmail{Value: post.Email, Verified: post.EmailVerified, UserID: id}}
		err = updateUser(ctx, us)
		if err != nil {
			log.Println(err)
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

func ConfirmUserEmail(confirmEmail userservice.ConfirmUserEmailFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Code  string `json:"code"`
			Email string `json:"email"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = confirmEmail(ctx, post.Email, post.Code)
		log.Println(err)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

func SendUserEmailConfirmation(
	sendUserConfirmationEmail userservice.SendUserConfirmationEmailFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Email string `json:"email"`
		}
		err := ectx.Bind(&post)
		log.Println(err, post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = sendUserConfirmationEmail(ctx, post.Email)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

func SendPasswordRecoveryCode(
	sendCode userservice.SendPasswordRecoveryCodeFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Email string `json:"email"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = sendCode(ctx, post.Email)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

func SendAccountPasswordRecoveryLink(
	recover userservice.SendAccountPasswordRecoveryLinkFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Email string `json:"email"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = recover(ctx, post.Email)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

func ConfirmPasswordRecoveryCode(
	confirmCode userservice.ConfirmPasswordRecoveryCodeFunc,
	getUserByEmail userservice.GetUserByEmailFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		log.Println(post)
		err = confirmCode(ctx, post.Email, post.Code)
		if err != nil {
			log.Println(err)
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		us, err := getUserByEmail(ctx, post.Email)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, us)
	}
}

func RecoverAccountPassword(
	recoverUserPassword userservice.RecoverUserPasswordFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Code     string `json:"code"`
			Password string `json:"password"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		log.Println(post)
		err = recoverUserPassword(ctx, post.Code, post.Password)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

var upgrader = websocket.Upgrader{}

func UserWebSocket(
	register userservice.RegisterUserWebSocketFunc,
	unregister userservice.UnregisterUserWebSocketFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		// ctx := ectx.Request().Context()
		log.Println("Wow")
		ws, err := upgrader.Upgrade(ectx.Response(), ectx.Request(), nil)
		if err != nil {
			log.Println("nol", err)
			return ectx.JSON(http.StatusBadRequest, nil)
		}

		err = register(ws)
		if err != nil {
			log.Println("satu", err)
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		defer func() {
			ws.Close()
			err = unregister(ws)
			log.Println("dua", err)
		}()
		for {
			_, _, err = ws.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
		}
		return nil
	}
}
