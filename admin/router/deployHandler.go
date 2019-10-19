package router

import (
	"encoding/json"
	"errors"
	"mail/admin/context"
	"mail/mongo/model"
	"net/http"
)

type DeployHandler struct {
	Ctx *context.Context
}

func (d DeployHandler) GetDeployInfo(w http.ResponseWriter, r *http.Request) error {
	deploy, err := d.Ctx.MongoClient.DeployRepository().Find()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	if deploy == nil {
		resError(w, "deploy信息为空", 400)
		return errors.New("deploy信息为空")
	}
	resJson(w, deploy)
	return nil
}

func (d DeployHandler) SetDeployInfo(w http.ResponseWriter, r *http.Request) error {
	deploy := new(model.Deploy)
	err := json.NewDecoder(r.Body).Decode(deploy)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	err = d.Ctx.MongoClient.DeployRepository().Upsert(deploy)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (d DeployHandler) GetMachineNodeInfo(w http.ResponseWriter, r *http.Request) error {
	node, err := d.Ctx.MongoClient.NodeRepository().Find()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	resJson(w, node)
	return nil
}

func (d DeployHandler) UpdateImage(w http.ResponseWriter, r *http.Request) error {
	return nil
}
