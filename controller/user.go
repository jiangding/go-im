package controller

import (
	"fmt"
	"im/model"
	"im/service"
	"im/util"
	"math/rand"
	"net/http"
)
var userService service.UserService

// UserLogin 用户登录对应的处理control
func UserLogin(writer http.ResponseWriter, request *http.Request){
	//数据库操作
	//逻辑处理
	//restapi json/xml返回
	//1.获取前端传递的参数
	//mobile,passwd
	//解析参数
	//如何获得参数
	//解析参数
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")

	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(writer,err.Error())
	}else{
		util.RespOk(writer,user)
	}

	// 如果正确
	//if loginOK {
	//	// str = `{"code":0, "data": {"id":1, "token":"test"}}`
	//	data := make(map[string]interface{})
	//	data["id"] = 2
	//	data["token"] = "the token"
	//	util.RespOk(writer,data)
	//}else{
	//	util.RespFail(writer,"密码不正确")
	//}
	//go get github.com/go-xorm/xorm
	//go get github.com/go-sql-driver/mysql
	//返回json ok
}

func UserRegister(writer http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()
	// 获取值
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")

	nickname := fmt.Sprintf("user%06d",rand.Int31())
	avatar :=""
	sex := model.SEX_UNKNOW

	// 调用用户模块, 执行注册
	user,err := userService.Register(mobile,passwd,nickname,avatar,sex)
	if err!=nil{
		util.RespFail(writer,err.Error())
	}else{
		util.RespOk(writer,user)
	}

}