package output

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"nacs/common"
	"nacs/utils/logger"
	"os"
	"strings"
)

var JsonData []string

func JsonOutput(result map[string]interface{}, operate string) {

	if operate == "save" {
		result["banner.byte"] = hex.EncodeToString(result["banner.byte"].([]byte))
		byteR, err := json.Marshal(result)
		if logger.DebugError(err) {
			logger.Error("json.Marshal failed to " + logger.Red("flag.OutJson"))
		}
		strR := string(byteR)
		JsonData = append(JsonData, strR)
	}
	if operate == "write" {
		_, err := os.Stat(common.RunningInfo.OutJson)
		if err == nil {
			logger.Error(logger.LightRed(fmt.Sprintf("json file %s already exists", common.RunningInfo.OutJson)))
		}
		jsonFile, err := os.OpenFile(common.RunningInfo.OutJson, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if logger.DebugError(err) {
			logger.Error(fmt.Sprintf("Failed to open file %s", logger.Red(common.RunningInfo.OutJson)))
		}
		byteContent := strings.Join(JsonData, ", ")
		_, err = jsonFile.Write([]byte(fmt.Sprintf("[%s]\r\n", byteContent)))
		if logger.DebugError(err) {
			logger.Error(fmt.Sprintf("Target %s write failed", logger.Red(result["uri"])))
		}
		err = jsonFile.Close()
		if logger.DebugError(err) {
			logger.Error(fmt.Sprintf("Close file %s exception", logger.Red(common.RunningInfo.OutJson)))
		} else {
			logger.Info("The identification json results are saved in " + logger.Yellow(common.RunningInfo.OutJson))
		}
	}
}
