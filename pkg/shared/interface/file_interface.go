package shared_interface

import (
	"github.com/ledaian41/go-cms-service/pkg/file/model"
	"mime/multipart"
)

type FileService interface {
	GetFileCachePath(path string) string
	SaveFile(file *multipart.FileHeader, uploadDir string) (*file_model.FileInfo, error)
}
