package file_service

import (
	"fmt"
	"github.com/ledaian41/go-cms-service/config"
	"github.com/ledaian41/go-cms-service/pkg/file/model"
	"github.com/ledaian41/go-cms-service/pkg/file/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (s FileService) GetFileCachePath(path string) string {
	fullPath := filepath.Join(config.Env.CachePath, "files", path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return ""
	}
	return fullPath
}

func (s FileService) SaveFile(file *multipart.FileHeader, uploadDir string) (*file_model.FileInfo, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir error: %v", err)
	}

	filepath, err := file_utils.GenerateUploadPath(file.Filename)
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("source file error: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("create new file error: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("copy file error: %v", err)
	}

	return &file_model.FileInfo{
		OriginalName: file.Filename,
		SavedPath:    filepath,
		Size:         file.Size,
		ContentType:  file.Header.Get("Content-Type"),
	}, nil
}
