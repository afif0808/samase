package uservoucherjunctionresthandler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	authenticationmiddleware "samase/authentication/middleware"
	"samase/uservoucherjunction"
	uservoucherjunctionsqlrepo "samase/uservoucherjunction/repository/sql"
	uservoucherjunctionservice "samase/uservoucherjunction/service"
	"samase/voucher"
	"strconv"

	"github.com/labstack/echo"
)

func InjectUserVoucherJunctionRESTHandler(conn *sql.DB, ee *echo.Echo) {
	createUserVoucherJunction := uservoucherjunctionsqlrepo.Createuservoucherjunction(conn)
	claimVoucher := uservoucherjunctionservice.ClaimVoucher(createUserVoucherJunction)
	ee.POST("/users/voucher/", ClaimVoucher(claimVoucher), authenticationmiddleware.InjectAuthenticate())
}

func ClaimVoucher(
	claimVoucher uservoucherjunctionservice.ClaimVoucherFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		user := ectx.Get("user").(map[string]interface{})
		var post struct {
			Voucher *voucher.Voucher `json:"voucher"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			log.Println(err)
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		userID, err := strconv.ParseInt(fmt.Sprint(user["id"]), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		usvo := uservoucherjunction.UserVoucherJunction{
			UserID:    userID,
			VoucherID: post.Voucher.ID,
		}
		usvo, err = claimVoucher(ctx, usvo)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, "Claim success")
	}
}
