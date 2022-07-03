package poc

import (
	"nacs/common"
	"nacs/utils/logger"
	"nacs/web/poc/internal/common/check"
	. "nacs/web/poc/internal/common/load"
	"nacs/web/poc/internal/common/output"
	nuclei_parse "nacs/web/poc/pkg/nuclei/parse"
	xray_requests "nacs/web/poc/pkg/xray/requests"
	"time"
)

func filterWebTargets() (targets []string) {
	for _, discoverResult := range common.DiscoverResults {
		if discoverResult["protocol"].(string) == "http" || discoverResult["protocol"].(string) == "https" {
			targets = append(targets, discoverResult["uri"].(string))
		}
	}

	return
}

func Poc() {
	if !common.RunningInfo.Nuclei {
		return
	}
	if common.RunningInfo.NoPoc {
		return
	}
	logger.Info("Start to send pocs to web services (nuclei type)")
	targets := filterWebTargets()

	var file = ""
	var json = false
	var xrayProxy = ""

	// 初始化dnslog平台
	//structs.InitReversePlatform(common.RunningInfo.CeyeApi, common.RunningInfo.CeyeDomain, time.Duration(common.RunningInfo.PocTimeout)*time.Second)
	//if common.RunningInfo.ReversePlatformType != xray_structs.ReverseType_Ceye {
	//	logger.Warning("No Ceye api, use dnslog.cn")
	//}

	// 初始化http客户端
	xray_requests.InitHttpClient(common.RunningInfo.PocThread, xrayProxy, time.Duration(common.RunningInfo.PocTimeout)*time.Second)

	// 初始化nuclei options
	nuclei_parse.InitExecuterOptions(common.RunningInfo.PocRate, common.RunningInfo.PocTimeout)

	// 加载poc
	xrayPocs, nucleiPocs := LoadPocs(common.RunningInfo.NucleiPocPath)

	// 计算xray的总发包量，初始化缓存
	xrayTotalRequests := 0
	totalTargets := len(targets)
	for _, poc := range xrayPocs {
		ruleLens := len(poc.Rules)
		// 额外需要缓存connectionID
		if poc.Transport == "tcp" || poc.Transport == "udp" {
			ruleLens += 1
		}
		xrayTotalRequests += totalTargets * ruleLens
	}
	if xrayTotalRequests == 0 {
		xrayTotalRequests = 1
	}
	xray_requests.InitCache(xrayTotalRequests)

	// 初始化输出
	outputChannel, outputWg := output.InitOutput(file, json)

	// 初始化check
	check.InitCheck(common.RunningInfo.PocThread, common.RunningInfo.PocRate, false)

	// check开始
	check.Start(targets, xrayPocs, nucleiPocs, outputChannel)
	check.Wait()

	// check结束
	close(outputChannel)
	check.End()
	outputWg.Wait()

}
