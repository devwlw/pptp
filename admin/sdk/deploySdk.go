package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	sjson "github.com/bitly/go-simplejson"
)

type DeploySdk struct {
	endpoint string
}

func NewDeploySdk(host string) *DeploySdk {
	return &DeploySdk{
		endpoint: fmt.Sprintf("http://%s:%s", host, "9100"),
	}
}

func (s *DeploySdk) CreateContainer(number int) error {
	urlStr := fmt.Sprintf("%s/docker/container/create?number=%d", s.endpoint, number)
	_, err := s.do(http.MethodPost, urlStr, nil)
	return err
}

func (s *DeploySdk) StopContainer(id string) error {
	urlStr := fmt.Sprintf("%s/docker/container/stop?id=%s", s.endpoint, id)
	_, err := s.do(http.MethodGet, urlStr, nil)
	return err
}
func (s *DeploySdk) StartContainer(id string) error {
	urlStr := fmt.Sprintf("%s/docker/container/start?id=%s", s.endpoint, id)
	_, err := s.do(http.MethodGet, urlStr, nil)
	return err
}
func (s *DeploySdk) DeleteContainer(id string) error {
	urlStr := fmt.Sprintf("%s/docker/container/delete?id=%s", s.endpoint, id)
	_, err := s.do(http.MethodDelete, urlStr, nil)
	return err
}

func (s *DeploySdk) SendMail(id, mailType, receiver, title, body, username, password, mode string) error {
	sj := sjson.New()
	sj.Set("id", id)
	sj.Set("mailType", mailType)
	sj.Set("receiver", receiver)
	sj.Set("title", title)
	sj.Set("body", body)
	sj.Set("username", username)
	sj.Set("password", password)
	sj.Set("mode", mode)
	log.Printf("send mail:%s,%s,%s", id, s.endpoint, receiver)
	urlStr := fmt.Sprintf("%s/docker/mail/send", s.endpoint)
	dd, _ := sj.Encode()
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(dd))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	sj, err = sjson.NewFromReader(res.Body)
	if err != nil {
		return err
	}
	dd, _ = sj.Encode()
	log.Printf("mail id:%s,receiver:%s , res:%s", id, receiver, string(dd))
	if sj.Get("success").MustBool() {
		return nil
	}
	return errors.New(sj.Get("errMsg").MustString())
}

func (s *DeploySdk) do(method, urlStr string, body io.Reader) (*sjson.Json, error) {
	log.Printf("%s %s", method, urlStr)
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	sj, err := sjson.NewFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	if !sj.Get("success").MustBool() {
		return nil, errors.New(sj.Get("errMsg").MustString())
	}
	return sj, nil
}
