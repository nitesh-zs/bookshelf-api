package main

import (
	"SampleFileUpload/handler"
	"os"

	"github.com/krogertechnology/krogo/pkg/krogo"
)

func main() {
	// Create the application object
	k := krogo.New()
	k.Server.ValidateHeaders = false
	rootPath, _ := os.Getwd()

	// overriding default template location.
	k.TemplateDir = rootPath + "/templates"
	k.GET("/image/{id}", handler.Image)
	k.POST("/upload", handler.UploadFile)
	k.Start()
}
