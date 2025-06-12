package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ledaian41/go-cms-service/pkg/shared/interface"
	"net/http"
)

type FileHandler struct {
	fileService shared_interface.FileService
}

func NewFileHandler(fileService shared_interface.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h FileHandler) ReadFile(context *gin.Context) {
	filepath := h.fileService.GetFileCachePath(context.Param("path")[1:])
	if len(filepath) == 0 {
		context.String(http.StatusNotFound, "sorry, resource not found!!")
		return
	}
	context.File(filepath)
}
