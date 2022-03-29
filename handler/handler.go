package handler

import (
	"SampleFileUpload/fileresponse"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
)

// Image handler demonstrates how to use `template.File` for responding with any Content-Type,
// in this example we respond with a PNG image
//func Image(c *krogo.Context) (interface{}, error) {
//	i := c.PathParam("id")
//	if i == "" {
//		return nil, errors.MissingParam{Param: []string{"id"}}
//	}
//	return template.Template{Directory: c.TemplateDir, File: i, Data: nil, Type: template.FILE}, nil
//}

func FetchImage(c *krogo.Context) (interface{}, error) {
	filePath := c.Request().FormValue("path")

	if filePath == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	filePath = c.TemplateDir + "/" + filePath
	return fileresponse.FetchImage(filePath)
}

func UploadFile(c *krogo.Context) (interface{}, error) {
	m := map[string]string{
		"Content-Type":                "text/html; charset=utf-8",
		"Access-Control-Allow-Origin": "*",
	}
	c.SetResponseHeader(m)
	r := c.Request()
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("myFile")
	if err != nil {
		return nil, errors.Error("cant get file")
	}
	category := r.FormValue("category")
	itemId := r.FormValue("id")
	path := c.TemplateDir
	return fileresponse.UploadFile(file, category, itemId, path)

}
