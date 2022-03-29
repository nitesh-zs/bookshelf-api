package fileresponse

import (
	"bytes"
	"fmt"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo/template"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

func UploadFile(file multipart.File, category string, itemId string, path string) (interface{}, error) {

	if category == "user" {
		path += "/users"
	} else if category == "book" {
		path += "/books"
	}

	// Create a temporary file within our images directory that follows
	// a random naming pattern
	tempFile, err := ioutil.TempFile(path, itemId+"_*")
	if err != nil {
		fmt.Println(err)
		return "Internal Error", err
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return nil, err
	}
	// return file path which tells that we have successfully uploaded our file!
	filePSlice := strings.Split(tempFile.Name(), "/")
	fileP := filePSlice[len(filePSlice)-2] + "/" + filePSlice[len(filePSlice)-1]
	return fileP, nil
}

func FetchImage(filePath string) (interface{}, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, errors.Error("Cant open file")
	}
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Error("decode error")
	}
	b := new(bytes.Buffer)
	//err = png.Encode(b, i)
	err = jpeg.Encode(b, i, nil)
	if err != nil {
		return nil, errors.Error("encoding error")
	}
	return template.File{
		Content:     b.Bytes(),
		ContentType: "image/jpeg",
	}, nil
}

//func DeleteImage(filePath string) (interface{}, error) {
//	f, err := os.Open(filePath)
//	defer f.Close()
//	if err != nil {
//		return nil, errors.Error("Cant open file")
//	}
//	i, _, err := image.Decode(f)
//	if err != nil {
//		return nil, errors.Error("decode error")
//	}
//	b := new(bytes.Buffer)
//	//err = png.Encode(b, i)
//	err = jpeg.Encode(b, i, nil)
//	if err != nil {
//		return nil, errors.Error("encoding error")
//	}
//	return template.File{
//		Content:     b.Bytes(),
//		ContentType: "image/jpeg",
//	}, nil
//}
