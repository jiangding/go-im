package service

import (
	"errors"
	"fmt"
	"im/model"
	"im/util"
	"math/rand"
	"time"

)

type UserService struct {

}

func (s *UserService)Register(
	mobile,
	pwd,
	name,
	avatar,
	sex string) (user model.User, err error) {

	// 初始化
	tmp := model.User{}
	_,err= DbEngin.Where("mobile=? ",mobile).Get(&tmp)
	if err!=nil{
		return tmp,err
	}
	//如果存在则返回提示已经注册
	if tmp.Id>0{
		return tmp,errors.New("该手机号已经注册")
	}

	//否则拼接插入数据
	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = name
	tmp.Sex = sex
	tmp.Salt = fmt.Sprintf("%06d",rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(pwd,tmp.Salt)
	tmp.Createat = time.Now()
	//token 可以是一个随机数
	tmp.Token = fmt.Sprintf("%08d",rand.Int31())
	//passwd =
	//md5 加密
	//返回新用户信息

	//插入 InserOne
	_,err = DbEngin.InsertOne(&tmp)
	//前端恶意插入特殊字符
	//数据库连接操作失败
	return tmp,err
}

func (s *UserService)Login(
	mobile,//手机
	pwd string )(user model.User,err error) {

	tmp := model.User{}
	// 手机号查询用户
	DbEngin.Where("mobile = ?", mobile).Get(&tmp)

	// 没找到用户
	if tmp.Id == 0 {
		return tmp, errors.New("该用户不存在")
	}

	// 对比密码是否正确
	if !util.ValidatePasswd(pwd,tmp.Salt, tmp.Passwd) {
		return tmp, errors.New("密码不正确")
	}

	//刷新token,安全
	str := fmt.Sprintf("%d",time.Now().Unix())
	token := util.MD5Encode(str)
	tmp.Token = token
	//返回数据
	DbEngin.ID(tmp.Id).Cols("token").Update(&tmp)

	// 返回数据
	return tmp,nil
}

// Find 查找某个用户
func (s *UserService)Find( userId int64 )(user model.User) {
	//首先通过手机号查询用户
	tmp :=model.User{}
	DbEngin.ID(userId).Get(&tmp)
	return tmp
}





