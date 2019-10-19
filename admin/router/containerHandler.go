package router

import (
	"log"
	"mail/admin/context"
	"mail/admin/sdk"
	"mail/mongo/model"
	"net/http"

	sjson "github.com/bitly/go-simplejson"
)

type ContainerHandler struct {
	Ctx *context.Context
}

func (h ContainerHandler) Create(w http.ResponseWriter, r *http.Request) error {
	sj, err := sjson.NewFromReader(r.Body)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	number := sj.Get("number").MustInt()
	host := sj.Get("host").MustString()
	log.Printf("number:%d,host:%s", number, host)
	dSdk := sdk.NewDeploySdk(host)
	err = dSdk.CreateContainer(number)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h ContainerHandler) Flush(w http.ResponseWriter, r *http.Request) error {
	sj, err := sjson.NewFromReader(r.Body)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	host := sj.Get("host").MustString()
	dSdk := sdk.NewDeploySdk(host)
	re, err := h.Ctx.MongoClient.ContainerRepository().FindByField("ip", host)
	log.Println("host:", host)
	log.Println("re:", re)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if len(re.List) == 0 {
		resJson(w, nil)
		return nil
	}
	for _, v := range re.List {
		err := dSdk.DeleteContainer(v.Id)
		if err != nil {
			resError(w, err.Error(), 400)
			return err
		}
	}
	resJson(w, nil)
	return nil
}

func (h ContainerHandler) Start(w http.ResponseWriter, r *http.Request) error {
	host := r.URL.Query().Get("host")
	id := r.URL.Query().Get("id")
	dSdk := sdk.NewDeploySdk(host)
	err := dSdk.StartContainer(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h ContainerHandler) Restart() {

}

func (h ContainerHandler) Update() { //pull最新镜像

}

func (h ContainerHandler) Stop(w http.ResponseWriter, r *http.Request) error {
	host := r.URL.Query().Get("host")
	id := r.URL.Query().Get("id")
	dSdk := sdk.NewDeploySdk(host)
	err := dSdk.StopContainer(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h ContainerHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	host := r.URL.Query().Get("host")
	id := r.URL.Query().Get("id")
	dSdk := sdk.NewDeploySdk(host)
	err := dSdk.DeleteContainer(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h ContainerHandler) List(w http.ResponseWriter, r *http.Request) error {
	sj, err := sjson.NewFromReader(r.Body)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	host := sj.Get("host").MustString()
	all := sj.Get("all").MustBool()
	if all == false && host == "" {
		resError(w, "缺少host参数", 404)
		return nil
	}
	list := make([]*model.Container, 0)
	if all {
		cInfo, err := h.Ctx.MongoClient.ContainerRepository().Find()
		if err != nil {
			resError(w, err.Error(), 404)
			return err
		}
		list = cInfo
	} else {
		cInfo, err := h.Ctx.MongoClient.ContainerRepository().FindByField("ip", host)
		if err != nil {
			resError(w, err.Error(), 404)
			return err
		}
		if cInfo == nil {
			resJson(w, nil)
			return nil
		}
		list = append(list, cInfo)
	}
	resJson(w, list)
	return nil
}
