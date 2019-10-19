package sender

import (
	"errors"
	"fmt"
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
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(uv.Encode()))
	if err != nil {
		return err
	}
	var res *http.Response
	if strings.ToLower(mode) == "IP" {
		res, err = http.DefaultClient.Do(req)
	} else {
		err = exec.Command("/pptp_start.sh").Run()
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 3) //等待三秒是为了保证能正常启动pptp client,这里可以优化
		res, err = http.DefaultClient.Do(req)
		exec.Command("/pptp_stop.sh").Run()
	}
	if err != nil {
		return err
	}
	sj, err := sjson.NewFromReader(res.Body)
	if err != nil {
		return err
	}
	if sj.GetPath("resphead", "success").MustBool() {
		return nil
	}
	return errors.New(sj.GetPath("resphead", "msg").MustString())
}
