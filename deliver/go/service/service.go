package service

import (
	"fmt"
	sender2 "mail/deliver/go/service/sender"
)

type Sender interface {
	Send(username, password, receiver, title, body, mode string) error
}

var serviceList map[string]Sender

func init() {
	serviceList = make(map[string]Sender)
	serviceList["163"] = &sender2.Sender163{}
}

func GetSender(t string) (Sender, error) {
	s := serviceList[t]
	if s == nil {
		return nil, fmt.Errorf("no sender found for:%s", t)
	}
	return s, nil
}
