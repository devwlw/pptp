package router

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mail/admin/sdk"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"

	sjson "github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"

	//sjson "github.com/bitly/go-simplejson"
	"html/template"
	"mail/admin/context"
	"mail/admin/tool"
	"mail/mongo/model"
	"net/http"
	"time"
)

type AdminHandler struct {
	Ctx *context.Context
}

func (h AdminHandler) ParseTemplate(path string) (*template.Template, error) {
	tPath := h.Ctx.Config.TemplatePath
	return template.ParseFiles(tPath + path)
}

func (h AdminHandler) handleLoginGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/login.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	return t.Execute(w, "")
}

func (h AdminHandler) handleLoginPost(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	fmt.Println(r.Form)
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	ar := h.Ctx.MongoClient.AdminRepository()
	admin, err := ar.Get(username)
	if err != nil {
		resTextError(w, "登录失败")
		return err
	}
	if admin == nil {
		resTextError(w, "无此用户")
		return err
	}
	if !admin.CheckPassword(password, admin.PassWord) {
		resTextError(w, "密码错误")
		return err
	}
	sid := tool.UUID()
	sv := tool.NonceStrN(32)
	expires := time.Now().Add(time.Hour * 6)
	cookie := http.Cookie{
		Name:    "SID",
		Value:   sid,
		Path:    "/",
		Expires: expires,
	}
	cookie1 := http.Cookie{
		Name:    "SV",
		Value:   sv,
		Path:    "/",
		Expires: expires,
	}
	s := &Session{
		Id:      cookie.Value,
		Value:   sv,
		Expired: expires,
	}
	SessionMapInstance.Add(s)
	http.SetCookie(w, &cookie)
	http.SetCookie(w, &cookie1)
	http.SetCookie(w, &http.Cookie{
		Name:    "avatarImg",
		Value:   admin.Avatar,
		Path:    "/",
		Expires: expires,
	})
	http.Redirect(w, r, "/admin/setting", http.StatusFound)
	return nil
}

func (h AdminHandler) handleDockerPageGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/docker.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	return t.Execute(w, "")
}

func (h AdminHandler) handleVariablePageGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/variable.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	m := make(map[string]string)
	v, err := h.Ctx.MongoClient.VariableRepository().Find()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if v != nil {
		m = v.Data
	}
	return t.Execute(w, m)
}

func (h AdminHandler) handleVariablePost(w http.ResponseWriter, r *http.Request) error {
	m := make(map[string]string)
	err := r.ParseMultipartForm(1024 * 1024 * 100)
	if err != nil {
		resTextError(w, "最多接收100M数据")
		return err
	}
	form := r.MultipartForm.Value
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("variable%d", i+1)
		m[name] = h.parseFormText(form, name)
	}
	variable := &model.Variable{
		Data: m,
	}
	err = h.Ctx.MongoClient.VariableRepository().Upsert(variable)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	http.Redirect(w, r, "/admin/variable", http.StatusFound)
	return nil
}

func (h AdminHandler) handleVariableGet(w http.ResponseWriter, r *http.Request) error {
	v, err := h.Ctx.MongoClient.VariableRepository().Find()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, v)
	return nil
}

func (h AdminHandler) handleTemplatesPageGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/templates.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	return t.Execute(w, "")
}

func (h AdminHandler) handleTemplatesListGet(w http.ResponseWriter, r *http.Request) error {
	tr := h.Ctx.MongoClient.TemplatesRepository()
	list, err := tr.Get()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	type t struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		CreatedTime string `json:"createdTime"`
	}
	arr := make([]t, 0)
	for _, v := range list {
		arr = append(arr, t{
			Id:          v.Id,
			Name:        v.Name,
			CreatedTime: time.Unix(v.CreatedTime, 0).Format("2006-01-02 15:04:05"),
		})
	}
	resJson(w, arr)
	return nil
}

