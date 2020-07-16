package controllers

import (
	ser "chatservice/service"
	"chatservice/service/logic"
	"chatservice/service/web"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// 聊天接口说明
type MsgController struct {
	beego.Controller
}

// @Title 创建websocket链接
// @Summary 创建websocket链接
// @Description 根据用户信息创建websocket链接
// @Param	userid	path 	int	true		"当前连接的用户id"
// @router /Conn [get]
func (m *MsgController) Conn() {
	userid, _ := m.GetInt("userid", 0)
	ws, err := websocket.Upgrade(m.Ctx.ResponseWriter, m.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error(m.Ctx.ResponseWriter, "不是websocket连接", 400)
		return
	} else if err != nil {
		beego.Error("无法设置websocket连接:", err)
		return
	}
	client := logic.NewUser(userid, ws, make(chan *logic.SendMessage, 256))
	client.Join(ws)
}

// @Title 创建会话
// @Summary 创建会话
// @Description 根据用户创建会话
// @Param	targetid		formData 	int		true	"目标用户ID(你想和谁聊天)"
// @Param	userid			formData	int		true	"当前用户ID(你是谁)"
// @Success 200 {object} models.web.WebResponse
// @router /CreateConversation [post]
func (c *MsgController) CreateConversation() {
	targetid, _ := c.GetInt("targetid", 0)
	userid, _ := c.GetInt("userid", 0)
	cvsid, msg := (&ser.Conversation{}).CreateConversation(targetid, userid)
	//返回会话ID
	web.ResponseData(cvsid, msg, &c.Controller)
}

// @Title 删除会话
// @Summary 删除会话
// @Description 删除会话
// @Param	cvsid		formData 	int		true	"删除的会话id"
// @Param	userid		formData	int	 	true	"操作的用户ID"
// @Success 200 {object} models.web.WebResponse
// @router /DelConversation [post]
func (c *MsgController) DelConversation() {
	cvsid, _ := c.GetInt("cvsid", 0)
	userid, _ := c.GetInt("userid", 0)
	code, msg := (&ser.Conversation{}).DelConversation(cvsid, userid)
	web.ResponseData(code, msg, &c.Controller)
}

// @Title 获取会话列表
// @Summary 获取会话列表
// @Description 获取会话列表包含未读消息数和最后一条消息内容
// @Param	userid			formData		int	true	"用户ID"
// @Success 200 {object} logic.Conversation
// @router /GetConversation [post]
func (c *MsgController) GetConversation() {
	userid, _ := c.GetInt("userid", 0)
	list, _ := (&ser.Conversation{}).GetConversation(userid)
	web.ResponseData(list, "获取会话列表成功", &c.Controller)

}

// @Title 获取会话消息数据
// @Summary 获取会话消息数据
// @Description 获取会话消息数据
// @Param	userid			formData		int	true	"用户ID"
// @Param	maxid			formData		int	false	"消息ID"
// @Param	cvsid			formData		int	true	"会话ID"
// @Success 200 {object} logic.Conversation
// @router /GetMsgList [post]
func (c *MsgController) GetMsgList() {
	userid, _ := c.GetInt("userid", 0)
	cvsid, _ := c.GetInt("cvsid", 0)
	maxid, _ := c.GetInt("maxid", 0)
	_, msg, data := ser.GetMsgList(maxid, cvsid, userid)
	beego.Debug(data)
	web.ResponseData(data, msg, &c.Controller)

}

// @Title 发消息
// @Summary 发消息
// @Description 发消息
// @Param	cvsid		formData 	int		true	会话id
// @Param	content		formData	string	false	消息内容
// @Param	userid		formData	int		true	发消息的用户ID
// @router /SendMsg [post]
func (c *MsgController) SendMsg() {
	//发送消息的用户ID
	userid, _ := c.GetInt("userid", 0)
	//会话id
	cvsid, _ := c.GetInt("cvsid", 0)
	//消息内容
	content := c.GetString("content", "")
	if len(content) == 0 {
		web.ResponseState(0, "消息不能为空", &c.Controller)
		return
	}
	code, msg := ser.SendMsg(userid, cvsid, content)
	web.ResponseData(code, msg, &c.Controller)
}
