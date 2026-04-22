package kenan

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/louismax/kenan/kTool"
	"github.com/rcrowley/go-metrics"
	etcdCli "github.com/rpcxio/rpcx-etcd/client"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/protocol"
	rpcPlugin "github.com/smallnest/rpcx/serverplugin"
)

var (
	etcdBasePath *string //ETCD 注册路径前缀
	etcdURL      *string //ETCD 注册中心
	svrName      *string //微服务节点前缀
	RpcxSvr *server.Server //当前服务
)

// RpcMapSync 线程安全的Rpc客户端实例map
var RpcMapSync sync.Map

type Args struct {
	HttpHeader http.Header
	Session    map[string]interface{}
	Data       interface{}
}

type Reply struct {
	Success      bool
	Code         int
	ErrorMessage string
	Data         interface{}
}

func RPCXInit() {
	etcdURL = flag.String("info", GetConfigDefault("etcdUrl", "127.0.0.1:2379"), "ETCD URL")
	etcdBasePath = flag.String("base", "/"+GetConfigDefault("EtcdBasePath", "RPCX_Root"), "prefix path")
	svrName = flag.String("RPC_"+settings.ServerTag, settings.ServerTag, "Service Node - "+settings.ServerTag)
}

// RpcNewClient RPC客户端注册
func RpcNewClient(node string) (client.XClient, error) {
	var serverNode *string
	if flag.Lookup("RPC_"+node) != nil {
		temp := (*flag.Lookup("RPC_" + node)).Value.(flag.Getter).Get().(string)
		serverNode = &temp
	} else {
		serverNode = flag.String("RPC_"+node, node, "Service Node - "+node)
	}

	d, err := etcdCli.NewEtcdV3Discovery(*etcdBasePath, *serverNode, []string{*etcdURL}, true, nil)
	if err != nil {
		return nil, err
	}
	option := client.DefaultOption
	option.ConnectTimeout = 10 * time.Second //建立连接的超时时间

	option.Heartbeat = true                    //启用心跳
	option.HeartbeatInterval = 5 * time.Second //心跳检测间隔，5秒一次

	// 内网环境可以设置较短的检测间隔
	option.TCPKeepAlivePeriod = 120 * time.Second // 240秒探测

	option.IdleTimeout = 0 // 禁用空闲超时

	option.Retries = 3 // 失败重试次数
	//断路器实现
	option.GenBreaker = func() client.Breaker {
		return client.NewConsecCircuitBreaker(uint64(GetConfigDefaultByInt("BreakerFuseCount", 30)), time.Duration(GetConfigDefaultByInt("BreakerRecoveryTimes", 10))*time.Second)
	}

	return client.NewXClient(*serverNode, client.Failover, client.RandomSelect, d, option), nil
}

// RegisterMS 微服务注册 port注册端口
func RegisterMS(svrArity interface{}, port string, limits ...int64) error {
	flag.Parse()
	addr := flag.String("addr", kTool.GetInternal()+":"+port, "service address")
	LogInfo("微服务注册在:tcp@%s", *addr)
	//注册服务
	RpcxSvr = server.NewServer()
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*etcdURL},
		BasePath:       *etcdBasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Second,
	}
	err := r.Start()
	if err != nil {
		return err
	}
	RpcxSvr.Plugins.Add(r)
	RpcxSvr.Plugins.Add(&ConnectionPlugin{})

	if len(limits) == 2 && time.Duration(limits[0]) > 0 && limits[1] > 0 {
		//限流器实现
		rateLimiter := rpcPlugin.NewReqRateLimitingPlugin(time.Second, 1, true)
		RpcxSvr.Plugins.Add(rateLimiter)
	}

	err = RpcxSvr.RegisterName(*svrName, svrArity, "")
	if err != nil {
		return err
	}
	err = RpcxSvr.Serve("tcp", *addr)
	if err != nil {
		return err
	}
	return nil
}

func (args *Args) CallService(node, fun string, data interface{}) (*Reply, error) {
	var err error
	var nextArgs Args
	nextArgs.HttpHeader = args.HttpHeader
	nextArgs.Session = args.Session
	nextArgs.Data = data

	var mz client.XClient

	if v, ok := RpcMapSync.Load(node); ok {
		mz = *v.(*client.XClient)
	} else {
		newMz, err := RpcNewClient(node)
		if err != nil {
			return nil, err
		}
		RpcMapSync.Store(node, &newMz)
		mz = newMz
	}
	var replyX Reply

	if err = mz.Call(context.Background(), fun, &nextArgs, &replyX); err != nil {
		fmt.Printf(">>>>>>>>>>>>>>>>>err:%+v,尝试重新建立连接！\n", err)
		if e := mz.Close(); e != nil { //先尝试关闭之前的连接!
			fmt.Println(e)
		}
		//重新建立一个连接
		mz, err = RpcNewClient(node)
		if err != nil {
			return nil, err
		}
		RpcMapSync.Store(node, &mz)
		if err = mz.Call(context.Background(), fun, &nextArgs, &replyX); err != nil {
			return nil, err
		}
	}

	return &replyX, nil
}

