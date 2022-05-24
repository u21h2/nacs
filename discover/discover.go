package discover

import (
	"nacs/common"
	"nacs/discover/parse"
	"nacs/discover/protocol"
	"nacs/utils/logger"
	"strconv"
	"strings"
	"sync"
)

func Discover() {

	thread := common.RunningInfo.Thread
	ports := common.DefaultPorts
	actualHosts := common.AliveHosts

	wg := &sync.WaitGroup{}

	Args := make(map[string]interface{})
	Args["Timeout"] = common.RunningInfo.DiscoverTimeout
	Args["Mode"] = common.RunningInfo.DiscoverMode
	Args["Type"] = common.RunningInfo.DiscoverType

	logger.Info("Start to discover the ports")
	intSyncThread := 0
	intAll := 0
	intIde := 0
	for _, host := range actualHosts {
		for _, port := range ports {
			wg.Add(1)
			intSyncThread++
			go func(host string, port int, Args map[string]interface{}) {
				res := protocol.Discover(host, port, Args)
				if res["status"].(string) == "open" {
					intAll++
					parse.VerboseParse(res)
					//output.JsonOutput(res, "save")
					common.DiscoverResults = append(common.DiscoverResults, res)
					if strings.Contains(res["uri"].(string), "://") {
						intIde++
					}
				}
				wg.Done()
			}(host, port, Args)
			if intSyncThread >= thread {
				intSyncThread = 0
				wg.Wait()
			}
		}
	}
	wg.Wait()
	logger.Info(logger.LightGreen("A total of ") +
		logger.White(strconv.Itoa(intAll)) +
		logger.LightGreen(" targets, the rule base hits ") +
		logger.White(strconv.Itoa(intIde)) +
		logger.LightGreen(" targets"))
	//output.JsonOutput(make(map[string]interface{}), "write")
}
