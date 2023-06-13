## 一个可以正向从socket接受往http发送数据，并可以以http接受反向往socket发送数据

####查看在线客户端数量
http://ip:port/showClientOnline

####以http方式 发送消息到客户端
http://ip:port/sendMsgToClient?ipPort=192.168.0.2:11&content=abc123


---
**1.编译方式** 

Mac下编译Linux可执行程序：
`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aaaaa main.go`  

Mac下编译Windows平台的64位可执行程序
`CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o aaaa main.go` 

Windows下编译Linux平台的64位可执行程序：注意使用管理员启动命令提示符（不是PowerShell）
`SET CGO_ENABLED=0`
`SET GOOS=linux`
`SET GOARCH=amd64`
`go build -o aaaa main.go`  废气
---

**2.运行程序** 注意config.yml文件中配置修改
`./run.sh  start | restart | stop`


