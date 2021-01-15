package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	authenticationresthandler "samase/authentication/handler/rest"
	eventresthandler "samase/event/handler/rest"
	imagemanagerresthandler "samase/imagemanager/handler/rest"
	notificationresthandler "samase/notification/handler/rest"
	userresthandler "samase/user/handler/rest"
	userservice "samase/user/service"
	useremailresthandler "samase/useremail/handler/rest"
	userpasswordresthandler "samase/userpassword/handler/rest"
	uservoucherjunctionresthandler "samase/uservoucherjunction/handler/rest"
	voucherresthandler "samase/voucher/handler/rest"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	// msg := messaging.Message{
	// 	Topic: "topic",
	// 	Notification: &messaging.Notification{
	// 		Title: "Another Notification",
	// 		Body:  "Another Body",
	// 	},
	// }
	// _, err = msging.Send(context.Background(), &msg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	conn, err := sql.Open("mysql", "root:@tcp(localhost:3306)/"+fmt.Sprint("samaseapp")+"?parseTime=true")
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

	dsn := "root:@tcp(localhost:3306)/samase?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	userresthandler.InjectUserRESTHandler(conn, ee)
	notificationresthandler.InjectNotificationRESTHandler(conn, ee)
	authenticationresthandler.InjectAuthenticationRESTHandler(conn, ee, rc)
	uservoucherjunctionresthandler.InjectUserVoucherJunctionRESTHandler(conn, ee)
	useremailresthandler.InjectUserEmailRESTHandler(conn, ee)
	userpasswordresthandler.InjectUserPasswordRESTHandler(conn, ee)
	voucherresthandler.InjectVoucherRESTHandler(conn, gormDB, ee)
	eventresthandler.InjectEventRESTHandler(conn, gormDB, ee)
	imagemanagerresthandler.InjectImageManagerRESTHandler(ee)
	// ee.Static("/assets", "/root/go/src/samase/assets")
	ee.Static("/vouchers/images", "/root/go/src/samase/assets/vouchers")
	ee.Static("/events/images", "/root/go/src/samase/assets/events")
	ee.Static("/images", "C:/Users/bayur/go/src/samase/assets/images")
	// fs := http.FileServer(http.Dir("/media/afif0808/data/go/src/samase/assets"))
	// ee.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", fs)))

	ee.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// ee.GET("/ws", hello)

	go userservice.WebsocketStream()

	ee.Start(":757")
}
