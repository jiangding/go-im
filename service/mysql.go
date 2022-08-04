package service

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"log"

	_ "github.com/go-sql-driver/mysql"
)
var DbEngin *xorm.Engine
func init()  {
	drive :="mysql"
	DsName := "root:lvdouiot168@(120.24.184.167:3306)/chat?charset=utf8"
	err := errors.New("")
	DbEngin,err = xorm.NewEngine(drive,DsName)
	if nil!=err && ""!=err.Error() {
		log.Fatal(err.Error())
	}
	//是否显示SQL语句
	DbEngin.ShowSQL(true)
	//数据库最大打开的连接数
	DbEngin.SetMaxOpenConns(100)

	//自动创建User
	// DbEngin.Sync2(new(model.User), new(model.Contact), new(model.Community))

	fmt.Println("init data base ok")
}