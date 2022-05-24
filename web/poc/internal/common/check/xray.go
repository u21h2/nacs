package check

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"nacs/web/poc/internal/common/errors"
	"nacs/web/poc/pkg/xray/cel"
	"nacs/web/poc/pkg/xray/requests"
	xray_structs "nacs/web/poc/pkg/xray/structs"
	"nacs/web/poc/utils"

	"github.com/google/cel-go/checker/decls"
	"gopkg.in/yaml.v2"
)

var (
	BodyBufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	BodyPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}
	VariableMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{})
		},
	}
)

type RequestFuncType func(ruleName string, rule xray_structs.Rule) error

func executeXrayPoc(oReq *http.Request, target string, poc *xray_structs.Poc) (isVul bool, err error) {
	isVul = false

	var (
		milliseconds int64
		tcpudpType   string = ""

		request       *http.Request
		response      *http.Response
		oProtoRequest *xray_structs.Request
		protoRequest  *xray_structs.Request
		protoResponse *xray_structs.Response
		variableMap   map[string]interface{} = VariableMapPool.Get().(map[string]interface{})

		oReqUrlString string

		requestFunc cel.RequestFuncType
	)

	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "Run Xray Poc[%s] error", poc.Name)
			isVul = false
		}
	}()
	// 回收
	defer func() {
		if protoRequest != nil {
			requests.PutUrlType(protoRequest.Url)
			requests.PutRequest(protoRequest)

		}
		if oProtoRequest != nil {
			requests.PutUrlType(oProtoRequest.Url)
			requests.PutRequest(oProtoRequest)

		}
		if protoResponse != nil {
			requests.PutUrlType(protoResponse.Url)
			if protoResponse.Conn != nil {
				requests.PutAddrType(protoResponse.Conn.Source)
				requests.PutAddrType(protoResponse.Conn.Destination)
				requests.PutConnectInfo(protoResponse.Conn)
			}
			requests.PutResponse(protoResponse)
		}

		for _, v := range variableMap {
			switch v.(type) {
			case *xray_structs.Reverse:
				cel.PutReverse(v)
			default:
			}
		}
		VariableMapPool.Put(variableMap)
	}()

	// 初始赋值
	if oReq != nil {
		oReqUrlString = oReq.URL.String()
	}
	utils.DebugF("Run Xray Poc[%s] for %s", poc.Name, target)

	// 设置原始请求变量
	oProtoRequest, _ = requests.ParseHttpRequest(oReq)
	variableMap["request"] = oProtoRequest

	// 判断transport，如果不合法则跳过
	transport := poc.Transport
	if transport == "tcp" || transport == "udp" {
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			utils.InfoF("Invalid target[%s], skip", target)
			return
		}
	} else {
		_, err = url.ParseRequestURI(target)
		if err != nil {
			utils.InfoF("Invalid target[%s], skip", target)
			return
		}
	}

	// 初始化cel-go环境，并在函数返回时回收
	c := cel.NewEnvOption()
	defer cel.PutCustomLib(c)

	env, err := cel.NewEnv(&c)
	if err != nil {
		wrappedErr := errors.Wrap(err, "Environment creation error")
		utils.ErrorP(wrappedErr)
		return false, err
	}

	// 请求中的全局变量

	// 定义渲染函数
	render := func(v string) string {
		for k1, v1 := range variableMap {
			_, isMap := v1.(map[string]string)
			if isMap {
				continue
			}
			v1Value := fmt.Sprintf("%v", v1)
			t := "{{" + k1 + "}}"
			if !strings.Contains(v, t) {
				continue
			}
			v = strings.ReplaceAll(v, t, v1Value)
		}
		return v
	}
	ReCreateEnv := func() error {
		env, err = cel.NewEnv(&c)
		if err != nil {
			wrappedErr := errors.Newf(errors.EnvInitializationError, "Environment re-creation error: %v", err)
			return wrappedErr
		}
		return nil
	}

	// 定义evaluateUpdateVariableMap
	evaluateUpdateVariableMap := func(set yaml.MapSlice) {
		for _, item := range set {
			k, expression := item.Key.(string), item.Value.(string)
			// ? 需要重新生成一遍环境，否则之前增加的变量定义不生效
			if err := ReCreateEnv(); err != nil {
				utils.ErrorP(err)
			}

			out, err := cel.Evaluate(env, expression, variableMap)
			if err != nil {
				wrappedErr := errors.Wrapf(err, "Evalaute expression error: %s", expression)
				utils.ErrorP(wrappedErr)
				continue
			}

			// 设置variableMap并且更新CompileOption
			switch value := out.Value().(type) {
			case *xray_structs.UrlType:
				if _, ok := variableMap[k]; !ok {
					c.UpdateCompileOption(k, cel.UrlTypeType)
				}
				variableMap[k] = cel.UrlTypeToString(value)
			case *xray_structs.Reverse:
				if _, ok := variableMap[k]; !ok {
					c.UpdateCompileOption(k, cel.ReverseType)
				}
				variableMap[k] = value
			case int64:
				if _, ok := variableMap[k]; !ok {
					c.UpdateCompileOption(k, decls.Int)
				}
				variableMap[k] = int(value)
			case map[string]string:
				if _, ok := variableMap[k]; !ok {
					c.UpdateCompileOption(k, cel.StrStrMapType)
				}
				variableMap[k] = value
			default:
				if _, ok := variableMap[k]; !ok {
					c.UpdateCompileOption(k, decls.String)
				}
				variableMap[k] = value
			}
		}
		// ? 最后再生成一遍环境，否则之前增加的变量定义不生效
		if err := ReCreateEnv(); err != nil {
			utils.ErrorP(err)
		}
	}

	// 处理set
	evaluateUpdateVariableMap(poc.Set)

	// 渲染detail
	detail := &poc.Detail
	detail.Author = render(detail.Author)
	for k, v := range poc.Detail.Links {
		detail.Links[k] = render(v)
	}
	fingerPrint := &detail.FingerPrint
	for _, info := range fingerPrint.Infos {
		info.ID = render(info.ID)
		info.Name = render(info.Name)
		info.Version = render(info.Version)
		info.Type = render(info.Type)
	}
	fingerPrint.HostInfo.Hostname = render(fingerPrint.HostInfo.Hostname)
	vulnerability := &detail.Vulnerability
	vulnerability.ID = render(vulnerability.ID)
	vulnerability.Match = render(vulnerability.Match)

	// transport=http: request处理
	HttpRequestInvoke := func(rule xray_structs.Rule) error {
		var (
			ok               bool
			err              error
			ruleReq          xray_structs.RuleRequest = rule.Request
			rawHeaderBuilder strings.Builder
		)

		// 渲染请求头，请求路径和请求体
		for k, v := range ruleReq.Headers {
			ruleReq.Headers[k] = render(v)
		}
		ruleReq.Path = render(strings.TrimSpace(ruleReq.Path))
		ruleReq.Body = render(strings.TrimSpace(ruleReq.Body))

		// 尝试获取缓存
		if request, protoRequest, protoResponse, ok = requests.XrayGetHttpRequestCache(&ruleReq); !ok || !rule.Request.Cache {
			// 获取protoRequest
			protoRequest, err = requests.ParseHttpRequest(oReq)
			if err != nil {
				wrappedErr := errors.Wrapf(err, "Run poc[%v] parse request error", poc.Name)
				return wrappedErr
			}

			// 处理Path
			if strings.HasPrefix(ruleReq.Path, "/") {
				protoRequest.Url.Path = strings.Trim(oReq.URL.Path, "/") + "/" + ruleReq.Path[1:]
			} else if strings.HasPrefix(ruleReq.Path, "^") {
				protoRequest.Url.Path = "/" + ruleReq.Path[1:]
			}

			if !strings.HasPrefix(protoRequest.Url.Path, "/") {
				protoRequest.Url.Path = "/" + protoRequest.Url.Path
			}

			// 某些poc没有区分path和query，需要处理
			protoRequest.Url.Path = strings.ReplaceAll(protoRequest.Url.Path, " ", "%20")
			protoRequest.Url.Path = strings.ReplaceAll(protoRequest.Url.Path, "+", "%20")

			// 克隆请求对象
			request, err = http.NewRequest(ruleReq.Method, fmt.Sprintf("%s://%s%s", protoRequest.Url.Scheme, protoRequest.Url.Host, protoRequest.Url.Path), strings.NewReader(ruleReq.Body))
			if err != nil {
				return err
			}

			// 处理请求头
			request.Header = oReq.Header.Clone()
			for k, v := range ruleReq.Headers {
				request.Header.Set(k, v)
				rawHeaderBuilder.WriteString(k)
				rawHeaderBuilder.WriteString(": ")
				rawHeaderBuilder.WriteString(v)
				rawHeaderBuilder.WriteString("\n")
			}

			protoRequest.RawHeader = []byte(strings.Trim(rawHeaderBuilder.String(), "\n"))

			// 额外处理protoRequest.Raw
			protoRequest.Raw, _ = httputil.DumpRequestOut(request, true)

			// 发起请求
			response, milliseconds, err = requests.DoRequest(request, ruleReq.FollowRedirects)
			if err != nil {
				return err
			}

			// 获取protoResponse
			protoResponse, err = requests.ParseHttpResponse(response, milliseconds)
			if err != nil {
				wrappedErr := errors.Wrapf(err, "Run poc[%s] parse response error", poc.Name)
				return wrappedErr
			}

			// 设置缓存
			requests.XraySetHttpRequestCache(&ruleReq, request, protoRequest, protoResponse)

		} else {
			utils.DebugF("Hit http request cache[%s%s]", oReqUrlString, ruleReq.Path)
		}

		return nil
	}

	// transport=tcp/udp: request处理
	TCPUDPRequestInvoke := func(rule xray_structs.Rule) error {
		var (
			tcpudpTypeUpper = strings.ToUpper(tcpudpType)
			buffer          = BodyBufPool.Get().([]byte)

			content             = rule.Request.Content
			connectionID string = rule.Request.ConnectionID
			conn         net.Conn
			connCache    *net.Conn
			responseRaw  []byte
			readTimeout  int

			ok  bool
			err error
		)
		defer BodyBufPool.Put(buffer)

		// 获取response缓存
		if responseRaw, protoResponse, ok = requests.XrayGetTcpUdpResponseCache(rule.Request.Content); !ok || !rule.Request.Cache {
			responseRaw = BodyPool.Get().([]byte)
			defer BodyPool.Put(responseRaw)

			// 获取connectionID缓存
			if connCache, ok = requests.XrayGetTcpUdpConnectionCache(connectionID); !ok {
				// 处理timeout
				readTimeout, err = strconv.Atoi(rule.Request.ReadTimeout)
				if err != nil {
					wrappedErr := errors.Wrapf(err, "Parse read_timeout[%s] to int  error", rule.Request.ReadTimeout)
					return wrappedErr
				}

				// 发起连接
				conn, err = net.Dial(tcpudpType, target)
				if err != nil {
					wrappedErr := errors.Wrapf(err, "%s connect to target[%s] error", tcpudpTypeUpper, target)
					return wrappedErr
				}

				// 设置读取超时
				err := conn.SetReadDeadline(time.Now().Add(time.Duration(readTimeout) * time.Second))
				if err != nil {
					wrappedErr := errors.Wrapf(err, "Set read_timeout[%d] error", tcpudpTypeUpper, readTimeout)
					return wrappedErr
				}

				// 设置连接缓存
				requests.XraySetTcpUdpConnectionCache(connectionID, &conn)
			} else {
				conn = *connCache
				utils.DebugF("Hit connection_id cache[%s]", connectionID)
			}

			// 获取protoRequest
			protoRequest, _ = requests.ParseTCPUDPRequest([]byte(content))

			// 发送数据
			_, err = conn.Write([]byte(content))
			if err != nil {
				wrappedErr := errors.Wrapf(err, "%s[%s] write error", tcpudpTypeUpper, connectionID)
				return wrappedErr
			}

			// 接收数据
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					if err == io.EOF {
					} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					} else {
						wrappedErr := errors.Wrapf(err, "%s[%s] read error", tcpudpTypeUpper, connectionID)
						return wrappedErr
					}
					break
				}
				responseRaw = append(responseRaw, buffer[:n]...)
			}

			// 获取protoResponse
			protoResponse, _ = requests.ParseTCPUDPResponse(responseRaw, &conn, tcpudpType)

			// 设置响应缓存
			requests.XraySetTcpUdpResponseCache(content, responseRaw, protoResponse)

		} else {
			utils.DebugF("Hit tcp/udp request cache[%s]", responseRaw)
		}

		return nil
	}

	// reqeusts总处理
	RequestInvoke := func(requestFunc cel.RequestFuncType, ruleName string, rule xray_structs.Rule) (bool, error) {
		var (
			flag bool
			ok   bool
			err  error
		)
		err = requestFunc(rule)
		if err != nil {
			return false, err
		}

		variableMap["request"] = protoRequest
		variableMap["response"] = protoResponse

		utils.DebugF("raw requests: \n%#s", string(protoRequest.Raw))
		utils.DebugF("raw response: \n%#s", string(protoResponse.Raw))

		// 执行表达式
		out, err := cel.Evaluate(env, rule.Expression, variableMap)

		if err != nil {
			wrappedErr := errors.Wrapf(err, "Evalute rule[%s] expression of [%s] error: %s", ruleName, "", rule.Expression)
			return false, wrappedErr
		}

		// 判断表达式结果
		flag, ok = out.Value().(bool)
		if !ok {
			flag = false
		}

		// 处理output
		evaluateUpdateVariableMap(rule.Output)

		return flag, nil
	}

	// 判断transport类型，设置requestInvoke
	if poc.Transport == "tcp" {
		tcpudpType = "tcp"
		requestFunc = TCPUDPRequestInvoke
	} else if poc.Transport == "udp" {
		tcpudpType = "udp"
		requestFunc = TCPUDPRequestInvoke
	} else {
		requestFunc = HttpRequestInvoke
	}

	ruleSlice := poc.Rules
	// 提前定义名为ruleName的函数
	for _, ruleItem := range ruleSlice {
		c.DefineRuleFunction(requestFunc, ruleItem.Key, ruleItem.Value, RequestInvoke)
	}

	// ? 最后再生成一遍环境，否则之前增加的变量定义不生效
	if err := ReCreateEnv(); err != nil {
		utils.ErrorP(err)
	}

	// 执行rule 并判断poc总体表达式结果
	run := func() (bool, error) {
		successVal, err := cel.Evaluate(env, poc.Expression, variableMap)
		if err != nil {
			wrappedErr := errors.Wrapf(err, "Evalute poc[%s] expression error: %s", poc.Name, poc.Expression)
			return false, wrappedErr
		}

		isVul, ok := successVal.Value().(bool)
		if !ok {
			isVul = false
		}

		return isVul, nil
	}

	// 如果没设置payload，则直接评估rules并返回
	if len(poc.Payloads.Payloads) == 0 {
		return run()
	}

	// 如果设置了payload，则遍历执行
	isVul = false

	for _, setMapVal := range poc.Payloads.Payloads {
		setMap := setMapVal.Value.(yaml.MapSlice)
		evaluateUpdateVariableMap(setMap)
		isVul, err = run()
		if err != nil {
			return false, err
		}

		if isVul && !poc.Payloads.Continue {
			return isVul, nil
		}
	}
	return isVul, nil
}
