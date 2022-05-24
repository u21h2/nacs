package judge

import (
	"nacs/utils/logger"
	"regexp"
)

func TlsRedisSsl(result map[string]interface{}) bool {
	var buff []byte
	buff, _ = result["banner.byte"].([]byte)
	ok, err := regexp.Match(`(^-ERR(.*)command|^-(.*).Redis)`, buff)
	if logger.DebugError(err) {
		return false
	}
	if ok {
		result["protocol"] = "redis-ssl"
		return true
	}
	return false
}
