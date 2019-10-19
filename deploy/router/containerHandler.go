package router

import (
	"errors"
	"fmt"
	"log"
	"mail/deploy/context"
	"net/http"
	"strconv"
	"strings"
	"sync"

	sjson "github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
)

var isProcess bool
var lock sync.Mutex

func init() {
	lock = sync.Mutex{}
}
func SetIsProcess(b bool) {
	lock.Lock()
	isProcess = b
	lock.Unlock()
}

func GetIsProcess() bool {
	return isProcess
}

type ContainerHandler struct {
	Ctx *context.Context
}

func (h *ContainerHandler) CreateContainer(w http.ResponseWriter, r *http.Request) error {
	var errMsg string
	if GetIsProcess() {
		errMsg = "有docker创建进程正在处理,请稍后重试"
		resError(w, errMsg, 400)
		return errors.New(errMsg)
	}
	numberStr := r.URL.Query().Get("number")
	if numberStr == "" {
		resError(w, "缺少number参数", 404)
		return nil
	}
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if number < 1 {
		number = 1
	}
	deploy, err := h.Ctx.MongoClient.DeployRepository().Find()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	cList, err := h.Ctx.SDK.Containers()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	if len(cList)+number > deploy.MaxInstance {
		errMsg = fmt.Sprintf("最多部署%d个实例,当前实例数:%d,请删掉一些已停止的实例后再重试", deploy.MaxInstance, len(cList))
		resError(w, errMsg, 400)
		return errors.New(errMsg)
	}
	go func() {
		SetIsProcess(true)
		sdk := h.Ctx.SDK
		for i := 0; i < number; i++ {
			var successInt, failedInt int
			name := uuid.NewV4().String()
			err := sdk.CreateContainer(name)
			if err != nil {
				failedInt++
				log.Printf("create %s failed:%s", name, err)
				continue
			}
			log.Printf("create %s success", name)
			successInt++
		}
		SetIsProcess(false)
	}()
	resJson(w, nil)
	return nil
}

func (h *ContainerHandler) DeleteContainer(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	sdk := h.Ctx.SDK
	err := sdk.DeleteContainer(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h *ContainerHandler) StartContainer(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	sdk := h.Ctx.SDK
	err := sdk.StartContainer(id)
	fmt.Println(err)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}

	resJson(w, nil)
	return nil
}

func (h *ContainerHandler) StopContainer(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	sdk := h.Ctx.SDK
	err := sdk.StopContainer(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h *ContainerHandler) UpdatePptpIp(w http.ResponseWriter, r *http.Request) error {
	ip := r.URL.Query().Get("ip")
	containerIp := strings.Split(r.RemoteAddr, ":")[0]
	container, err := h.Ctx.MongoClient.ContainerRepository().FindByField("ip", h.Ctx.Config.HostIp)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if container == nil {
		errMsg := fmt.Sprintf("no mechine found for ip:%s", h.Ctx.Config.HostIp)
		resError(w, errMsg, 400)
		return err
	}
	clist := container.List
	var isFound bool
	for i := 0; i < len(clist); i++ {
		if clist[i].Ip == containerIp {
			isFound = true
			clist[i].PptpIp = ip
			break
		}
	}
	if !isFound {
		errMsg := fmt.Sprintf("no container found for ip:%s", containerIp)
		resError(w, errMsg, 400)
		return err
	}
	container.List = clist
	err = h.Ctx.MongoClient.ContainerRepository().UpsertByField("ip", h.Ctx.Config.HostIp, container)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h *ContainerHandler) SendMail(w http.ResponseWriter, r *http.Request) error {
	sj, err := sjson.NewFromReader(r.Body)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	id := sj.Get("id").MustString()
	mailType := sj.Get("mailType").MustString()
	receiver := sj.Get("receiver").MustString()
	title := sj.Get("title").MustString()
	body := sj.Get("title").MustString()
	username := sj.Get("username").MustString()
	password := sj.Get("password").MustString()
	mode := sj.Get("mode").MustString()
	err = h.Ctx.SDK.SendMail(id, mailType, receiver, title, body, username, password, mode)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}
