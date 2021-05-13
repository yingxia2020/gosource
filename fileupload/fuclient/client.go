/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, paramName, path string) (*fasthttp.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	//for key, val := range params {
	//	_ = writer.WriteField(key, val)
	//}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethodBytes([]byte("POST"))
	req.SetBody(body.Bytes())

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func main() {
	client := &fasthttp.Client{}
	request, err := newfileUploadRequest("http://localhost:8073/upload", "uploadfile", "./simple.tiff")
	if err != nil {
		log.Fatal(err)
	}
	response := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetConnectionClose()

	if err := client.Do(request, response); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(response.Header.Header()))
		fmt.Println(string(response.Body()))
	}
}
