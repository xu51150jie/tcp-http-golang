package tcp

import (
	"encoding/hex"
	"github.com/spf13/viper"
	"io"
	"net"
	"sync"
	"tcp-http-golang/global"
	"tcp-http-golang/http/request"
)

type Server struct {
	IpPort string

	//在线客户端的列表
	OnlineMap map[string]*Client
	MapLock   sync.RWMutex //锁

}

// NewServer 创建一个server的接口
func NewServer(ipPort string) *Server {
	server := &Server{
		IpPort:    ipPort,
		OnlineMap: make(map[string]*Client),
	}

	return server
}

// Start 启动服务器的接口
func (this *Server) Start() {
	// 1：建立端口服务
	listener, err := net.Listen("tcp", this.IpPort)
	if err != nil {
		global.Logger.Errorf("启动失败! err:%v", err)
		return
	}
	//close listen socket
	defer listener.Close()

	// 2.监听客户端连接
	for {
		conn, err := listener.Accept() //accept
		if err != nil {
			global.Logger.Errorf("接受连接失败, err:%v", err)
			continue
		}

		//3.处理接受到的数据
		go this.Handler(conn) // 启动一个goroutine处理连接 多协程处理
	}
}

// 3.处理接受到的数据
func (this *Server) Handler(conn net.Conn) {
	defer conn.Close()
	//...当前链接的业务
	//fmt.Println("链接建立成功")

	client := NewClient(conn)
	this.OnlineClient(client) //客户端上线

	// 处理接受数据
	buf := make([]byte, 4096) //2048表示每次接受的一包数据的字节数量，设置太小的话，会出现长内容被分割几段
	for {
		read, err := conn.Read(buf) // 读取数据字节
		if read == 0 {
			global.Logger.Errorf("客户端断线了 类型1，%s - err:%s", conn.RemoteAddr(), err)
			this.OfflineClient(client)
			return
		}
		if err == io.EOF {
			global.Logger.Errorf("客户端断线了 类型2，%s - err:%s", conn.RemoteAddr(), err)
			this.OfflineClient(client)
			return
		}
		if err != nil {
			global.Logger.Errorf("客户端读取数据失败，%s - err:%s", conn.RemoteAddr(), err)
			this.OfflineClient(client)
			return
		}

		var readData string
		if viper.GetString("SEND_DATA_TYPE") == "HEX" { // 根据配置文件发送ASCII HEX
			readData = hex.EncodeToString(buf[:read]) // 十六进制转十六进制字符串
		} else {
			readData = string(buf[:read]) // ASCII字符串
		}
		global.Logger.Infof("接受数据：%v -> %v", conn.RemoteAddr().String(), readData)

		// 以http方式发送
		go request.SendDataHttp(readData, conn.RemoteAddr().String())
	}

	//当前handler阻塞
	//select {}
}

// 客户端上线
func (this *Server) OnlineClient(client *Client) {
	global.Logger.Infof("新客户端上线%v %v", client.IpPort, client.Time)
	//客户端上线,将用户加入到onlineMap中
	this.MapLock.Lock()
	this.OnlineMap[client.IpPort] = client
	this.MapLock.Unlock()
	//fmt.Println("this.OnlineMap", this.OnlineMap)
}

// 客户端下线
func (this *Server) OfflineClient(client *Client) {
	//客户端下线,将用户从onlineMap中删除
	this.MapLock.Lock()
	delete(this.OnlineMap, client.IpPort)
	this.MapLock.Unlock()

	client.Conn.Close() //关闭连接
}
