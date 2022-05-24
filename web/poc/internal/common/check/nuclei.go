package check

import (
	"nacs/web/poc/internal/common/errors"
	nuclei_structs "nacs/web/poc/pkg/nuclei/structs"
	"nacs/web/poc/utils"

	"github.com/projectdiscovery/nuclei/v2/pkg/output"
)

func executeNucleiPoc(target string, poc *nuclei_structs.Poc) (results []*output.ResultEvent, isVul bool, err error) {
	isVul = false

	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "Run Nuclei Poc[%s] error", poc.ID)
			isVul = false
		}
	}()

	utils.DebugF("Run Nuclei Poc[%s] for %s", poc.Info.Name, target)

	e := poc.Executer
	results = make([]*output.ResultEvent, 0, e.Requests())

	err = e.ExecuteWithResults(target, func(result *output.InternalWrappedEvent) {
		if len(result.Results) > 0 {
			isVul = true
		}
		results = append(results, result.Results...)
	})

	if len(results) == 0 {
		results = append(results, &output.ResultEvent{TemplateID: poc.ID, Matched: target})
	}
	return results, isVul, err
}
