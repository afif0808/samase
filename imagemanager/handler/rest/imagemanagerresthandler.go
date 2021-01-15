package imagemanagerresthandler

import (
	"net/http"
	filemanagerservice "samase/filemanager/service"
	imagemanagerservice "samase/imagemanager/service"

	"github.com/labstack/echo"
)

func InjectImageManagerRESTHandler(ee *echo.Echo) {
	saveFile := filemanagerservice.SaveFile("C:/Users/bayur/go/src/samase/assets/images/")
	// saveFile := filemanagerservice.SaveFile("/root/go/src/samase/assets/images/")

	ee.POST("/upload/image", UploadImage(saveFile))
	listFolderFiles := filemanagerservice.ListFolderFiles()
	listFolderImages := imagemanagerservice.ListFolderImages(listFolderFiles)
	ee.GET("/images/list", GetImages(listFolderImages, "C:/Users/bayur/go/src/samase/assets/images"))

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

func GetImages(
	listFolderImages imagemanagerservice.ListFolderImagesFunc,
	imagesFolder string,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		images, err := listFolderImages(imagesFolder)
		type image struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		resp := make([]image, len(images))
		for i, img := range images {
			resp[i] = image{
				ID:   int64(i),
				Name: img,
			}
		}

		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, resp)
	}
}