// ConnectionPlugin RPC连接插件
type ConnectionPlugin struct {
}

// HandleConnAccept Rpc客户端连接
func (p *ConnectionPlugin) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	LogInfo("RPCX服务间通讯已成功连接,RemoteAddr:%s,LocalAddr:%s", conn.RemoteAddr().String(), conn.LocalAddr().String())
	return conn, true
}

// HandleConnClose Rpc客户端关闭
func (p *ConnectionPlugin) HandleConnClose(conn net.Conn) bool {

	LogInfo("RPCX服务间通讯已关闭,RemoteAddr:%s,LocalAddr:%s", conn.RemoteAddr().String(), conn.LocalAddr().String())
	return true
}

func (args *Args) ReturnErrorNative(ack int, msg string, reply *Reply) error {
	reply.Success = false
	reply.Code = ack
	reply.ErrorMessage = msg

	return nil
}

func (args *Args) ReturnSuccessNative(d interface{}, reply *Reply, isByte ...bool) error {
	reply.Success = true
	if len(isByte) > 0 && isByte[0] {
		reply.Data = args.UseComplicated(d)
	} else {
		reply.Data = d
	}
	return nil
}

func (args *Args) ReturnError(ack int, msg string, reply *Reply) error {
	var rtn map[string]interface{}
	rtn = make(map[string]interface{})
	rtn["Success"] = false
	rtn["Result"] = ack
	rtn["ErrorMessage"] = msg
	data, err := json.Marshal(rtn)
	if err != nil {
		return err
	}
	reply.Data = data
	return nil
}

func (args *Args) ReturnSuccess(d interface{}, reply *Reply) error {
	var rtn map[string]interface{}
	rtn = make(map[string]interface{})
	rtn["Success"] = true
	if d != nil {
		rtn["Result"] = d
	}
	data, err := json.Marshal(rtn)
	if err != nil {
		return err
	}
	reply.Data = data
	return nil
}

// UseComplicated 复杂对象转为[]byte
func (args *Args) UseComplicated(body interface{}) []byte {
	jsons, errs := json.Marshal(body) //转换成JSON返回的是byte[]
	if errs != nil {
		LogError("json转移失败,err:%+v", errs)
	}
	return jsons
}

// ComplexAnalysis 解析指定的复杂Data到内存对象
func (args *Args) ComplexAnalysis(body interface{}, headerObj interface{}) error {
	if reflect.TypeOf(headerObj).Kind() != reflect.Ptr {
		return errors.New("headerObj必须是一个指针对象")
	}
	return json.Unmarshal(body.([]byte), &headerObj)
}

// ComplexAnalysisSelf 解析来自调用方返回的复杂Data到内存对象
func (args *Args) ComplexAnalysisSelf(headerObj interface{}) error {
	if reflect.TypeOf(headerObj).Kind() != reflect.Ptr {
		return errors.New("headerObj必须是一个指针对象")
	}
	return json.Unmarshal(args.Data.([]byte), &headerObj)
}

// ComplexAnalysis 解析Call返回的复杂Data到内存对象
func (reply *Reply) ComplexAnalysis(headerObj interface{}) error {
	if reflect.TypeOf(headerObj).Kind() != reflect.Ptr {
		return errors.New("headerObj必须是一个指针对象")
	}
	return json.Unmarshal(reply.Data.([]byte), &headerObj)
}

// RpcNewBidirectionalClient RPC双向客户端注册
func RpcNewBidirectionalClient(node string,msgChan chan<- *protocol.Message) (client.XClient, error) {
	var serverNode *string
	if flag.Lookup("RPC_"+node) != nil {
		temp := (*flag.Lookup("RPC_" + node)).Value.(flag.Getter).Get().(string)
		serverNode = &temp
	} else {
		serverNode = flag.String("RPC_"+node, node, "Service Node - "+node)
	}

	d, err := etcdCli.NewEtcdV3Discovery(*etcdBasePath, *serverNode, []string{*etcdURL}, true, nil)
	if err != nil {
		return nil, err
	}
	option := client.DefaultOption
	option.ConnectTimeout = 10 * time.Second //建立连接的超时时间

	option.Heartbeat = true                    //启用心跳
	option.HeartbeatInterval = 5 * time.Second //心跳检测间隔，5秒一次

	// 内网环境可以设置较短的检测间隔
	option.TCPKeepAlivePeriod = 120 * time.Second // 240秒探测

	option.IdleTimeout = 0 // 禁用空闲超时

	option.Retries = 3 // 失败重试次数
	//断路器实现
	option.GenBreaker = func() client.Breaker {
		return client.NewConsecCircuitBreaker(uint64(GetConfigDefaultByInt("BreakerFuseCount", 30)), time.Duration(GetConfigDefaultByInt("BreakerRecoveryTimes", 10))*time.Second)
	}

	return client.NewBidirectionalXClient(*serverNode, client.Failover, client.RandomSelect, d, option,msgChan), nil
}
