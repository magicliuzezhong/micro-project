//
// Package jeager
// @Description：jeager配置
// @Author：liuzezhong 2021/6/29 3:54 下午
// @Company cloud-ark.com
//
package jeager

import (
	"github.com/kataras/iris/v12/context"
	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"micro-project/internal/pkg/util"
)

var (
	tracer         opentracing.Tracer                         //tracer
	serviceName    = util.GetApplication().Server.ServerName  //服务名称
	jeagerHostPort = util.GetApplication().JeagerConf.Url     //jeager服务端地址
	jeagerEnabled  = util.GetApplication().JeagerConf.Enabled //是否启用jeager
	jeagerType     = util.GetApplication().JeagerConf.Type    // 采样类型 const：固定采样
	jeagerParam    = util.GetApplication().JeagerConf.Param   // 1:全采样，0:不采样
)

//
// init
// @Description: 初始化方法
//
func init() {
	if jeagerEnabled {
		initTrace(serviceName, jeagerHostPort)
	}
}

//
// initTrace
// @Description: 初始化jeager
// @param serviceName 服务名称，一般取成服务名称 + tag
// @param jaegerHostPort jeager服务端的地址
//
func initTrace(serviceName string, jaegerHostPort string) {
	var localJeagerType = "const" //固定采样
	if jeagerType == "const" {
		localJeagerType = jeagerType
	}
	var localJeagerParam float64 = 1 //全采样
	if jeagerParam == 0 || jeagerParam == 1 {
		localJeagerParam = float64(jeagerParam)
	}
	var cfg = &jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  localJeagerType,  // const：固定采样
			Param: localJeagerParam, // 1:全采样，0:不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           false, // 是否在终端输出Reporting span信息
			LocalAgentHostPort: jaegerHostPort,
		},
	}
	var innerTracer, _, err = cfg.NewTracer()
	if err != nil {
		panic("注册jeager异常" + err.Error())
	}
	tracer = innerTracer
	opentracing.SetGlobalTracer(tracer)
}

//
// JeagerTraceFilter
// @Description: jeager全链路调用过滤器
// @param context iris上下文
//
func JeagerTraceFilter(context context.Context) {
	if !jeagerEnabled {
		context.Next()
		return
	}
	var url = context.Request().URL
	if url.String() == "/health" || url.String() == "/metrics" { //如果是健康检查或者普罗米修斯那么不进行全链路追踪
		context.Next()
		return
	}
	var tracerId = context.GetHeader("Uber-Trace-Id")
	var span opentracing.Span
	if tracerId == "" {
		span = tracer.StartSpan(serviceName)
		_ = tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(context.Request().Header))
	} else {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(context.Request().Header))
		span = tracer.StartSpan(serviceName, opentracing.ChildOf(spanCtx))
	}
	span.SetTag("url", url)
	span.SetTag("method", context.Request().Method)
	span.Finish()
	context.Next()
}