func (h AdminHandler) handleTemplatesGet(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	if id == "" {
		resError(w, "缺少id参数", 404)
		return nil
	}
	tm, err := h.Ctx.MongoClient.TemplatesRepository().GetById(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if tm == nil {
		resError(w, "id错误", 404)
		return nil
	}
	t, err := h.ParseTemplate("/templates-detail.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	m := make(map[string]string)
	m["name"] = tm.Name
	m["data"] = tm.Data
	return t.Execute(w, m)
}

func (h AdminHandler) handleTemplatesDelete(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	if id == "" {
		resError(w, "缺少id参数", 404)
		return nil
	}
	err := h.Ctx.MongoClient.TemplatesRepository().DeleteById(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h AdminHandler) handleTemplatesPost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 * 1024 * 100)
	if err != nil {
		resTextError(w, "最多接收100M数据")
		return err
	}
	form := r.MultipartForm.Value
	val := h.parseFormText(form, "data")
	name := h.parseFormText(form, "name")
	if name == "" {
		resTextError(w, "name不能为空")
		return nil
	}
	if val == "" {
		resTextError(w, "data不能为空")
		return nil
	}
	tr := h.Ctx.MongoClient.TemplatesRepository()
	old, _ := tr.GetByName(name)
	if old != nil {
		resTextError(w, fmt.Sprintf("模板%s已存在", old.Name))
		return nil
	}
	t := &model.Templates{
		Id:          tool.UUID(),
		Name:        name,
		Data:        val,
		CreatedTime: time.Now().Unix(),
	}
	err = tr.Upsert(t)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h AdminHandler) handleReceiverGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/receiver.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	return t.Execute(w, "")
}

func (h AdminHandler) handleReceiverPost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 * 1024 * 100)
	if err != nil {
		resTextError(w, "最多接收100M数据")
		return err
	}
	f, err := h.parseFile(r, "receiver")
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	rr := h.Ctx.MongoClient.ReceiverRepository()
	err = rr.Flush()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	data = bytes.Replace(data, []byte("\r"), []byte{}, -1)
	for _, v := range bytes.Split(data, []byte("\n")) {
		err = rr.Upsert(&model.ReceiverInfo{
			Id:    uuid.NewV4().String(),
			Email: string(v),
		})
		if err != nil {
			resError(w, err.Error(), 400)
			return err
		}
	}
	resJson(w, nil)
	return nil
}

func (h AdminHandler) handleReceiverListGet(w http.ResponseWriter, r *http.Request) error {
	rr := h.Ctx.MongoClient.ReceiverRepository()
	draw := r.URL.Query().Get("draw")
	length := r.URL.Query().Get("length")
	value := r.URL.Query().Get("search[value]")
	_ = value
	start := r.URL.Query().Get("start")
	startNum, _ := strconv.Atoi(start)
	lengthNum, _ := strconv.Atoi(length)
	re, total, err := rr.Get(startNum, lengthNum)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	sj := sjson.New()
	sj.Set("draw", draw)
	sj.Set("recordsTotal", total)
	sj.Set("recordsFiltered", total)
	sj.Set("data", re)
	dd, _ := sj.Encode()
	w.Write(dd)
	return nil
}

func (h AdminHandler) handleReceiverDetailGet(w http.ResponseWriter, r *http.Request) error {
	/*	id := mux.Vars(r)["id"]
		rr := h.Ctx.MongoClient.ReceiverRepository()
		receiver, err := rr.GetById(id)
		if err != nil {
			resError(w, err.Error(), 400)
			return err
		}
		if receiver == nil {
			resTextError(w, "无此id")
			return nil
		}
		w.WriteHeader(200)
		w.Write([]byte(strings.Join(receiver.Data, "\n")))*/
	return nil
}

func (h AdminHandler) handleReceiverDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	rr := h.Ctx.MongoClient.ReceiverRepository()
	err := rr.DeleteById(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

//
func (h AdminHandler) handleSenderGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/sender.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	return t.Execute(w, "")
}

//
func (h AdminHandler) handleSenderPost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 * 1024 * 100)
	if err != nil {
		resTextError(w, "最多接收100M数据")
		return err
	}
	//form := r.MultipartForm.Value
	f, err := h.parseFile(r, "sender")
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	sr := h.Ctx.MongoClient.SenderRepository()
	err = sr.Flush()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	str := string(data)
	str = strings.Replace(str, "\r", "", -1)
	for _, v := range strings.Split(str, "\n") {
		tmp := strings.Replace(v, " ", "", 1)
		if tmp == "" {
			continue
		}
		arr := strings.Split(v, " ")
		if len(arr) != 2 {
			resTextError(w, "文件内容不合法")
			return nil
		}
		err := sr.Upsert(&model.SendInfo{
			Id:       uuid.NewV4().String(),
			Email:    arr[0],
			Password: arr[1],
		})
		if err != nil {
			resError(w, err.Error(), 400)
			return err
		}
	}
	resJson(w, nil)
	return nil
}

