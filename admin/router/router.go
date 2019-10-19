package router

import (
	"encoding/json"
	"fmt"
	"log"
	"mail/admin/context"
	"net/http"

	sjson "github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

type httpHandleFunc func(w http.ResponseWriter, r *http.Request) error
type middlewareFunc func(w http.ResponseWriter, r *http.Request) bool
type Router struct {
	mux *mux.Router
}

func NewRouter() *Router {
	h := mux.NewRouter()
	r := &Router{
		mux: h,
	}
	deployHandler := &DeployHandler{
		Ctx: context.ContextSingle,
	}
	nodeHandler := &NodeHandler{
		Ctx: context.ContextSingle,
	}
	containerHandler := &ContainerHandler{
		Ctx: context.ContextSingle,
	}
	adminHandler := &AdminHandler{
		Ctx: context.ContextSingle,
	}
	r.HandleFunc("/docker/deploy", deployHandler.GetDeployInfo).Methods(http.MethodGet)
	r.HandleFunc("/docker/deploy", deployHandler.SetDeployInfo).Methods(http.MethodPost)
	r.HandleFunc("/docker/node", deployHandler.GetMachineNodeInfo).Methods(http.MethodGet)        //已注册节点信息
	r.HandleFunc("/docker/node/register/{ip}", nodeHandler.RegisterNode).Methods(http.MethodPost) //注册节点信息
	r.HandleFunc("/docker/container/create", containerHandler.Create).Methods(http.MethodPost)    //创建docker实例
	r.HandleFunc("/docker/container/flush", containerHandler.Flush).Methods(http.MethodPost)      //清空
	r.HandleFunc("/docker/container/list", containerHandler.List).Methods(http.MethodPost)        //获取节点docker实例列表
	r.HandleFunc("/docker/container/start", containerHandler.Start).Methods(http.MethodGet)
	r.HandleFunc("/docker/container/stop", containerHandler.Stop).Methods(http.MethodGet)
	r.HandleFunc("/docker/container/delete", containerHandler.Delete).Methods(http.MethodDelete)
	//
	templatesDir := context.ContextSingle.Config.TemplatePath
	r.mux.PathPrefix("/vendors/").Handler(http.StripPrefix("/vendors/", http.FileServer(http.Dir(templatesDir+"/vendors"))))
	r.mux.PathPrefix("/build/").Handler(http.StripPrefix("/build/", http.FileServer(http.Dir(templatesDir+"/build"))))
	r.HandleAdminFunc("/admin/login", adminHandler.handleLoginGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/login", adminHandler.handleLoginPost).Methods(http.MethodPost)
	r.HandleAdminFunc("/admin/docker", adminHandler.handleDockerPageGet).Methods(http.MethodGet)
	//
	r.HandleAdminFunc("/admin/variable", adminHandler.handleVariablePageGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/variable", adminHandler.handleVariablePost).Methods(http.MethodPost)
	r.HandleAdminFunc("/admin/variable/list", adminHandler.handleVariableGet).Methods(http.MethodGet)
	//
	r.HandleAdminFunc("/admin/templates", adminHandler.handleTemplatesPageGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/templates", adminHandler.handleTemplatesPost).Methods(http.MethodPost)
	r.HandleAdminFunc("/admin/templates/list", adminHandler.handleTemplatesListGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/templates/{id}", adminHandler.handleTemplatesGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/templates/{id}", adminHandler.handleTemplatesDelete).Methods(http.MethodDelete)
	//
	r.HandleAdminFunc("/admin/receiver", adminHandler.handleReceiverGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/receiver", adminHandler.handleReceiverPost).Methods(http.MethodPost)
	r.HandleAdminFunc("/admin/receiver/list", adminHandler.handleReceiverListGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/receiver/{id}", adminHandler.handleReceiverDetailGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/receiver", adminHandler.handleReceiverDelete).Methods(http.MethodDelete)
	//
	r.HandleAdminFunc("/admin/sender", adminHandler.handleSenderGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/sender", adminHandler.handleSenderPost).Methods(http.MethodPost)
	r.HandleAdminFunc("/admin/sender", adminHandler.handleSenderDelete).Methods(http.MethodDelete)
	r.HandleAdminFunc("/admin/sender/list", adminHandler.handleSenderListGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/sender/{id}", adminHandler.handleSenderDetailGet).Methods(http.MethodGet)

	//
	r.HandleAdminFunc("/admin/setting", adminHandler.handleSettingGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/setting", adminHandler.handleSettingPost).Methods(http.MethodPost)
	//
	r.HandleAdminFunc("/admin/mail", adminHandler.handleMailPageGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/mail", adminHandler.handleMailPost).Methods(http.MethodPost)
	r.HandleAdminFunc("/admin/mail/log/list", adminHandler.handleMailLogListGet).Methods(http.MethodGet)
	r.HandleAdminFunc("/admin/mail/log/detail", adminHandler.handleMailLogDetail).Methods(http.MethodGet)
	return r
}

func (h *Router) HandleFunc(path string, action httpHandleFunc, middlewareFunc ...middlewareFunc) *mux.Route {
	return h.mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		for _, v := range middlewareFunc {
			ok := v(writer, request)
			if !ok {
				return
			}
		}
		err := action(writer, request)
		if err != nil {
			dd, _ := json.Marshal(request.URL)
			log.Println("URL:", string(dd))
			log.Printf("%s %s", request.RequestURI, err)
		}
	})
}

func (h *Router) HandleAdminFunc(path string, action httpHandleFunc, middlewareFunc ...middlewareFunc) *mux.Route {
	return h.mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		for _, v := range middlewareFunc {
			ok := v(writer, request)
			if !ok {
				http.Redirect(writer, request, "/admin/login", http.StatusFound)
				return
			}
		}
		err := action(writer, request)
		if err != nil {
			log.Println(err)
		}
	})
}

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func resError(w http.ResponseWriter, errMsg string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	log.Printf("err msg:%s, code: %d", errMsg, code)
	w.Write([]byte(fmt.Sprintf(`{"errMsg":"%s","success":false}`, errMsg)))
}

func resJson(w http.ResponseWriter, re interface{}) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	w.Header().Set("Access-Control-Allow-Headers", "*") //header的类型
	w.Header().Set("Access-Control-Allow-Methods", "*")
	sj := sjson.New()
	if re != nil {
		dd, _ := json.Marshal(re)
		subSj, _ := sjson.NewJson(dd)
		sj.Set("data", subSj)
	}
	sj.Set("success", true)
	data, err := sj.Encode()
	if err != nil {
		log.Println(err)
		resError(w, "未知错误", 500)
		return
	}
	w.Write(data)
}

func resTextError(w http.ResponseWriter, err string) {
	w.WriteHeader(400)
	w.Write([]byte(err))
}
