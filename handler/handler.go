package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/krogertechnology/krogo/pkg/krogo/template"
)

// Image handler demonstrates how to use `template.File` for responding with any Content-Type,
// in this example we respond with a PNG image
func Image(c *krogo.Context) (interface{}, error) {
	i := c.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	return template.Template{Directory: c.TemplateDir, File: i, Data: nil, Type: template.FILE}, nil
}

func FetchImage(c *krogo.Context) (interface{}, error) {
	fileName := c.PathParam("id")
	if fileName == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}
	filePath := c.TemplateDir + "/" + fileName
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return nil, errors.Error("Cant open file")
	}

	i, formatName, err := image.Decode(f)
	//fmt.Println(formatName)
	if err != nil {
		return nil, errors.Error("Cant decode file. Format : " + formatName)
	}

	b := new(bytes.Buffer)
	png.Encode(b, i)

	return template.File{
		Content:     b.Bytes(),
		ContentType: "image/png",
	}, nil
}

func UploadFile(c *krogo.Context) (interface{}, error) {
	m := map[string]string{
		"Content-Type":                "text/html; charset=utf-8",
		"Access-Control-Allow-Origin": "*",
	}

	c.SetResponseHeader(m)

	r := c.Request()
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return "Internal Error", err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our templates directory that follows
	// a random naming pattern
	tempFile, err := ioutil.TempFile("templates", "*.png")
	if err != nil {
		fmt.Println(err)
		return "Internal Error", err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "Internal Error", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return file path which tells that we have successfully uploaded our file!
	return tempFile.Name(), nil
}