func (h AdminHandler) handleSenderListGet(w http.ResponseWriter, r *http.Request) error {
	sr := h.Ctx.MongoClient.SenderRepository()

	draw := r.URL.Query().Get("draw")
	length := r.URL.Query().Get("length")
	value := r.URL.Query().Get("search[value]")
	_ = value
	start := r.URL.Query().Get("start")
	startNum, _ := strconv.Atoi(start)
	lengthNum, _ := strconv.Atoi(length)
	re, total, err := sr.Get(startNum, lengthNum)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	sj := sjson.New()
	sj.Set("draw", draw)
	sj.Set("recordsTotal", total)
	sj.Set("recordsFiltered", total)
	sj.Set("data", re)
	dd, _ := sj.Encode()
	w.Write(dd)
	return nil
}

/*func (h AdminHandler) handleSenderDetailGet(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	sr := h.Ctx.MongoClient.SenderRepository()
	sender, err := sr.GetById(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if sender == nil {
		resTextError(w, "无此id")
		return nil
	}
	w.WriteHeader(200)
	arr := make([]string, 0)
	for _, v := range sender.Data {
		arr = append(arr, v.Email+" "+v.Password)
	}
	w.Write([]byte(strings.Join(arr, "\n")))
	return nil
}*/

func (h AdminHandler) handleSenderDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")
	sr := h.Ctx.MongoClient.SenderRepository()
	err := sr.DeleteById(id)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	resJson(w, nil)
	return nil
}

func (h AdminHandler) handleSettingGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/setting.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	de, err := h.Ctx.MongoClient.DeployRepository().Find()
	if err != nil {
		resError(w, err.Error(), 500)
		return err
	}
	var maxMachine, maxInstance int
	m := make(map[string]int)
	if de != nil {
		maxMachine = de.MaxMachine
		maxInstance = de.MaxInstance
	}
	m["maxMachine"] = maxMachine
	m["maxInstance"] = maxInstance
	return t.Execute(w, m)
}

//todo 新部署admin时,如果没设置deploy的值,则在注册还有不部署docker时应该报错
func (h AdminHandler) handleSettingPost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		resTextError(w, "最多接收10M数据")
		return err
	}
	form := r.MultipartForm.Value
	maxMachineStr := h.parseFormText(form, "maxMachine")
	if maxMachineStr == "" {
		resTextError(w, "maxMachine不能为空")
		return nil
	}
	maxMachine, err := strconv.Atoi(maxMachineStr)
	if err != nil {
		resTextError(w, err.Error())
		return nil
	}
	if maxMachine%2 != 0 && maxMachine != 1 {
		resTextError(w, "maxMachine必须为1,或者为2的倍数")
		return nil
	}
	maxNet, err := tool.GetMaxNetworks()
	if err != nil {
		resTextError(w, err.Error())
		return err
	}
	deploy := &model.Deploy{
		MaxMachine:  maxMachine,
		MaxInstance: maxNet / maxMachine,
	}
	err = h.Ctx.MongoClient.DeployRepository().Upsert(deploy)
	if err != nil {
		resTextError(w, err.Error())
		return err
	}
	resJson(w, nil)
	return nil
}

func (h AdminHandler) handleMailPageGet(w http.ResponseWriter, r *http.Request) error {
	t, err := h.ParseTemplate("/mail.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return err
	}
	return t.Execute(w, "")
}

