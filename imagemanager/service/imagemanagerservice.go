package imagemanagerservice

import (
	"path/filepath"
	filemanagerservice "samase/filemanager/service"
)

func ListFolderImages(listFolderFiles filemanagerservice.ListFolderFilesFunc) ListFolderImagesFunc {
	return func(folder string) ([]string, error) {
		files, err := listFolderFiles(folder)
		if err != nil {
			return nil, err
		}
		images := []string{}
		for _, f := range files {
			if filepath.Ext(f) == ".jpg" || filepath.Ext(f) == ".png" || filepath.Ext(f) == ".jpeg" {
				images = append(images, f)
			}
		}
		return images, nil
	}
}
