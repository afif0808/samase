package filemanagerservice

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
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