func (h AdminHandler) handleMailPost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 * 1024 * 10)
	form := r.MultipartForm.Value
	nr := h.Ctx.MongoClient.NodeRepository()
	cr := h.Ctx.MongoClient.ContainerRepository()
	rr := h.Ctx.MongoClient.ReceiverRepository()
	sr := h.Ctx.MongoClient.SenderRepository()
	tr := h.Ctx.MongoClient.TemplatesRepository()
	lr := h.Ctx.MongoClient.LogRepository()
	vr := h.Ctx.MongoClient.VariableRepository()
	proxy := h.parseFormText(form, "proxy")
	minNumber := h.parseFormText(form, "minNumber")
	maxNumber := h.parseFormText(form, "maxNumber")
	if proxy == "" || minNumber == "" || maxNumber == "" {
		resError(w, "字段不允许为空", 400)
		return errors.New("字段不允许为空")
	}
	min, err := strconv.Atoi(minNumber)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	max, err := strconv.Atoi(maxNumber)
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	node, err := nr.Find()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if len(node.Nodes) == 0 {
		resError(w, "没有正在运行的机器", 400)
		return nil
	}
	list, err := cr.Find()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if len(list) == 0 {
		resError(w, "没有正在运行的实例", 400)
		return err
	}
	if GetRunningStatus() {
		resError(w, "有正在运行的任务,请稍后重试", 400)
		return nil
	}
	/*
		host := r.URL.Query().Get("host")
			id := r.URL.Query().Get("id")
			dSdk := sdk.NewDeploySdk(host)
			err := dSdk.DeleteContainer(id)
			if err != nil {
				resError(w, err.Error(), 400)
				return err
			}
			resJson(w, nil)
	*/
	runningIns := make(map[string][]string)
	for _, v := range list {
		for _, j := range v.List {
			if strings.Contains(strings.ToUpper(j.Status), "UP") {
				ids := runningIns[v.Ip]
				ids = append(ids, v.Id)
				runningIns[v.Ip] = ids
			}
		}
	}
	total := 0
	for _, v := range runningIns {
		total += len(v)
	}
	if total == 0 {
		resError(w, "没有正在运行的实例", 400)
		return err
	}
	receivers, err := rr.GetAll()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if len(receivers) == 0 {
		resError(w, "收件人为空", 400)
		return nil
	}
	senders, err := sr.GetAll()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if len(senders) == 0 {
		resError(w, "发件人为空", 400)
		return nil
	}
	templates, err := tr.Get()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	if len(templates) == 0 {
		resError(w, "模板为空", 400)
		return nil
	}
	variable, err := vr.Find()
	if err != nil {
		resError(w, err.Error(), 400)
		return nil
	}
	var varMap map[string]string
	if variable != nil {
		varMap = variable.Data
	}
	resJson(w, nil)
	go func() {
		SetRunningStatus(true)
		err := lr.Flush()
		if err != nil {
			log.Println(err)
			SetRunningStatus(false)
			return
		}
		tasks := h.genMailTask(proxy, min, max, runningIns, receivers, senders, templates, varMap)
		start := time.Now().Unix()
		for _, v := range tasks {
			for _, i := range v {
				err := lr.Upsert(&i)
				if err != nil {
					log.Println(err)
				}
			}
		}
		var wg sync.WaitGroup
		taskCount := 0
		var lock sync.Mutex
		var addCount = func() {
			lock.Lock()
			taskCount++
			lock.Unlock()
		}
		var getCount = func() int {
			return taskCount
		}
		for _, v := range tasks {
			wg.Add(1)
			go func(arr []model.Log) {
				for _, task := range v {
					dSdk := sdk.NewDeploySdk(task.MachineIp)
					isOk := false
					err := dSdk.SendMail(task.ContainerId, task.MailType, task.Receiver.Email, task.Title, task.Body, task.Sender.Email, task.Sender.Password, task.Proxy)
					if err != nil {
						log.Printf("总共%d个任务,当前:%d,sender:%s,receiver:%s,失败:%s", len(tasks), getCount()+1, task.Sender.Email, task.Receiver.Email, err)
					} else {
						isOk = true
						log.Printf("总共%d个任务,当前:%d,sender:%s,receiver:%s,成功", len(tasks), getCount()+1, task.Sender.Email, task.Receiver.Email)
					}
					addCount()
					task.Success = isOk
					task.IsProcess = true
					task.CreatedTime = time.Now().Unix()
					err = lr.Upsert(&task)
					if err != nil {
						log.Println("upsert task err:", err)
					}
					err = sr.Upsert(task.Sender)
					if err != nil {
						log.Println("upsert sender err:", err)
					}
					err = rr.Upsert(task.Receiver)
					if err != nil {
						log.Println("upsert receiver err:", err)
					}
				}
				wg.Done()
			}(v)
		}
		wg.Wait()
		log.Printf("运行完成,总共耗时%d秒", time.Now().Unix()-start)
		SetRunningStatus(false)
	}()
	return nil
}

