package router

import (
	"encoding/json"
	"fmt"
	"log"
	"mail/deploy/context"
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
	containerHandler := &ContainerHandler{
		Ctx: context.ContextSingle,
	}
	r.HandleFunc("/docker/container/create", containerHandler.CreateContainer).Methods(http.MethodPost)
	r.HandleFunc("/docker/container/delete", containerHandler.DeleteContainer).Methods(http.MethodDelete)
	r.HandleFunc("/docker/container/start", containerHandler.StartContainer).Methods(http.MethodGet)
	r.HandleFunc("/docker/container/stop", containerHandler.StartContainer).Methods(http.MethodGet)
	r.HandleFunc("/docker/container/pptpip", containerHandler.UpdatePptpIp).Methods(http.MethodPost)
	r.HandleFunc("/docker/mail/send", containerHandler.SendMail).Methods(http.MethodPost)
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
