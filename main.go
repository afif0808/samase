package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	authenticationresthandler "samase/authentication/handler/rest"
	notificationresthandler "samase/notification/handler/rest"
	userresthandler "samase/user/handler/rest"
	useremailresthandler "samase/useremail/handler/rest"
	userpasswordresthandler "samase/userpassword/handler/rest"
	uservoucherjunctionresthandler "samase/uservoucherjunction/handler/rest"

	"github.com/gomodule/redigo/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func main() {
	// mc := samasemail.MailConfig{
	// 	Email:     "afifsamase@gmail.com",
	// 	Password:  "samaseafif87",
	// 	SMTP_HOST: "smtp.gmail.com",
	// 	SMTP_PORT: 587,
	// 	Mailer:    gomail.NewMessage(),
	// }
	// log.Println(samasemailservice.SendEmail(mc)(context.Background(), []string{"afif.panai@gmail.com"}, "Hello WOrld", "<h1>Comment alez vous?</h1>"))
	// config := map[string]interface{}{}
	// configFile, err := os.Open("../src/fifentory/config.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// json.NewDecoder(configFile).Decode(&config)

	conn, err := sql.Open("mysql", "root:@tcp(localhost:3306)/"+fmt.Sprint("samase")+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	ee := echo.New()
	// parseJWT := jsonwebtokenservice.ParseJWT()
	// authenticate := authenticationmiddleware.Authenticate(
	// 	parseJWT,
	// 	[]byte("itssignaturekey"),
	// 	jwt.SigningMethodHS256,
	// )

	// conf := &oauth2.Config{
	// 	ClientID:     "744967159273-34m45ct0lc0pas8ao3a4d9o7o1v6b7lp.apps.googleusercontent.com",
	// 	ClientSecret: "pe3BLwj424W3rxnV11kwXA8S",
	// 	RedirectURL:  "http://localhost:3000",
	// 	Scopes: []string{
	// 		"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	// 	},
	// 	Endpoint: google.Endpoint,
	// }
	rp := redisPool()
	rc, _ := rp.Dial()

	userresthandler.InjectUserRESTHandler(conn, ee)
	notificationresthandler.InjectNotificationRESTHandler(conn, ee)
	authenticationresthandler.InjectAuthenticationRESTHandler(conn, ee, rc)
	uservoucherjunctionresthandler.InjectUserVoucherJunctionRESTHandler(conn, ee)
	useremailresthandler.InjectUserEmailRESTHandler(conn, ee)
	userpasswordresthandler.InjectUserPasswordRESTHandler(conn, ee)
	ee.Static("/assets", "/media/afif0808/data/go/src/samase/assets")
	ee.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	ee.Start(":777")
}
