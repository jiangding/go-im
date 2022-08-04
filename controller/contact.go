package controller

import (
	"im/args"
	"im/model"
	"im/service"
	"im/util"
	"net/http"
)

var contactService service.ContactService

func LoadFriend(w http.ResponseWriter, req *http.Request){
	// 获取arg
	var arg args.ContactArg
	// 直接跟参数绑定
	util.Bind(req, &arg)

	// 查找所有的名下的用户
	users := contactService.GetFriend(arg.Userid)

	util.RespOkList(w,users,len(users))
}

// AddFriend 添加用户
func AddFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	util.Bind(req,&arg)
	//调用service
	err := contactService.AddFriend(arg.Userid,arg.Dstid)
	if err != nil {
		util.RespFail(w,err.Error())
	}else{
		util.RespOk(w,"好友添加成功")
	}
}

// CreateCommunity 创建群
func CreateCommunity(w http.ResponseWriter, req *http.Request){
	var arg model.Community
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	com,err := contactService.CreateCommunity(arg);
	if err!=nil{
		util.RespFail(w,err.Error())
	}else {
		util.RespOk(w,com)
	}
}

// JoinCommunity 加入群
func JoinCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	util.Bind(req,&arg)
	//调用service
	err := contactService.JoinCommunity(arg.Userid,arg.Dstid)
	if err != nil {
		util.RespFail(w,err.Error())
	}else{
		util.RespOk(w,"")
	}
}
// LoadCommunity 加载群列表
func LoadCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	comunitys := contactService.SearchComunity(arg.Userid)
	util.RespOkList(w,comunitys,len(comunitys))
}

