package parse

import (
	"nacs/common"
	xray_structs "nacs/web/poc/pkg/xray/structs"
	"net/http"
	"strings"
)

func ParseReversePlatform(api, domain string) {
	if api != "" && domain != "" && strings.HasSuffix(domain, ".ceye.io") {
		common.RunningInfo.CeyeApi = api
		common.RunningInfo.CeyeDomain = domain
		common.RunningInfo.ReversePlatformType = xray_structs.ReverseType_Ceye
	} else {
		common.RunningInfo.ReversePlatformType = xray_structs.ReverseType_DnslogCN

		// 设置请求相关参数
		common.RunningInfo.DnslogCNGetDomainRequest, _ = http.NewRequest("GET", "http://dnslog.cn/getdomain.php", nil)
		common.RunningInfo.DnslogCNGetRecordRequest, _ = http.NewRequest("GET", "http://dnslog.cn/getrecords.php", nil)

	}
}
