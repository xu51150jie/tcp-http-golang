package request

import (
	"fmt"
	"github.com/spf13/viper"
	"tcp-http-golang/global"
	"tcp-http-golang/utils"
)

// 上报数据得到接口
func SendDataHttp(readData string, ipPort string) {
	// 判断配置文件是否开启http发送
	//_sd := viper.GetStringMapString("SEND_DATA")
	if viper.GetBool("SEND_DATA.status") {
		if len(readData) < 2 { //屏蔽短内容数据 正常数据长度应大于2以上
			global.Logger.Infof("非需数据%v", readData)
		} else {
			/**
			 ** 这里根据接口需求发送对应格式的body
			 **/
			_bodyStr := fmt.Sprintf("{\"data\":\"%s\"}", isStrKG(readData)) // 单发送内容
			//_bodyStr := fmt.Sprintf("{\"ipAddr\":\"%s\",\"count\":\"%s\"}", ipPort, isStrKG(readData)) // 发送IP和内容

			body := utils.HttpPostJson(_bodyStr, viper.GetString("SEND_DATA.url")) //发送接受端http接口
			global.Logger.Infof("上报接口返回数据：%v %v ", ipPort, string(body))
		}
	}
}

// 字符串每隔个字符添加一个空格
func isStrKG(str string) string {
	var _data string
	//   如果HEX格式 则需要添加字符串空格
	if viper.GetString("SEND_DATA_TYPE") == "HEX" {
		runes := []rune(str)
		var result []rune
		for i, r := range runes {
			result = append(result, r)
			if (i+1)%2 == 0 {
				result = append(result, ' ')
			}
		}
		_data = string(result)
	}
	return _data
}
