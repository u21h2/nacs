package structs

import (
	"fmt"
	xray_structs "nacs/web/poc/pkg/xray/structs"
	"net/http"
	"strings"
	"time"
)

var (
	CeyeApi                  string
	CeyeDomain               string
	ReversePlatformType      xray_structs.ReverseType
	DnslogCNGetDomainRequest *http.Request
	DnslogCNGetRecordRequest *http.Request
)

func InitReversePlatform(api, domain string, timeout time.Duration) {
	if api != "" && domain != "" && strings.HasSuffix(domain, ".ceye.io") {
		CeyeApi = api
		CeyeDomain = domain
		fmt.Println(CeyeApi, CeyeDomain)
		ReversePlatformType = xray_structs.ReverseType_Ceye
	} else {
		ReversePlatformType = xray_structs.ReverseType_DnslogCN

		// 设置请求相关参数
		DnslogCNGetDomainRequest, _ = http.NewRequest("GET", "http://dnslog.cn/getdomain.php", nil)
		DnslogCNGetRecordRequest, _ = http.NewRequest("GET", "http://dnslog.cn/getrecords.php", nil)

	}
}
