package tcp

import (
	"net"
	"time"
)

type Client struct {
	IpPort  string
	MsgChan chan []byte
	Conn    net.Conn
	Time    string
}

// 监听当前Client channel的 方法,一旦有消息，就直接发送给对端客户端
func (this *Client) ListenMessage() {
	for {
		msg := <-this.MsgChan
		this.Conn.Write(msg)
	}
}

// 创建一个客户端的API
func NewClient(conn net.Conn) *Client {
	clientAddr := conn.RemoteAddr().String()

	client := &Client{
		IpPort:  clientAddr,
		MsgChan: make(chan []byte),
		Conn:    conn,
		Time:    time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
	}

	//启动监听当前client channel消息的goroutine
	go client.ListenMessage()

	return client
}
