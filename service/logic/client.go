package logic

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

const (
	// 设置websocket链接的超时时间
	writeWait = 10 * time.Second
	// 设置心跳的时间1分钟
	pongWait = 60 * time.Second
	// 心跳检测的间隔时间 且小于心跳的时间
	pingPeriod = (pongWait * 9) / 10
	// 消息大小限制
	maxMessageSize = 512
)

//客户端结构
type Client struct {
	Userid   int
	Conn     *websocket.Conn
	Sendchan chan *SendMessage
}

//创建一个用户
func NewUser(userid int, ws *websocket.Conn, send chan *SendMessage) *Client {
	user := &Client{
		Userid:   userid,
		Conn:     ws,
		Sendchan: send,
	}
	return user
}

//加入聊天
func (c *Client) Join(ws *websocket.Conn) {
	register <- &Client{Userid: c.Userid, Conn: ws, Sendchan: c.Sendchan}
	//开始监听发送消息
	go c.beginSend()
	//开始监听读取消息
	go c.beginRead()
	//执行心跳检测
	go c.ProcLoop()
}

//退出
func (c *Client) Leave() {
	unregister <- c
}

//发送一条消息
func (c *Client) Send(msg *SendMessage) {
	c.Sendchan <- msg
}

//开启一个携程执行监听读取消息
func (c *Client) beginRead() {
	defer func() {
		c.Leave()
	}()
	c.Conn.SetReadLimit(maxMessageSize) //设置获取消息的大小限制
	//设置websocket连接超时时间
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		//从websocket获取消息缓冲
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				beego.Info("error: %v", err)
			}
			break
		}
		beego.Info(c.Userid, "发来消息：", string(message))
		//将消息放入发送通道中
		//publish <- &SendMessage{}
	}
}

//开启一个携程执行监听发送消息
func (c *Client) beginSend() {
	defer func() {
		c.Leave()
	}()
	for {
		message, ok := <-c.Sendchan
		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		if !ok {
			//通道关闭
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		ws := c.Conn
		if ws != nil {
			//发送消息
			if ws.WriteMessage(websocket.TextMessage, []byte{}) != nil {
				unregister <- c
			}
		}

		//发送消息
		c.Conn.WriteJSON(message)
	}
}

// 开启一个携程执行 心跳检测
func (c *Client) ProcLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//心跳停止了 关掉当前连接
		beego.Debug(c.Userid, "心跳停止")
		ticker.Stop()
		c.Leave()
	}()

	for {
		<-ticker.C
		//beego.Info(c.Name, "`•.¸¸.•´´¯`••.¸¸.•´´❤`•.¸¸.•´´¯`••.¸¸.•´´")
		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("heartbeat")); err != nil {
			beego.Debug(err.Error())
			return
		}
	}

}
