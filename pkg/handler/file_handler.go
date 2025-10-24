package handler

import (
	"learn/fiber/pkg/service"
	"learn/fiber/utils"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

// @Summary		Upload File
// @Description	Upload File to S3
// @Tags			File
// @Accept			multipart/form-data
// @Produce		json
// @Param			X-Api-Key	header		string	true	"Api Key for access Public Endpoint"
// @Param			file		formData	file	true	"File"
// @Router			/file/upload [post]
func (h *FileHandler) UploadFileHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("file")

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	url, err := h.fileService.Upload(file)

	if err != nil {
		return err
	}

	return utils.SuccessResponse(c, int(fiber.StatusCreated), "Success Upload file", url)
}
