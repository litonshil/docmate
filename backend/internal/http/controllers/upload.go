package controllers

import (
	"docmate/response"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (u *UploadController) UploadFile(c echo.Context) error {
	// 1. Get the file from form
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "No file uploaded")
	}

	src, err := file.Open()
	if err != nil {
		return response.InternalServerError(c, "Could not open file")
	}
	defer src.Close()

	// 2. Ensure uploads directory exists
	uploadDir := "uploads"
	if _, err := os.Stat("/project"); err == nil {
		uploadDir = "/project/uploads"
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		cwd, _ := os.Getwd()
		return response.InternalServerError(c, fmt.Sprintf("Could not create uploads directory (CWD: %s, Dir: %s): %v", cwd, uploadDir, err))
	}

	// 3. Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	dstPath := filepath.Join(uploadDir, filename)

	// 4. Create the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return response.InternalServerError(c, fmt.Sprintf("Could not create file on server: %v", err))
	}
	defer dst.Close()

	// 5. Copy the uploaded file to the destination
	if _, err = io.Copy(dst, src); err != nil {
		return response.InternalServerError(c, fmt.Sprintf("Could not save file: %v", err))
	}

	// 6. Return the accessible URL
	// Note: We'll serve /uploads/ statically in routes.go
	scheme := "http"
	if c.IsTLS() {
		scheme = "https"
	}
	fileURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, c.Request().Host, filename)

	return response.Success(c, "File uploaded successfully", map[string]string{
		"url":      fileURL,
		"filename": filename,
	})
}
