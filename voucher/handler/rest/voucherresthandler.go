package voucherresthandler

import (
	"database/sql"
	"net/http"
	"samase/voucher"
	vouchersqlrepo "samase/voucher/repository/sql"
	voucherservice "samase/voucher/service"

	"github.com/labstack/echo"
)

func InjectVoucherRESTHandler(conn *sql.DB, ee *echo.Echo) {
	getVouchers := voucherservice.GetVouchers(vouchersqlrepo.GetVouchers(conn))
	ee.GET("/vouchers", GetVouchers(getVouchers))
}

func GetVouchers(getVouchers voucherservice.GetVouchersFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		vos, err := getVouchers(ctx)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, vos)
	}
}

func CreateVoucher(
	createVoucher voucherservice.CreateVoucherFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		imgHeader, err := ectx.FormFile("image")
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		img, err := imgHeader.Open()

		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}

		name := ectx.FormValue("name")
		vo := voucher.Voucher{
			Name: name,
		}

		vo, err = createVoucher(ctx, vo, img)

		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		return ectx.JSON(http.StatusOK, vo)
	}
}
