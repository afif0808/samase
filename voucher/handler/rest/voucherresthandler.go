package voucherresthandler

import (
	"database/sql"
	"net/http"
	"samase/voucher"
	vouchersqlrepo "samase/voucher/repository/sql"
	voucherservice "samase/voucher/service"
	"strconv"

	"github.com/labstack/echo"
)

func InjectVoucherRESTHandler(conn *sql.DB, ee *echo.Echo) {
	getVouchers := voucherservice.GetVouchers(vouchersqlrepo.GetVouchers(conn))
	ee.GET("/vouchers", GetVouchers(getVouchers))
	deleteVoucherByID := voucherservice.DeleteVoucherByID(vouchersqlrepo.DeleteVouchers(conn))
	ee.DELETE("/vouchers/:id", DeleteVoucherByID(deleteVoucherByID))
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

func DeleteVoucherByID(deleteVoucherByID voucherservice.DeleteVoucherByIDFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = deleteVoucherByID(ctx, id)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}
