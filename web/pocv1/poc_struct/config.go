package poc_struct

type PocInfoStruct struct {
	Num        int
	Timeout    int64
	Proxy      string
	PocName    string
	PocDir     string
	Target     string
	TargetFile string
	RawFile    string
	Cookie     string
	ForceSSL   bool
	ApiKey     string
	CeyeDomain string
}

var PocInfo PocInfoStruct
