package service

import (
	"chatservice/models"
	"chatservice/service/logic"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Conversation struct {
}

//创建会话ID
func (c *Conversation) CreateConversation(targetid, userid int) (int, string) {
	o := orm.NewOrm()
	cvsinfo := models.Conversation{}
	sql := "select * from conversation where userlist like ? and userlist like ?"
	tid := strconv.Itoa(targetid)
	uid := strconv.Itoa(userid)
	err := o.Raw(sql, "%,"+tid+",%", "%,"+uid+",%").QueryRow(&cvsinfo)
	cvsid := 0
	if err != nil {
		//如果没有和对方聊过天就新创建主会话信息
		var conversation models.Conversation
		conversation.Userlist = "," + uid + "," + tid + ","
		cvsnum, _ := o.Insert(&conversation)
		cvsid, _ = strconv.Atoi(strconv.FormatInt(cvsnum, 10))
	} else {
		//返回已经存在的会话ID
		cvsid = cvsinfo.Id
	}

	var cvslist []models.ConversationUser
	o.Raw("select * from conversation_user where cvsid=?", cvsid).QueryRows(&cvslist)
	//创建map存储
	cvslistmap := make(map[int]int, len(cvslist))
	for _, v := range cvslist {
		cvslistmap[v.Userid] = v.Userid
	}
	//用户如果删除会话列表 再次创建会话列表数据
	//获取删除之前最后消息ID
	var msglist []logic.SendMessage
	o.Raw("select max(id) as id,source from message where cvsid=? GROUP BY source", cvsid).QueryRows(&msglist)
	msglistmap := make(map[int]int, len(msglist))
	if len(msglist) > 0 {
		for _, v := range msglist {
			msglistmap[v.Source] = v.Id
		}
	}
	_, ok1 := cvslistmap[userid]
	_, ok2 := cvslistmap[targetid]
	//删除之后以前的数据最大ID设置成最大阅读ID
	wlasid, _ := msglistmap[userid]
	tlasid, _ := msglistmap[targetid]
	//创建对应的两条用户关联数据
	chatcvsuser := []models.ConversationUser{}
	if !ok1 && !ok2 {
		chatcvsuser = []models.ConversationUser{
			{Cvsid: cvsid, Userid: userid, Lastid: wlasid},
			{Cvsid: cvsid, Userid: targetid, Lastid: tlasid},
		}
	} else if ok1 && !ok2 {
		chatcvsuser = []models.ConversationUser{
			{Cvsid: cvsid, Userid: targetid, Lastid: tlasid},
		}
	} else if !ok1 && ok2 {
		chatcvsuser = []models.ConversationUser{
			{Cvsid: cvsid, Userid: userid, Lastid: wlasid},
		}
	}
	//插入数据
	o.InsertMulti(100, chatcvsuser)
	return cvsid, "创建会话成功"
}

//删除会话
func (c *Conversation) DelConversation(cvsid, userid int) (int, string) {
	o := orm.NewOrm()
	res, err := o.Raw("delete FROM conversation_user where cvsid=? and userid=?", cvsid, userid).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		beego.Debug(num)
		return 1, "删除会话成功！"
	}
	return 0, "删除会话失败！"
}

//获取会话列表数据
func (this *Conversation) GetConversation(userid int) ([]*logic.Conversation, error) {
	list := []*logic.Conversation{}
	o := orm.NewOrm()
	//当前用户的会话消息列表
	var cvslist []logic.Conversation
	sql := "select cu.*,c.userlist,m.* from conversation_user cu,conversation c,(select cvsall.id as msgid,cvsall.cvsid as mcvsid,cvsall.source,cvsall.content,cvsall.sendtime  from (select m.* from message m,(select cvsid from conversation_user where userid =?) s where s.cvsid=m.cvsid  ORDER BY m.id desc ) cvsall GROUP BY cvsall.cvsid) m where cu.cvsid=c.id and m.mcvsid=cu.cvsid and cu.userid=?"
	_, e := o.Raw(sql, userid, userid).QueryRows(&cvslist)
	if e != nil {
		return nil, e
	}
	beego.Debug(len(cvslist))
	//查询所有会话未读消息数
	var CvsNotseelist []logic.CvsNotsee
	cvsmsgsql := "select cm.cvsid,count(1) as notsee from (select msg.*,cv.userlist from message msg ,conversation cv where cv.id=msg.cvsid and cv.userlist like '%," + strconv.Itoa(userid) + ",%' ) cm,conversation_user cu where cu.cvsid=cm.cvsid and cu.userid=?  and cm.id>cu.lastid GROUP BY cvsid  "
	beego.Debug(cvsmsgsql)
	o.Raw(cvsmsgsql, userid).QueryRows(&CvsNotseelist)
	NotseeMap := make(map[int]int, len(CvsNotseelist))
	for _, v := range CvsNotseelist {
		NotseeMap[v.Cvsid] = v.Notsee
	}
	//获取会话列表中用户信息 获取对方用户信息
	userlist := []int{}
	cvsMap := make(map[int]int, len(cvslist))
	for _, c := range cvslist {
		if c.Userlist != "" {
			//获取对方的用户id
			otherid := strings.Replace(strings.Replace(c.Userlist, ","+strconv.Itoa(userid)+",", "", -1), ",", "", -1)
			otherids, _ := strconv.Atoi(otherid)
			userlist = append(userlist, otherids)
			cvsMap[c.Cvsid] = otherids
		}
	}
	//会话列表中对方的用户信息
	db := orm.NewOrm()
	var userinfo []models.User
	if len(userlist) > 0 {
		_, e = db.QueryTable(new(models.User)).Filter("Id__in", userlist).All(&userinfo)
		if e != nil {
			return nil, e
		}
	}
	userMap := make(map[int]models.User, len(userinfo))
	for _, u := range userinfo {
		userMap[u.Id] = u
	}
	for _, c := range cvslist {
		ortherid, _ := cvsMap[c.Cvsid]
		user, ok := userMap[ortherid]
		notsee, sok := NotseeMap[c.Cvsid]
		if !sok {
			notsee = 0
		}
		var users *models.User
		if ok {
			users = &user
		}
		list = append(list, NewConversation(&c, users, notsee))
	}

	return list, nil
}

func NewConversation(conversation *logic.Conversation, user *models.User, notsee int) *logic.Conversation {
	conv := &logic.Conversation{
		Cvsid:    conversation.Cvsid,
		Lastid:   conversation.Lastid,
		Msgid:    conversation.Msgid,
		Userlist: conversation.Userlist,
		Source:   conversation.Source,
		Content:  conversation.Content,
		Sendtime: conversation.Sendtime,
		Type:     conversation.Type,
		Notsee:   notsee,
	}
	if user != nil {
		conv.Username = user.Name
		conv.Userimg = user.Img
	}
	return conv
}
