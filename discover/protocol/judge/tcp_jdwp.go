package judge

import (
	"bytes"
	"nacs/utils/logger"
	"strings"

	"nacs/discover/parse"
	"nacs/discover/proxy"
)

func TcpJdwp(result map[string]interface{}, Args map[string]interface{}) bool {
	timeout := Args["Timeout"].(int)
	host := result["host"].(string)
	port := result["port"].(int)

	conn, err := proxy.ConnProxyTcp(host, port, timeout)
	if logger.DebugError(err) {
		return false
	}

	msg := "JDWP-Handshake"
	_, err = conn.Write([]byte(msg))
	if logger.DebugError(err) {
		return false
	}

	reply := make([]byte, 256)
	_, _ = conn.Read(reply)
	if conn != nil {
		_ = conn.Close()
	}

	var buffer [256]byte
	if bytes.Equal(reply[:], buffer[:]) {
		return false
	}
	if strings.Contains(string(reply), "JDWP-Handshake") == false {
		return false
	}
	result["protocol"] = "jdwp"
	result["banner.string"] = parse.ByteToStringParse2(reply[0:16])
	result["banner.byte"] = reply
	return true
}
