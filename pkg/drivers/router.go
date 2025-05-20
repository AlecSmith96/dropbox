package drivers

import (
	"github.com/AlecSmith96/dopbox/pkg/adapters"
	"github.com/AlecSmith96/dopbox/pkg/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewRouter is a function that reates a simple Gin router for the http server. It uses handler funcs to allow for
// dependency injection at the endpoint level, this restricts access for each endpoint to the exact dependencies they
// need.
func NewRouter(destinationDir string, fileWriter *adapters.FileWriter) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/health/live", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		v1.POST("/file", usecases.NewCreateNewFile(fileWriter.CreateFile, destinationDir))
		v1.DELETE("/file", usecases.NewDeleteFile(fileWriter.DeleteFile, destinationDir))
		v1.PATCH("/file", usecases.NewRenameFile(fileWriter.RenameFile, destinationDir))
		v1.PUT("/file", usecases.NewUpdateFileContents(fileWriter.UpdateFile, destinationDir))
	}

	return r
}
