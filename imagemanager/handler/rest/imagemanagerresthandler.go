package imagemanagerresthandler

import (
	"net/http"
	filemanagerservice "samase/filemanager/service"

	"github.com/labstack/echo"
)

func InjectImageManagerRESTHandler(ee *echo.Echo) {
	saveFile := filemanagerservice.SaveFile("/root/go/src/samase/assets/images/")
	ee.POST("/upload/image", UploadImage(saveFile))
}

func UploadImage(saveFile filemanagerservice.SaveFileFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		imgHeader, err := ectx.FormFile("image")
		fileName := imgHeader.Filename
		f, err := imgHeader.Open()
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = saveFile(ctx, f, fileName)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}
