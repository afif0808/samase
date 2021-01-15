package voucherresthandler

import (
	"database/sql"
	"net/http"
	"samase/voucher"
	vouchersqlrepo "samase/voucher/repository/sql"
	voucherservice "samase/voucher/service"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func InjectVoucherRESTHandler(conn *sql.DB, gormDB *gorm.DB, ee *echo.Echo) {
	getVouchers := voucherservice.GetVouchers(vouchersqlrepo.GetVouchers(conn))
	ee.GET("/vouchers", GetVouchers(getVouchers))
	deleteVoucherByID := voucherservice.DeleteVoucherByID(vouchersqlrepo.DeleteVouchers(conn))
	ee.DELETE("/vouchers/:id", DeleteVoucherByID(deleteVoucherByID))
	createVoucher := voucherservice.CreateVoucher(vouchersqlrepo.CreateVoucher(conn))
	ee.POST("/vouchers", CreateVoucher(createVoucher))
	// gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	updateVouchers := vouchersqlrepo.UpdateVouchers(gormDB)
	updateVoucherByID := voucherservice.UpdateVoucherByID(updateVouchers)
	ee.POST("/vouchers/:id", UpdateVoucherByID(updateVoucherByID))
}

func GetVouchers(getVouchers voucherservice.GetVouchersFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		search := ectx.Request().URL.Query().Get("s")
		vos, err := getVouchers(ctx, search)
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
		var post struct {
			Voucher voucher.Voucher `json:"voucher"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		vo, err := createVoucher(ctx, post.Voucher)
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

func UpdateVoucherByID(updateVoucherByID voucherservice.UpdateVoucherByIDFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		var post struct {
			Voucher voucher.Voucher `json:"voucher"`
		}
		err = ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		post.Voucher.ID = id
		err = updateVoucherByID(ctx, post.Voucher)

		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)

	}
}
