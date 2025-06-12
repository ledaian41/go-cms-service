package nodeType_service

import (
	"errors"
	"fmt"
	"go-cms-service/config"
	"go-cms-service/pkg/shared/dto"
	"go-cms-service/pkg/shared/utils"
	"go-cms-service/pkg/valuetype"
	"log"
	"mime/multipart"
	"os"
	"sync"
)

type fileInfo struct {
	pid        string
	fileHeader *multipart.FileHeader
}

type fileResult struct {
	pid  string
	path string
	err  error
}

var (
	ErrFileTooLarge      = errors.New("file size exceeds maximum limit")
	ErrTotalSizeTooLarge = errors.New("total file size exceeds maximum limit")
)

func validateFileSize(fileHeader *multipart.FileHeader) error {
	fileSize := shared_utils.FileSize(fileHeader.Size)
	maxSize := shared_utils.FileSize(config.Env.MaxUploadFileSize * shared_utils.MB)
	if fileSize > maxSize {
		return fmt.Errorf("%w: file %s is %s, max allowed is %s",
			ErrFileTooLarge,
			fileHeader.Filename,
			fileSize.String(),
			maxSize.String())
	}
	return nil
}

func (s *NodeTypeService) PreprocessFile(nodeTypeDTO shared_dto.NodeTypeDTO, rawData map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	rawFiles := make(map[string]*multipart.FileHeader)
	for k, v := range rawData {
		fh, ok := v.(*multipart.FileHeader)
		if ok {
			rawFiles[k] = fh
			continue
		}
		result[k] = v
	}

	var totalSize shared_utils.FileSize
	filesChan := make(chan fileResult)

	filesToProcess := make([]fileInfo, 0)
	for _, pt := range nodeTypeDTO.PropertyTypes {
		valueType, err := valuetype.ParseValueType(pt.ValueType)
		if err != nil {
			continue
		}

		if valueType != valuetype.File {
			continue
		}

		fileHeader, exists := rawFiles[pt.PID]
		if !exists || fileHeader == nil {
			continue
		}

		if err := validateFileSize(fileHeader); err != nil {
			return nil, err
		}

		totalSize += shared_utils.FileSize(fileHeader.Size)
		maxTotalSize := shared_utils.FileSize(config.Env.MaxTotalUploadFileSize * shared_utils.MB)
		if totalSize > maxTotalSize {
			return nil, fmt.Errorf("%w: total size %s exceeds limit of %s",
				ErrTotalSizeTooLarge,
				totalSize.String(),
				maxTotalSize.String())
		}

		filesToProcess = append(filesToProcess, fileInfo{pid: pt.PID, fileHeader: fileHeader})
	}

	var wg sync.WaitGroup
	for _, fi := range filesToProcess {
		wg.Add(1)

		go func(pid string, fh *multipart.FileHeader) {
			defer wg.Done()

			fileInfo, err := s.fileService.SaveFile(fh, fmt.Sprintf("./%s/files/", config.Env.CachePath))
			if err != nil {
				filesChan <- fileResult{
					pid: pid,
					err: fmt.Errorf("failed to save file for property %s: %w", pid, err),
				}
				if (fileInfo != nil) && (fileInfo.SavedPath != "") {
					os.Remove(fileInfo.SavedPath)
				}
				return
			}

			filesChan <- fileResult{
				pid:  pid,
				path: fileInfo.SavedPath,
			}
		}(fi.pid, fi.fileHeader)
	}

	go func() {
		wg.Wait()
		close(filesChan)
	}()

	for fr := range filesChan {
		if fr.err != nil {
			log.Println(fr.err)
			continue
		}
		result[fr.pid] = fr.path
	}

	return result, nil
}
