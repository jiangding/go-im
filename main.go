package main

import (
	"fmt"
	"html/template"
	"im/controller"
	"log"
	"net/http"

)

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login", controller.UserLogin)
	http.HandleFunc("/user/register", controller.UserRegister)

	http.HandleFunc("/contact/addfriend", controller.AddFriend)
	http.HandleFunc("/contact/loadfriend", controller.LoadFriend)

	http.HandleFunc("/contact/loadcommunity", controller.LoadCommunity)
	http.HandleFunc("/contact/joincommunity", controller.JoinCommunity)

	http.HandleFunc("/contact/createcommunity", controller.CreateCommunity)

	http.HandleFunc("/chat", controller.Chat)

	http.HandleFunc("/attach/upload", controller.Upload)

	// 指定目录的静态文件
	http.Handle("/asset/",http.FileServer(http.Dir(".")))
	http.Handle("/mnt/",http.FileServer(http.Dir(".")))
	// 注册view
	RegisterView()

	// 启动web服务器
	http.ListenAndServe(":8888", nil)

}


func RegisterView(){
	//一次解析出全部模板
	tpl,err := template.ParseGlob("view/**/*")
	if err != nil {
		log.Fatal(err)
	}
	//通过for循环做好映射
	for _,v := range tpl.Templates(){
		//
		tplName := v.Name()
		fmt.Println("HandleFunc+:"+v.Name())
		// 解析
		http.HandleFunc(tplName, func(w http.ResponseWriter,
			request *http.Request) {
			err := tpl.ExecuteTemplate(w,tplName,nil)
			if err!=nil{
				log.Fatal(err.Error())
			}
		})
	}
}

func RegisterTemplate(){
	//全局扫描模板
	GlobTemplete := template.New("root")
	GlobTemplete ,err:=GlobTemplete.ParseGlob("view/**/*")
	if err!=nil {
		//打印错误信息
		//退出系统
		log.Fatal(err)
	}
	//分别对每一个模板进行注册
	for _,tpl := range  GlobTemplete.Templates(){
		patern := tpl.Name()
		http.HandleFunc(patern,
			func(w http.ResponseWriter,
				r *http.Request) {
				GlobTemplete.ExecuteTemplate(w,patern,nil)
			})
		fmt.Println("register=>"+patern)
	}
}
