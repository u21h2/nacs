package main

import (
	"fmt"
	"nacs/common"
	"nacs/discover"
	"nacs/nonweb"
	"nacs/parse"
	"nacs/probe"
	"nacs/utils/logger"
	"nacs/web/poc"
	"nacs/web/pocv1"
	"time"
)

func main() {
	startTime := time.Now()
	common.PrintBanner()
	parse.Flag(&common.InputInfo)
	parse.Parse(&common.InputInfo, &common.RunningInfo)
	probe.Probe()
	discover.Discover()
	pocv1.Poc()
	poc.Poc()
	nonweb.Service()
	logger.Info(fmt.Sprintf("Task finish, consumption of time: %s", time.Now().Sub(startTime)))
}
