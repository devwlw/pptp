package router

import (
	"errors"
	"mail/admin/context"
	"mail/mongo/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type NodeHandler struct {
	Ctx *context.Context
}

/*func (d NodeHandler) DeleteNode(w http.ResponseWriter, r *http.Request) error {

}*/

//获取节点的信息
func (d NodeHandler) GetNodeDockerInfo(w http.ResponseWriter, r *http.Request) error {
	node, err := d.Ctx.MongoClient.NodeRepository().Find()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, node)
	return nil
}

func (d NodeHandler) RegisterNode(w http.ResponseWriter, r *http.Request) error {
	ip := mux.Vars(r)["ip"]
	if ip == "" {
		errMsg := "缺少ip参数"
		resError(w, errMsg, 404)
		return errors.New(errMsg)
	}
	aux := r.URL.Query().Get("aux")
	ipRange := r.URL.Query().Get("ipRange")
	subnet := r.URL.Query().Get("subnet")
	nic := r.URL.Query().Get("nic")
	deploy, err := d.Ctx.MongoClient.DeployRepository().Find()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	if deploy == nil {
		resError(w, "deploy信息为空", 400)
		return errors.New("deploy信息为空")
	}
	node, err := d.Ctx.MongoClient.NodeRepository().Find()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	if node == nil {
		node = new(model.Node)
	}
	nodes := node.Nodes
	if len(nodes) >= deploy.MaxMachine {
		errMsg := "超过最大节点数"
		resError(w, errMsg, 400)
		return errors.New(errMsg)
	}
	for _, v := range nodes {
		if v.Ip == ip {
			errMsg := "此节点已存在"
			resError(w, errMsg, 400)
			return errors.New(errMsg)
		}
	}

	nodes = append(nodes, model.NodeInfo{
		Ip:          ip,
		Aux:         aux,
		IpRange:     ipRange,
		Subnet:      subnet,
		Nic:         nic,
		CreatedTime: time.Now().Unix(),
	})
	node.Nodes = nodes
	err = d.Ctx.MongoClient.NodeRepository().Upsert(node)
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	resJson(w, nil)
	return nil
}
