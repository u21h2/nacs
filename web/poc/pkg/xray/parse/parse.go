package parse

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
	"nacs/web/poc/errors"
	"nacs/web/poc/pkg/xray/structs"
)

var PocPool = sync.Pool{
	New: func() interface{} {
		return new(structs.Poc)
	},
}

func ParsePoc(filename string) (*structs.Poc, error) {
	poc := PocPool.Get().(*structs.Poc)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(f).Decode(poc)

	if err != nil {
		return nil, err
	}
	if poc.Name == "" {
		return nil, errors.Newf("Xray poc[%s] name can't be nil", filename)
	}

	if poc.Transport == "" {
		poc.Transport = "http"
	}
	return poc, nil
}
