package filemanagerservice

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFile(basename string) SaveFileFunc {
	return func(ctx context.Context, f multipart.File, filename string) error {
		dest, err := os.Create(basename + filename)
		if err != nil {
			log.Println(err)
			return err

		}

		_, err = io.Copy(dest, f)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}

func ListFolderFiles() ListFolderFilesFunc {
	return func(folder string) ([]string, error) {
		files := []string{}
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				log.Println(err)
				return err
			}
			if !info.IsDir() {
				files = append(files, info.Name())
			}

			return nil
		})
		if err != nil {
			log.Println(err)
			return nil, err

		}
		return files, nil
	}
}
