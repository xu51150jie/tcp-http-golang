package main

import (
	"github.com/spf13/viper"
	"tcp-http-golang/global"
	"tcp-http-golang/http/response"
	"tcp-http-golang/tcp"
	"tcp-http-golang/utils"
	"time"
)

func init() {
	global.StartTime = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05") //记录启动时间
	utils.InitConfig()                                                               // 初始化配置文件
	global.Logger = utils.InitLogger()                                               // 初始化日志功能
	// 初始化http接口
}

func main() {
	global.Logger.Infof("运行...")
	global.Logger.Errorf("3333333%s", "22222")

	//新建tcp服务器 IP端口读取配置文件
	var tcp_ip_port = viper.GetString("TCP_IP_PORT")
	tcpServer := tcp.NewServer(tcp_ip_port)
	go showClient(tcpServer) // 定时日志输出 当前在线客户端

	go response.StartHandle(tcpServer) // http服务

	tcpServer.Start() // 启动TCP服务
	global.Logger.Infof("TCP服务器已经开启：%s", tcp_ip_port)
}

// 显示当前在线客户端
func showClient(server *tcp.Server) {
	for {
		global.Logger.Infof(" +----------------------------------------------------------------")
		server.MapLock.Lock()
		global.Logger.Infof(" |  总数：%v  服务启动日期：%v", len(server.OnlineMap), global.StartTime)
		global.Logger.Infof(" +----------------------------------------------------------------")
		for _, cli := range server.OnlineMap {
			global.Logger.Infof(" |  %v  %v", cli.IpPort, cli.Time)
		}
		server.MapLock.Unlock()
		global.Logger.Infof(" +----------------------------------------------------------------")
		global.Logger.Infof("")
		time.Sleep(50 * time.Second)
	}
}

//// 过滤屏蔽非需数据 后期需要改成可配置
//func isTrueData(readData string) bool {
//
//	if len(readData) <= 10 { //屏蔽短内容数据 正常数据长度应大于10以上
//		logrus.Warnf("非需数据%v", readData)
//		return false
//	} else if readData[0:2] != "##" { //过滤不以##开头的数据
//		logrus.Warnf("非需数据%v", readData)
//		return false
//	} else if readData[0:8] == "##PRO=AC" {
//		logrus.Infof("PRO数据%v", readData)
//		return true
//	} else if !strings.Contains(readData, "ST=31") && !strings.Contains(readData, "ST=32") {
//		//判断数据需要是 ST=31或ST=32或ST=40或ST=99
//		logrus.Warnf("非需数据%v", readData)
//		return false
//	} else if !strings.Contains(readData, "CN=2011") { //判断数据需要是 CN=2011
//		logrus.Warnf("非需数据%v", readData)
//		return false
//	}
//
//	return true
//}
