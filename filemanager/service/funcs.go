package filemanagerservice

import (
	"context"
	"mime/multipart"
)

type SaveFileFunc func(ctx context.Context, f multipart.File, filename string) error
type ListFolderFilesFunc func(folder string) ([]string , error)