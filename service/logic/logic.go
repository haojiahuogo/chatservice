package logic

import (
	"github.com/astaxie/beego"
)

//初始化代码
func init() {
	//使用协程的方式开启
	go chatroom()
}

var (
	register    = make(chan *Client)      //注册加入用户通道
	unregister  = make(chan *Client)      //注销用户通道
	publish     = make(chan *SendMessage) //发送消息通道
	Onlineusers = make(map[int]*Client)   //所有在线用户
)

func chatroom() {
	for {
		select {
		case sub := <-register:
			//加入在线用户
			Onlineusers[sub.Userid] = sub
			beego.Debug(sub.Userid, "加入系统")
			beego.Debug("当前在线人数:", len(Onlineusers))
		case unsub := <-unregister:
			//注销当前用户 断开websocket连接
			ws := unsub.Conn
			if ws != nil {
				ws.Close()
			}
			//将当前用户从在线用户列表中移除
			delete(Onlineusers, unsub.Userid)
			beego.Debug(unsub.Userid, "退出")
			beego.Debug("当前在线人数:", len(Onlineusers))
		case msg := <-publish:
			beego.Debug("消息通道中收到消息:", msg)
			//TODO 后续代码处理
			//如果需要给所有用户发送消息，则可以循环Onlineusers在线用户，给每一个用户发送
			//也可以根据消息中msg.Target 给指定用户发送消息

		}
	}
}
