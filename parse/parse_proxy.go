package parse

import (
	"fmt"
	"nacs/common"
	"nacs/utils/logger"
	"net/url"
)

func ProxyParse() {
	if common.InputInfo.Proxy == "" {
		common.RunningInfo.ProxySchema = ""
		common.RunningInfo.ProxyHost = ""
		common.RunningInfo.ProxyProxy = ""
		common.RunningInfo.ProxyErr = nil
		return
	}

	u, err := url.Parse(common.InputInfo.Proxy)
	if logger.DebugError(err) {
		logger.Fatal(fmt.Sprintf("The proxy address %s is formatted incorrectly", logger.LightRed(common.InputInfo.Proxy)))
		common.RunningInfo.ProxySchema = ""
		common.RunningInfo.ProxyHost = ""
		common.RunningInfo.ProxyProxy = ""
		common.RunningInfo.ProxyErr = nil
	}
	if u.Scheme == "http" || u.Scheme == "socks5" {
		common.RunningInfo.ProxySchema = u.Scheme
		common.RunningInfo.ProxyHost = u.Host
		common.RunningInfo.ProxyProxy = common.InputInfo.Proxy
		common.RunningInfo.ProxyErr = nil
		return
	}
	logger.Fatal(logger.Red("Unsupported proxy protocol"))
}
