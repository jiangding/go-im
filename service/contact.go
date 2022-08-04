package service

import (
	"errors"
	"im/model"
	"time"
)

type ContactService struct {

}

func (cs *ContactService) GetFriend(userId int64) []model.User{
	// 创建contact
	contacts := make([]model.Contact, 0)
	// 找出 own是当前userId，类型是用户
	DbEngin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_USER).Find(&contacts)

	objIds := make([]int64, 0)
	for _,v := range contacts {
		objIds = append(objIds, v.Dstobj) // 收集目标用户
	}
	coms := make([]model.User,0)
	if len(objIds) == 0 {
		return coms
	}
	// 找到对应的用户们返回
	DbEngin.In("id", objIds).Find(&coms)

	return coms
}

// AddFriend 添加好友
func (cs *ContactService) AddFriend(userId, dstid int64) error {
	if userId == dstid {
		return errors.New("不能添加自己为好友")
	}
	//判断是否已经加了好友
	tmp := model.Contact{}
	//查询是否已经是好友
	// 条件的链式操作
	DbEngin.Where("ownerid = ?",userId).And("dstid = ?",dstid).And("cate = ?",model.CONCAT_CATE_USER).Get(&tmp)
	//获得1条记录
	//count()
	//如果存在记录说明已经是好友了不加
	if tmp.Id>0{
		return errors.New("该用户已经被添加过啦")
	}

	//user := model.User{}
	//// 当前用户不存在
	//DbEngin.Where("id", dstid).Find(&user)
	//if user.Id == 0 {
	//	return errors.New("目标用户不存在")
	//}

	//事务,
	session := DbEngin.NewSession()
	session.Begin()
	//插自己的
	_,e2 := session.InsertOne(model.Contact{
		Ownerid:userId,
		Dstobj:dstid,
		Cate:model.CONCAT_CATE_USER,
		Createat:time.Now(),
	})
	//插对方的
	_,e3 := session.InsertOne(model.Contact{
		Ownerid:dstid,
		Dstobj:userId,
		Cate:model.CONCAT_CATE_USER,
		Createat:time.Now(),
	})
	//没有错误
	if  e2==nil && e3==nil{
		//提交
		session.Commit()
		return nil
	}else{
		//回滚
		session.Rollback()
		if e2!=nil{
			return e2
		}else{
			return e3
		}
	}
}

// SearchComunityIds 某个用户加入的所有群
func (cs *ContactService) SearchComunityIds(userId int64) (comIds []int64){
	//todo 获取用户全部群ID
	conconts := make([]model.Contact,0)
	comIds =make([]int64,0)

	DbEngin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _,v := range conconts{
		comIds = append(comIds,v.Dstobj);
	}
	return comIds
}


// SearchComunity 当前用户加入过的群列表
func (cs *ContactService) SearchComunity(userId int64) []model.Community {
	conconts := make([]model.Contact,0)
	comIds :=make([]int64,0)

	DbEngin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _,v := range conconts{
		comIds = append(comIds,v.Dstobj);
	}
	coms := make([]model.Community,0)
	if len(comIds)== 0{
		return coms
	}
	DbEngin.In("id",comIds).Find(&coms)
	return coms
}


// JoinCommunity 加群
func (cs *ContactService) JoinCommunity(userId,comId int64) error {
	cot := model.Contact{
		Ownerid:userId,
		Dstobj:comId,
		Cate:model.CONCAT_CATE_COMUNITY,
	}
	DbEngin.Get(&cot)
	if cot.Id==0 {
		cot.Createat = time.Now()
		_,err := DbEngin.InsertOne(cot)
		return err
	}else{
		return nil
	}
}

// CreateCommunity 建群
func (cs *ContactService) CreateCommunity(comm model.Community) (ret model.Community,err error){
	if len(comm.Name)==0{
		err = errors.New("缺少群名称")
		return ret,err
	}
	if comm.Ownerid==0{
		err = errors.New("请先登录")
		return ret,err
	}
	com := model.Community{
		Ownerid:comm.Ownerid,
	}
	num,err := DbEngin.Count(&com)
	if num>5 {
		err = errors.New("一个用户最多只能创见5个群")
		return com,err
	}else{
		comm.Createat=time.Now()
		session := DbEngin.NewSession()
		session.Begin()
		_,err = session.InsertOne(&comm)
		if err!=nil{
			session.Rollback();
			return com,err
		}
		_,err =session.InsertOne(
			model.Contact{
				Ownerid:comm.Ownerid,
				Dstobj:comm.Id,
				Cate:model.CONCAT_CATE_COMUNITY,
				Createat:time.Now(),
			})
		if err!=nil{
			session.Rollback();
		}else{
			session.Commit()
		}
		return com,err
	}
}
