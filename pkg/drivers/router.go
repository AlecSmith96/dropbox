package drivers

import (
	"github.com/AlecSmith96/dopbox/pkg/adapters"
	"github.com/AlecSmith96/dopbox/pkg/usecases"
	"github.com/gin-gonic/gin"
)

func NewRouter(destinationDir string, fileWriter *adapters.FileWriter) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/file", usecases.NewCreateNewFile(fileWriter.CreateFile, destinationDir))
		v1.DELETE("/file", usecases.NewDeleteFile(fileWriter.DeleteFile, destinationDir))
	}

	return r
}
