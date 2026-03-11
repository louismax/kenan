package constant

const (
	Exception         int = 9900 + iota //系统异常
	IllegalRequest                      //非法请求
	LoginExpired                        //登录过期
	InterfaceNotExist                   //接口不存在
)

const (
	ReqDataError          int = 10000 + iota //请求参数异常
	ReqArgsMarshalFail                       //请求参数编码失败
	ReqArgsUnmarshalFail                     //请求参数解析失败
	RespArgsMarshalFail                      //返回参数编码失败
	RespArgsUnmarshalFail                    //返回参数解析失败
	RespBusinessFail                         //通用业务异常
	InvalidRequest                           //非法请求
	VerificationCodeError                    //验证码错误

	DBOpenError       //数据库连接失败
	DBGetError        //数据库数据获取失败
	DBOperatingError  //数据库数据操作失败
	DBNoDataError     //数据库无数据
	DBDataRepeatError //数据库数据已存在

	RCOpenError      //Redis连接失败
	RCGetError       //缓存获取失败
	RCOperatingError //缓存操作失败

	CallServiceFail //调用服务异常

)
