package sender

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	sjson "github.com/bitly/go-simplejson"
)

type Sender163 struct{}

func (s *Sender163) Send(username, password, receiver, title, body, mode string) error {
	urlStr := fmt.Sprintf("http://127.0.0.1:8090/mail/send")
	uv := url.Values{}
	uv.Set("username", username)
	uv.Set("password", password)
	uv.Set("receiver", receiver)
	uv.Set("mailtitle", title)
	uv.Set("mailcontent", body)
	log.Printf("%s %s", urlStr, uv.Encode())
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(uv.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res *http.Response
	log.Println("mode:", mode)
	start := time.Now().Unix()
	if strings.ToUpper(mode) == "IP" {
		res, err = http.DefaultClient.Do(req)
	} else {
		err = exec.Command("/pptp_start.sh").Run()
		if err != nil {
			log.Println("start pptp failed:", err)
			return err
		}
		log.Println("start pptp success")
		time.Sleep(time.Second * 2) //等待2秒是为了保证能正常启动pptp client,这里可以优化
		res, err = http.DefaultClient.Do(req)
		stopErr := exec.Command("/pptp_stop.sh").Run()
		if stopErr != nil {
			log.Println("stop pptp failed:", stopErr)
		}
		log.Println("stop pptp success")
	}
	var end int64
	if strings.ToUpper(mode) == "VPN" {
		end = time.Now().Unix() - start - 2
	}
	log.Printf("java发信耗时%d秒", end-start)
	log.Println("res err:", err)
	if err != nil {
		return err
	}
	sj, err := sjson.NewFromReader(res.Body)
	if err != nil {
		return err
	}
	dd, _ := sj.Encode()
	log.Println("java res data:", string(dd))
	if sj.GetPath("resphead", "success").MustBool() {
		return nil
	}
	return errors.New(sj.GetPath("resphead", "msg").MustString())
}
