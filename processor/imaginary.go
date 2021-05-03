package processor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/agravelot/image_optimizer/config"
)

type ImaginaryProcessor struct {
	Url    string
	client http.Client
}

func isValidUrl(s string) error {

	if s == "" {
		return fmt.Errorf("url cannot be empty")
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unvalid imaginary scheme")
	}

	return nil
}

func NewImaginary(conf config.Config) (*ImaginaryProcessor, error) {

	err := isValidUrl(conf.Imaginary.Url)
	if err != nil {
		return nil, err
	}

	return &ImaginaryProcessor{
		client: http.Client{},
		Url:    conf.Imaginary.Url,
	}, nil
}

func (ip *ImaginaryProcessor) Optimize(media []byte, originalFormat string, targetFormat string, quality, width int) ([]byte, error) {

	println(width)

	u := fmt.Sprintf("%s/convert?type=webp&field=file", ip.Url)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	fileWriter, err := writer.CreateFormFile("file", "tmp.jpg")
	if err != nil {
		return nil, err
	}

	_, err = fileWriter.Write(media)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := ip.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
