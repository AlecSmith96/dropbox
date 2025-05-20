package usecases

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type DeleteFileRequestBody struct {
	Path string `json:"path" binding:"required"`
}

func NewDeleteFile(fileDeleter func(string) error, destinationDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateFileRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("failed to bind json body for request", "err", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "a bad request error occurred",
			})
		}

		filePathInDestinationDir := fmt.Sprintf("%s/%s", destinationDir, request.Path)
		err = fileDeleter(filePathInDestinationDir)
		if err != nil {
			slog.Error("deleting file", "err", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "an internal server error occurred",
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
