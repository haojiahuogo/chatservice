package service

import (
	"chatservice/models"
	"chatservice/service/logic"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//获取所有的消息
func GetMsgList(maxid, cvsid, userid int) (int, string, []logic.SendMessage) {
	var cmlist []logic.SendMessage
	db := orm.NewOrm()
	var query orm.RawSeter
	if maxid == 0 {
		query = db.Raw("select * from message m,user w where m.source=w.id and m.cvsid=?  ORDER BY m.id desc LIMIT ?", cvsid, 20)
	} else {
		query = db.Raw("select * from message m,user w where m.source=w.id and m.cvsid=? and m.id<? ORDER BY m.id desc LIMIT ?", cvsid, maxid, 20)
	}

	count, e := query.QueryRows(&cmlist)
	var laststruct lastStruct
	lastid := 0
	errlast := db.Raw("select max(id) as lastid from message m where m.cvsid=?", cvsid).QueryRow(&laststruct)
	if errlast == nil {
		lastid = laststruct.Lastid
	}
	if e == nil && count > 0 {
		SetUserMsgLastId(cvsid, userid, lastid)
	}
	return 1, "获取消息列表成功", cmlist
}

type lastStruct struct {
	Lastid int
}

//获取用户的最后阅读消息ID
func GetUserMessageLastId(userid, cvsid int) int {
	db := orm.NewOrm()
	var rows []orm.Params
	count, e := db.Raw("select lastid from conversation_user where cvsid = ? and userid=?", cvsid, userid).Values(&rows)
	if e == nil && count > 0 {
		lastid, _ := strconv.Atoi(rows[0]["lastid"].(string))
		return lastid
	}
	return 0
}

//修改用户的最后阅读消息ID
func SetUserMsgLastId(cvsid, userid, lastid int) {
	db := orm.NewOrm()
	num, err := db.QueryTable(new(models.ConversationUser)).Filter("cvsid", cvsid).Filter("userid", userid).Filter("lastid__lt", lastid).Update(orm.Params{
		"lastid": lastid,
	})
	beego.Debug(num, err)
}

//发送消息
func SendMsg(userid, cvsid int, content string) (int, string) {
	//消息的发送时间
	sendtime := int(time.Now().Unix())
	db := orm.NewOrm()
	//保存消息到数据库中
	var instermsg models.Message
	instermsg.Cvsid = cvsid
	instermsg.Content = content
	instermsg.Sendtime = sendtime
	instermsg.Source = userid
	msgid, err := db.Insert(&instermsg)
	if err != nil || msgid == 0 {
		return 0, err.Error()
	}

	//创建用于存储的消息数据
	msg := &logic.SendMessage{
		Id:       int(msgid),
		Source:   userid,
		Cvsid:    cvsid,
		Sendtime: sendtime,
		Content:  content,
	}
	//给用户发送消息
	state, sedmsg := SendToUser(msg)

	//发送消息更新最后阅读ID
	SetUserMsgLastId(cvsid, userid, int(msgid))
	return state, sedmsg

}

//发送消息给指定用户
func SendToUser(msg *logic.SendMessage) (int, string) {
	db := orm.NewOrm()
	userid := msg.Source
	cvsid := msg.Cvsid
	conversation := &models.Conversation{}
	targetid := 0
	err := db.Raw("select * from conversation where id=?", cvsid).QueryRow(conversation)
	if err == nil {
		if len(conversation.Userlist) > 0 {
			conversationlist := strings.TrimLeft(strings.TrimRight(conversation.Userlist, ","), ",")
			userarray := strings.Split(conversationlist, ",")
			for _, v := range userarray {
				vid, _ := strconv.Atoi(v)
				if vid != userid {
					targetid = vid
				}
			}
		}
	}

	//如果对方不存在会话 重新创建会话信息
	var cvslist []models.ConversationUser
	db.Raw("select * from conversation_user where cvsid=?", cvsid).QueryRows(&cvslist)
	cvslistmap := make(map[int]int, len(cvslist))
	for _, v := range cvslist {
		cvslistmap[v.Userid] = v.Userid
	}
	_, ok2 := cvslistmap[targetid]
	if !ok2 {
		//用户如果删除会话列表 再次创建会话列表数据
		//获取删除之前最后消息ID
		var msglist []logic.SendMessage
		db.Raw("select max(id) as id,source from message where cvsid=? GROUP BY source", cvsid).QueryRows(&msglist)
		msglistmap := make(map[int]int, len(msglist))
		if len(msglist) > 0 {
			for _, v := range msglist {
				msglistmap[v.Source] = v.Id
			}
		}
		tlasid, _ := msglistmap[targetid]
		chatcvsuser := []models.ConversationUser{
			{Cvsid: cvsid, Userid: targetid, Lastid: tlasid},
		}
		db.InsertMulti(10, chatcvsuser)
	}
	//获取发送者头像和姓名信息
	userinfo := models.User{}
	db.QueryTable(new(models.User)).Filter("id", userid).One(&userinfo)
	sendmsg := &logic.SendMessage{
		Id:       msg.Id,
		Cvsid:    msg.Cvsid,
		Source:   msg.Source,
		Content:  msg.Content,
		Sendtime: msg.Sendtime,
		Userimg:  userinfo.Img,
		Username: userinfo.Name,
	}
	//目标用户ID
	beego.Debug(msg.Source, "给目标用户", targetid, "发来消息：", msg.Content)
	if targetid == 0 {
		return 0, "目标用户不存在"
	}
	u, ok := logic.Onlineusers[targetid]
	if ok {
		//用户在线直接发送消息
		u.Send(sendmsg)
	} else {
		//用户不在线推送消息
		beego.Debug("目标用户不在线")
	}
	return 1, "已发送"
}
