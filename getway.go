package kenan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type Router struct {
	Ctx iris.Context
}

func (app *Router) CreateNewLink(ctx iris.Context) error {
	app.Ctx = ctx
	return nil
}

func (app *Router) ReturnToClientMap(m1 *map[string]interface{}, statusCode ...int) {
	//ctx.Header("Access-Control-Allow-Origin", "*")
	app.Ctx.Header("Access-Control-Allow-Origin", app.Ctx.GetHeader("Origin"))
	app.Ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	app.Ctx.Header("Access-Control-Allow-Headers", "X-Requested-With,content-type")
	app.Ctx.Header("Access-Control-Allow-Credentials", "true")
	if len(statusCode) > 0 {
		app.Ctx.StatusCode(statusCode[0])
	} else {
		app.Ctx.StatusCode(iris.StatusOK)
	}

	if err := app.Ctx.JSON(&m1); err != nil {
		app.Ctx.StatusCode(iris.StatusRequestTimeout)
	}
}

func (app *Router) ReturnSSEMessage(m []byte, flusher http.Flusher) {
	_, _ = app.Ctx.Writef("data: %s\n\n", m)
	flusher.Flush()
}

func (app *Router) ReturnSSEResult(m *map[string]interface{}, flusher http.Flusher) {
	_, _ = app.Ctx.Writef("event: result\n")
	_, _ = app.Ctx.Writef("data: %s\n\n", app.UseComplicated(m))
	flusher.Flush()
}

func ReturnErrorMap(ack int, msg string) *map[string]interface{} {
	var rtn map[string]interface{}
	rtn = make(map[string]interface{})
	rtn["Success"] = false
	rtn["Result"] = ack
	rtn["ErrorMessage"] = msg
	return &rtn
}

func ReturnSuccessMap(d interface{}) *map[string]interface{} {
	var rtn map[string]interface{}
	rtn = make(map[string]interface{})
	rtn["Success"] = true
	if d != nil {
		rtn["Result"] = d
	}
	return &rtn
}

func (app *Router) CallServiceReply(serverPath, fun string, session map[string]interface{}) (*Reply, error) {
	var err error
	var nextArgs Args

	nextArgs.HttpHeader = app.Ctx.Request().Header
	nextArgs.Session = session
	nextArgs.Data, _ = io.ReadAll(app.Ctx.Request().Body)

	var mz client.XClient

	if v, ok := RpcMapSync.Load(serverPath); ok {
		mz = v.(client.XClient)
	} else {
		mz, err = RpcNewClient(serverPath)
		if err != nil {
			LogError("客户端注册失败,err:%+v", err)
			return nil, err
		}
		RpcMapSync.Store(serverPath, mz)
	}

	var replyX Reply
	err = mz.Call(context.Background(), fun, &nextArgs, &replyX)
	if err != nil {
		LogError("Call服务异常,err:%+v,尝试重新建立连接！\n", err)
		if e := mz.Close(); e != nil { //先尝试关闭之前的连接!
			LogError("关闭旧的客户端连接异常,err:%+v\n", err)
		}
		//重新建立一个连接
		newMZ, errX := RpcNewClient(serverPath)
		if errX != nil {
			LogError("客户端重新注册失败,err:%+v", errX)
			return nil, errX
		}
		RpcMapSync.Store(serverPath, newMZ)
		if err = newMZ.Call(context.Background(), fun, &nextArgs, &replyX); err != nil {
			return nil, err
		}
	}

	return &replyX, nil
}
func (app *Router) CallBidirectionalService(serverPath, fun string, session map[string]interface{}, data interface{}, msgChan chan<- *protocol.Message) (client.XClient, *map[string]interface{}, error) {
	var err error

	var nextArgs Args

	nextArgs.HttpHeader = app.Ctx.Request().Header
	nextArgs.Session = session
	nextArgs.Data = data

	xCli, err := RpcNewBidirectionalClient(serverPath, msgChan)
	if err != nil {
		LogError("客户端注册失败,err:%+v", err)
		return nil, nil, err
	}
	var replyX Reply

	err = xCli.Call(context.Background(), fun, &nextArgs, &replyX)
	if err != nil {
		LogError("Call服务异常,err: %+v", err)
		xCli.Close()
		return nil, nil, err
	}
	Result := make(map[string]interface{})
	err = json.Unmarshal(replyX.Data.([]byte), &Result)
	if err != nil {
		xCli.Close()
		return nil, nil, errors.New("RPC请求结果解析失败")
	}
	return xCli, &Result, nil
}

