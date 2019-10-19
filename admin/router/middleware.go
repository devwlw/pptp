package router

import (
	"fmt"
	"net/http"

	"log"
)

func handleCors(writer http.ResponseWriter, req *http.Request) bool {
	writer.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	writer.Header().Set("Access-Control-Allow-Headers", "*") //header的类型
	writer.Header().Set("Access-Control-Allow-Methods", "*")
	//writer.Header().Set("content-type", "application/json") //返回数据格式是json
	if req.Method == http.MethodOptions {
		return false
	}
	return true
}

func handleCheckCookie(res http.ResponseWriter, req *http.Request) bool {
	sidCookie, err := req.Cookie("SID")
	if err != nil {
		return false
		log.Println(err)
	}
	sid := sidCookie.Value
	svCookie, err := req.Cookie("SV")
	if err != nil {
		return false
		log.Println(err)
	}
	sv := svCookie.Value
	fmt.Println("sid,sv", sid, sv)
	session := SessionMapInstance.Get(sid)
	if session == nil {
		return false
	}
	if !session.IsValid(sv) {
		return false
	}
	if session.IsExpired() {
		SessionMapInstance.Del(sid)
		return false
	}
	return true
}
