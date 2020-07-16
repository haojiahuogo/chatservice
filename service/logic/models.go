package logic

//发送消息数据
type SendMessage struct {
	Id       int    `json:"id" description:"消息ID"`
	Cvsid    int    `json:"cvsid" description:"会话id"`
	Source   int    `json:"source" description:"消息发送者ID"`
	Content  string `json:"content" description:"消息内容"`
	Sendtime int    `json:"sendtime" description:"发送时间"`
	Userimg  string `json:"userimg" description:"发送者头像"`
	Username string `json:"username" description:"发送者名称"`
}

//会话消息列表数据
type Conversation struct {
	Cvsid    int    `json:"cvsid" description:"会话ID"`
	Msgid    int    `json:"msgid" description:"最后一条消息ID"`
	Userlist string `json:"userlist" description:"当前会话所有用户ID"`
	Source   int    `json:"source" description:"最后一条消息发送者ID"`
	Content  string `json:"content" description:"最后一条消息内容"`
	Type     int    `json:"type" description:"最后一条消息类型"`
	Sendtime int    `json:"sendtime" description:"最后一条消息发送时间"`
	Notsee   int    `json:"notsee" description:"当前会话未读消息条数"`
	Lastid   int    `json:"lastid" description:"最后阅读消息ID"`
	Userimg  string `json:"userimg" description:"好友头像"`
	Username string `json:"username" description:"好友名称"`
}

type CvsNotsee struct {
	Cvsid  int `json:"cvsid" description:"会话ID"`
	Notsee int `json:"notsee" description:"未读消息条数"`
}
