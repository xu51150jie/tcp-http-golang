package response

import (
	"encoding/hex"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"tcp-http-golang/global"
	"tcp-http-golang/tcp"
)

//	func showOnline(r *gin.Engine, server *tcp.Server) {
//		r.GET("/showOnline", func(context *gin.Context) {
//			onlineMap := []map[string]interface{}{}
//			//查询onlineMap
//			server.MapLock.Lock()
//			for _, cli := range server.OnlineMap {
//				onlineMap = append(onlineMap, map[string]interface{}{"ipPort": cli.IpPort, "time": cli.Time})
//			}
//			server.MapLock.Unlock()
//
//			context.JSON(200, gin.H{
//				"count":     len(onlineMap),
//				"onlineMap": onlineMap,
//			})
//		})
//	}
func StartHandle(server *tcp.Server) {
	//handlePostRequest()
	//引用方法
	showClientOnline(server) // 显示客户端信息接口
	sendMsgToClient(server)  // 以http方式 发送消息到客户端

	var http_server_port = viper.GetString("HTTP_SERVER_PORT")
	err := http.ListenAndServe(http_server_port, nil) // 启动http接口服务
	global.Logger.Infof("HTTP服务器已经开启：%s", http_server_port)
	if err != nil {
		global.Logger.Errorf("启用http服务错误：%s", err)
	}
}

// 显示客户端信息接口
func showClientOnline(server *tcp.Server) {
	http.HandleFunc("/showClientOnline/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var onlineMap []map[string]interface{}
			//查询onlineMap
			server.MapLock.Lock()
			for _, cli := range server.OnlineMap {
				onlineMap = append(onlineMap, map[string]interface{}{"ipPort": cli.IpPort, "time": cli.Time})
			}
			server.MapLock.Unlock()

			//global.Logger.Info(onlineMap)
			jsonBytes, err := json.Marshal(onlineMap) // map转json
			if err != nil {
				global.Logger.Errorf("JSON转换出错:%s", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonBytes) // 数据流输出
			if err != nil {
				global.Logger.Errorf("返回出错%s", err)
			}

			return
		} else { // 无法处理该请求
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}

// 返回接口
type ResJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 以http方式 发送消息到客户端
func sendMsgToClient(server *tcp.Server) {
	http.HandleFunc("/sendMsgToClient/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			msgHttp := "" // 记录状态

			txquery := r.URL.Query()
			// 判断是否为空
			if txquery.Get("ipPort") == "" || txquery.Get("content") == "" {
				global.Logger.Infof("传参不全")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":400,"msg":"传参不全"}`)) // 数据流输出
				return
			}
			TX_IP_PORT := txquery.Get("ipPort")
			TX_CONTENT := txquery.Get("content")

			//查询onlineMap
			server.MapLock.Lock()
			cli, ok := server.OnlineMap[TX_IP_PORT]
			if !ok {
				msgHttp = TX_IP_PORT + "客户端不存在"
				global.Logger.Infof("客户端不存在 %v", txquery.Get("ipPort"))
			} else { // 执行发送数据

				var b []byte
				// 判断传输的内容是以 ASCII HEX
				if viper.GetString("SEND_DATA_TYPE") == "ASCII" {
					b = []byte(TX_CONTENT) // 字符串转[]byte
					cli.MsgChan <- b       // 发送到管道
					msgHttp = TX_IP_PORT + "发送成功"
					global.Logger.Infof("下发成功 %s -> %s", txquery.Get("ipPort"), txquery.Get("content"))
				} else {

					var err error
					b, err = hex.DecodeString(TX_CONTENT) // 十六进制字符串转十六进制[]byte
					if err != nil {
						msgHttp = TX_IP_PORT + "数据格式错误"
						global.Logger.Infof("数据格式错误 %v", err)
					} else {
						cli.MsgChan <- b // 发送到管道
						msgHttp = TX_IP_PORT + "发送成功"
						global.Logger.Infof("下发成功 %s -> %s", txquery.Get("ipPort"), txquery.Get("content"))
					}
				}

			}
			server.MapLock.Unlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			global.Logger.Infof(msgHttp)

			_, err1 := w.Write([]byte(`{"code":200," msg":"` + msgHttp + `"}`)) // 数据流输出
			if err1 != nil {
				global.Logger.Errorf("返回出错%s", err1)
			}

			return
		} else { // 无法处理该请求
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

}