func (app *Router) CallService(serverPath, fun string, session map[string]interface{}) (*map[string]interface{}, error) {
	var err error

	var nextArgs Args

	nextArgs.HttpHeader = app.Ctx.Request().Header
	nextArgs.Session = session

	nextArgs.Data, _ = io.ReadAll(app.Ctx.Request().Body)
	if len(nextArgs.Data.([]byte)) == 0 {
		nextArgs.Data = nil
	}

	var mz client.XClient

	if v, ok := RpcMapSync.Load(serverPath); ok {
		mz = v.(client.XClient)
	} else {
		mz, err = RpcNewClient(serverPath)
		if err != nil {
			LogError("客户端注册失败,err:%+v", err)
			return nil, err
		}
		RpcMapSync.Store(serverPath, mz)
	}
	var replyX Reply
	err = mz.Call(context.Background(), fun, &nextArgs, &replyX)
	if err != nil {
		LogError("Call服务异常,svr:%s,fun:%s,err:%+v,尝试重新建立连接！\n", serverPath, fun, err)
		if e := mz.Close(); e != nil { //先尝试关闭之前的连接!
			LogError("关闭旧的客户端连接异常,err:%+v\n", err)
		}
		//重新建立一个连接
		newMZ, errX := RpcNewClient(serverPath)
		if errX != nil {
			LogError("客户端重新注册失败,err:%+v", errX)
			return nil, errX
		}
		RpcMapSync.Store(serverPath, newMZ)
		err = newMZ.Call(context.Background(), fun, &nextArgs, &replyX)
		if err != nil {
			return nil, err
		}
	}

	Result := make(map[string]interface{})
	err = json.Unmarshal(replyX.Data.([]byte), &Result)
	if err != nil {
		return nil, errors.New("RPC请求结果解析失败")
	}
	return &Result, nil
}

func (app *Router) CallServiceData(serverPath, fun string, session map[string]interface{}, data interface{}) (*map[string]interface{}, error) {
	var err error

	var nextArgs Args

	nextArgs.HttpHeader = app.Ctx.Request().Header
	nextArgs.Session = session
	nextArgs.Data = data

	var mz client.XClient

	if v, ok := RpcMapSync.Load(serverPath); ok {
		mz = v.(client.XClient)
	} else {
		mz, err = RpcNewClient(serverPath)
		if err != nil {
			LogError("客户端注册失败,err:%+v", err)
			return nil, err
		}
		RpcMapSync.Store(serverPath, mz)
	}
	var replyX Reply
	err = mz.Call(context.Background(), fun, &nextArgs, &replyX)
	if err != nil {
		LogError("Call服务异常,svr:%s,fun:%s,err:%+v,尝试重新建立连接！\n", serverPath, fun, err)
		if e := mz.Close(); e != nil { //先尝试关闭之前的连接!
			LogError("关闭旧的客户端连接异常,err:%+v\n", err)
		}
		//重新建立一个连接
		newMZ, errX := RpcNewClient(serverPath)
		if errX != nil {
			LogError("客户端重新注册失败,err:%+v", errX)
			return nil, errX
		}
		RpcMapSync.Store(serverPath, newMZ)
		err = newMZ.Call(context.Background(), fun, &nextArgs, &replyX)
		if err != nil {
			return nil, err
		}
	}

	Result := make(map[string]interface{})
	err = json.Unmarshal(replyX.Data.([]byte), &Result)
	if err != nil {
		return nil, errors.New("RPC请求结果解析失败")
	}
	return &Result, nil
}

// UseComplicated 复杂对象转为[]byte
func (app *Router) UseComplicated(body interface{}) []byte {
	jsons, errs := json.Marshal(body) //转换成JSON返回的是byte[]
	if errs != nil {
		LogError("json转移失败,err:%+v", errs)
	}
	return jsons
}
