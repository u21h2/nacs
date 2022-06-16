package common

import (
	xray_structs "nacs/web/poc/pkg/xray/structs"
	"net/http"
)

var version = "0.0.2"

type InputInfoStruct struct {
	Host           string
	HostFile       string
	PortsAdd       string
	PortsOnly      string
	Mode           string
	OutputFileName string

	SkipHost string
	LiveTop  int
	Silent   bool
	NoSave   bool
	NoColor  bool
	LogLevel int
	Thread   int
	NoProbe  bool

	Proxy        string
	Timeout      int
	DiscoverMode string
	DiscoverType string
	OutJson      string

	CeyeKey    string
	CeyeDomain string

	NucleiPocPath string
	FscanPocPath  string
	PocRate       int
	PocDebug      bool
	PocThread     int
	PocTimeout    int

	BruteTimeout int

	NoNuclei bool

	NoPoc   bool
	NoBrute bool

	Command string
	SSHKey  string

	PasswordAdd string
	UsernameAdd string

	BruteThread      int
	BruteDebug       bool
	BruteSocks5Proxy string
	RedisFile        string
	RedisShell       string

	DirectUrl      string
	DirectUrlFile  string
	DirectUrlForce bool
}
type FileInfoStruct struct {
}

type RunningInfoStruct struct {
	OutputFileName string
	Hosts          []string
	Ports          []int
	// 不变
	NoProbe  bool
	HostFile string
	SkipHost string
	LiveTop  int
	Silent   bool
	NoSave   bool
	NoColor  bool
	LogLevel int
	Thread   int

	ProxySchema string
	ProxyHost   string
	ProxyProxy  string
	ProxyErr    error

	DiscoverTimeout int
	DiscoverMode    string
	DiscoverType    string

	OutJson string

	CeyeApi                  string
	CeyeDomain               string
	ReversePlatformType      xray_structs.ReverseType
	DnslogCNGetDomainRequest *http.Request
	DnslogCNGetRecordRequest *http.Request

	NucleiPocPath  string
	FscanPocPath   string
	PocTimeout     int
	PocRate        int
	PocDebug       bool
	PocThread      int
	PocSocks5Proxy string

	NoNuclei bool

	NoPoc   bool
	NoBrute bool

	Command          string
	SSHKey           string
	BruteTimeout     int
	BruteDebug       bool
	BruteThread      int
	BruteSocks5Proxy string
	RedisFile        string
	RedisShell       string

	DirectUse       bool
	DirectUrls      []map[string]interface{}
	DirectUrlsCache []string
	DirectUrlForce  bool
}

var InputInfo InputInfoStruct
var RunningInfo RunningInfoStruct

var AliveHosts []string
var DefaultPorts = []int{21, 22, 25, 53, 69, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 110, 135, 137, 138, 139, 143, 389, 443, 445, 554,
	587, 631, 800, 801, 808, 880, 888, 1000, 1024, 1025, 1080, 1099, 1389, 1433, 1521, 2000, 2001, 2222, 2601, 3306, 3307, 3308, 3309, 3388,
	3389, 3443, 5432, 5800, 5900, 6379, 7000, 7001, 7007, 7010, 7788, 8000, 8001, 8002, 8003, 8004, 8005, 8006, 8007, 8008, 8009, 8010, 8011,
	8030, 8060, 8070, 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 8090, 8091, 8092, 8093, 8094, 8095, 8096, 8097, 8098, 8099,
	8161, 8175, 8188, 8189, 8443, 8445, 8448, 8554, 8800, 8848, 8880, 8881, 8888, 8899, 8983, 8989, 9000, 9001, 9002, 9008, 9010, 9043, 9060,
	9080, 9081, 9082, 9083, 9084, 9085, 9086, 9087, 9088, 9089, 9090, 9091, 9092, 9093, 9094, 9095, 9096, 9097, 9099, 9200, 9300, 9443, 9448,
	9600, 9628, 9800, 9899, 9981, 9986, 9988, 9998, 9999, 11001, 11211, 13443, 15000, 20000, 27017, 33890, 45554, 49155, 49156, 50050, 61616}

var DefaultHeader = map[string]string{
	"Accept-Language": "zh,zh-TW;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6",
	"User-agent":      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36",
	"Cookie":          "rememberMe=int",
}

var DiscoverResultStruct = map[string]interface{}{
	"status":          "", // open
	"banner.byte":     "", //
	"banner.string":   "", //
	"protocol":        "", // ssh http https redis
	"type":            "", // tcp tls udp
	"host":            "", // IP地址
	"port":            "", //
	"uri":             "", //
	"note":            "", //
	"path":            "",
	"identify.bool":   false, // 是否匹配到服务
	"identify.string": "",    // 匹配到服务的关键词
}

var DiscoverResults = make([]map[string]interface{}, 0)

//var DiscoverResults []interface{}

//type DiscoverResultStruct struct {
//
//}
//
//var DiscoverResults = []DiscoverResultStruct

var Userdict = map[string][]string{
	"ftp":        {"ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"},
	"mysql":      {"root", "mysql"},
	"mssql":      {"sa", "sql"},
	"smb":        {"administrator", "admin", "guest"},
	"rdp":        {"administrator", "admin", "guest"},
	"postgresql": {"postgres", "admin"},
	"ssh":        {"root", "admin"},
	"mongodb":    {"root", "admin"},
	"oracle":     {"sys", "system", "admin", "test", "web", "orcl"},
}

var Passwords = []string{"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "{user}", "{user}1", "{user}111", "{user}123", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "{user}@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "888888", "88888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "Charge123", "Aa123456789"}
