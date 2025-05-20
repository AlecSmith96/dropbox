package usecases

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type CreateFileRequestBody struct {
	Path        string `json:"path" binding:"required"`
	IsDirectory bool   `json:"isDirectory"`
	Data        []byte `json:"data"`
}

func NewCreateNewFile(fileCreator func(string, []byte, bool) error, destinationDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateFileRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("failed to bind json body for request", "err", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "a bad request error occurred",
			})
			return
		}

		filePathInDestinationDir := fmt.Sprintf("%s%s", destinationDir, request.Path)
		err = fileCreator(filePathInDestinationDir, request.Data, request.IsDirectory)
		if err != nil {
			slog.Error("creating new file", "err", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "an internal server error occurred",
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
