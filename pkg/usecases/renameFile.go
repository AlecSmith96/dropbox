package usecases

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type RenameFileRequestBody struct {
	Path         string `json:"path" binding:"required"`
	PreviousPath string `json:"previousPath"`
}

func NewRenameFile(fileRenamingFunc func(string, string) error, destinationDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RenameFileRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("failed to bind json body for request", "err", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "a bad request error occurred",
			})
			return
		}

		oldFilePathInDestinationDir := fmt.Sprintf("%s%s", destinationDir, request.PreviousPath)
		newFilePathInDestinationDir := fmt.Sprintf("%s%s", destinationDir, request.Path)

		err = fileRenamingFunc(oldFilePathInDestinationDir, newFilePathInDestinationDir)
		if err != nil {
			slog.Error("renaming file", "err", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "an internal server error occurred",
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
