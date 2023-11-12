package err

var MsgFlags = map[int]string{
	ConfigFileNotFound: "config.yaml未发现",
	ConfigFileReadErr:  "文件读取错误",
	ConfigReadErr:      "配置文件配置参数发现未知错误",
	ConfigReadSuccess:  "%s 读取成功",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
