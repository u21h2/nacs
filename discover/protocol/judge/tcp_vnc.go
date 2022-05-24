package judge

import (
	"nacs/utils/logger"
	"regexp"
)

func TcpVNC(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^RFB \d`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "vnc"
		return true
	}
	return false
}
