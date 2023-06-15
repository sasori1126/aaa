package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/pkg/entities"
)

type SendInBlueTempEmail struct {
	Sender     EmailContact      `json:"sender"`
	To         []EmailContact    `json:"to"`
	ReplyTo    EmailContact      `json:"replyTo"`
	TemplateID int64             `json:"templateId"`
	Params     map[string]string `json:"params"`
}

type EmailContact struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Headers struct {
	XMailinCustom string `json:"X-Mailin-custom"`
	Charset       string `json:"charset"`
}

type HttpService struct {
	BaseUrl string
	Client  *http.Client
}

func NewHttpService(baseUrl string) *HttpService {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	return &HttpService{BaseUrl: baseUrl, Client: httpClient}
}

func (h *HttpService) Post(path string, data []byte, headers ...entities.HttpHeader) {
	log.Println(headers)
}

func (h *HttpService) SiBTempEmail(tempData *SendInBlueTempEmail) error {
	body, err := json.Marshal(tempData)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", h.BaseUrl+"smtp/email", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/json")
	request.Header.Set("accept", "application/json")
	request.Header.Set("api-key", configs.AppConfig.SendInBlueKey)
	resp, err := h.httpRequest(request)
	if err != nil {
		return err
	}
	log.Println(string(resp))

	var res struct {
		Message string
		Code    string
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return err
	}

	if res.Code == "unauthorized" {
		return errors.New("unauthorized: " + res.Message)
	}

	return nil
}

func (h *HttpService) httpRequest(req *http.Request) ([]byte, error) {
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
