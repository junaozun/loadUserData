package main

import (
	"fmt"
	"html/template"
	"loadOnlineData/etcd"
	"loadOnlineData/load"
	"loadOnlineData/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type InputData struct {
	ServerName string
	ServerId int
	Ip string
	EndServerName string
}

var NeedData  *InputData

// 渲染页面并输出
func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	// 获取页面内容
	t, err := template.New(file).ParseFiles("views/" + file)
	utils.CheckErr(err)
	t.Execute(w, data)
}

func index(w http.ResponseWriter, r *http.Request) {
	// 渲染页面并输出
	renderHTML(w, "index.html", "no data")
}

func insert(w http.ResponseWriter, r *http.Request) {
	// 必须通过 POST 提交数据
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println("Handler:page:ParseForm: ", err)
		}

		NeedData = &InputData{}
		serName := r.Form.Get("serverName")
		if serName == "" {
			w.Write([]byte("服务器名不能为空，请返回重新输入！"))
			return
		}
		NeedData.ServerName = serName

		optime := r.Form.Get("openTime")
		stamp,err:= time.ParseInLocation("2006-01-02 15:04:05",optime,time.Local)
		if err != nil{
			w.Write([]byte("时间格式输入错误，请返回重新输入!"))
			return
		}
		serverId,ip,endServerName := load.LoadData(NeedData.ServerName,stamp.Unix())
		NeedData.ServerId = serverId
		NeedData.Ip = ip
		NeedData.EndServerName = endServerName

		renderHTML(w,"sync.html",NeedData)
	} else {
		// 如果不是通过 POST 提交的数据，则将页面重定向到主页
		renderHTML(w, "redirect.html", "/")
	}
}

func syncEtcd(w http.ResponseWriter, r *http.Request) {
	// 必须通过 POST 提交数据
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Println("Handler:page:ParseForm: ", err)
		}
		//serverId := r.Form.Get("serverId")
		//sid,err := strconv.Atoi(serverId)
		//utils.CheckErr(err)
		serverId := NeedData.ServerId
		strSid := strconv.Itoa(serverId)
		etcdKey := "/game/sanguo/dev/service/config/"+strSid
		etcdValue := load.QueryServerData(int32(serverId))
		etcd.PutEtcd(etcdKey,string(etcdValue))
		log.Println("同步成功！")
		w.Write([]byte("同步成功"))
	}else {
		// 如果不是通过 POST 提交的数据，则将页面重定向到主页
		renderHTML(w, "redirect.html", "/")
	}
}

func main() {
	fmt.Println("127.0.0.1:9090")
	http.HandleFunc("/", index)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/syncEtcd", syncEtcd)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println(NeedData.ServerName)
}
