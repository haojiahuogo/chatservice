package main

import (
	_ "chatservice/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//初始化加载数据库
func init() {
	sqlconn := beego.AppConfig.String("sqlconn")
	orm.RegisterDataBase("default", "mysql", sqlconn, 30, 100)
	orm.DefaultRowsLimit = -1
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
