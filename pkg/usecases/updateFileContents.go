package usecases

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type UpdateFileContentsRequestBody struct {
	Path string `json:"path" binding:"required"`
	Data []byte `json:"data"`
}

func NewUpdateFileContents(fileUpdater func(string, []byte) error, destinationDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request UpdateFileContentsRequestBody
		err := c.ShouldBindJSON(&request)
		if err != nil {
			slog.Warn("failed to bind json body for request", "err", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "a bad request error occurred",
			})
			return
		}

		filePathInDestinationDir := fmt.Sprintf("%s%s", destinationDir, request.Path)
		err = fileUpdater(filePathInDestinationDir, request.Data)
		if err != nil {
			slog.Error("updating file contents", "err", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "an internal server error occurred",
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