func (h AdminHandler) handleMailLogListGet(w http.ResponseWriter, r *http.Request) error {
	lr := h.Ctx.MongoClient.LogRepository()
	re, err := lr.Get()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}

	draw := r.URL.Query().Get("draw")
	length := r.URL.Query().Get("length")
	value := r.URL.Query().Get("search[value]")
	_ = value
	start := r.URL.Query().Get("start")
	startNum, _ := strconv.Atoi(start)
	lengthNum, _ := strconv.Atoi(length)
	var endNum int
	if startNum > len(re) {
		startNum = len(re)
	}
	if startNum+lengthNum > len(re) {
		endNum = len(re)
	}
	resData := re[startNum:endNum]
	sj := sjson.New()
	sj.Set("draw", draw)
	sj.Set("recordsTotal", len(re))
	sj.Set("recordsFiltered", len(resData))
	sj.Set("data", resData)
	dd, _ := sj.Encode()
	w.Write(dd)
	return nil
}

func (h AdminHandler) handleMailLogDetail(w http.ResponseWriter, r *http.Request) error {
	lr := h.Ctx.MongoClient.LogRepository()
	re, err := lr.Get()
	if err != nil {
		resError(w, err.Error(), 400)
		return err
	}
	var success, fail, isProcess, unProcess int
	for _, v := range re {
		if v.IsProcess {
			isProcess++
		} else {
			unProcess++
		}
		if v.IsProcess && v.Success {
			success++
		}
		if v.IsProcess && !v.Success {
			fail++
		}
	}
	sj := sjson.New()
	sj.Set("success", success)
	sj.Set("fail", fail)
	sj.Set("isProcess", isProcess)
	sj.Set("unProcess", unProcess)
	sj.Set("total", len(re))
	resJson(w, sj)
	return nil
}

func (h AdminHandler) genMailTask(proxy string, min, max int, ins map[string][]string, receivers []*model.ReceiverInfo, senders []*model.SendInfo, templates []*model.Templates, varMap map[string]string) map[string][]model.Log {
	re := make(map[string][]model.Log)
	dockerInfo := make([]string, 0)
	for k, v := range ins {
		for _, i := range v {
			dockerInfo = append(dockerInfo, k+":"+i) //将map打平成一个不重复host:id数组,便于分配任务
		}
	}
	senderCount := 0
	dockerCount := 0
	for i := 0; i < len(receivers); i++ {
		if senderCount >= len(senders) {
			senderCount = 0
		}
		if dockerCount >= len(dockerInfo) {
			dockerCount = 0
		}
		host := strings.Split(dockerInfo[dockerCount], ":")[0]
		dockerId := strings.Split(dockerInfo[dockerCount], ":")[1]
		temp := templates[rand.Intn(len(templates))]
		arr := re[dockerId]
		arr = append(arr, model.Log{
			Id:          uuid.NewV4().String(),
			Sender:      senders[senderCount],
			Body:        temp.FillTemplate(min, max, varMap),
			Title:       temp.Name,
			Receiver:    receivers[i],
			Proxy:       proxy,
			MachineIp:   host,
			MailType:    "163",
			ContainerId: dockerId,
			IsProcess:   false,
		})
		re[dockerId] = arr
		senderCount++
		dockerCount++
	}
	return re
}

func (h AdminHandler) parseRich(form map[string][]string, field string) string {
	for _, v := range form[field] {
		if v != "" {
			return v
		}
	}
	return ""
}

func (h AdminHandler) parseFormText(form map[string][]string, field string) string {
	if len(form[field]) != 0 {
		return form[field][0]
	}
	return ""
}

func (h AdminHandler) parseFile(r *http.Request, field string) (io.Reader, error) {
	f, he, err := r.FormFile(field)
	if err != nil {
		log.Println(err)
	}
	if he != nil {
		if he.Size > 100*1024*1024 {
			return nil, errors.New("文件最大为10M")
		}
		contentType := strings.ToUpper(he.Header.Get("Content-Type"))
		log.Println(contentType)
		return f, nil
	}
	return nil, errors.New("no file received")
}
