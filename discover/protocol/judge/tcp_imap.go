package judge

import (
	"nacs/utils/logger"
	"regexp"
)

func TcpIMAP(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`^* OK`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "imap"
		return true
	}
	return false
}
