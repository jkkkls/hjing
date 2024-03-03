package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jkkkls/hjing/utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func getRemote2(ctx *routing.Context) string {
	ip := string(ctx.Request.Header.Peek("x-forwarded-for"))
	if ip != "" && ip != "unknown" {
		ip = strings.Split(ip, ",")[0]
	}
	if ip == "" {
		ip = strings.Split(ctx.RemoteIP().String(), ":")[0]
	}
	return ip
}

func RunHttpGateway(port int) error {
	router := routing.New()

	router.Options("/*", func(ctx *routing.Context) error {
		return nil
	})
	router.Post("/rpcapi/*", rpcHandler)
	utils.Info("启动rpc http服务器", "port", port)

	fasthttp.AcquireResponse().Header.SetServer("hjing server")
	return fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), router.HandleRequest)
}

// rpcHandler 入口
// url格式：/rpcapi/{ServiceName}/{MethodName}
// url格式：/rpcapi/{NodeName}/{ServiceName}/{MethodName}
func rpcHandler(ctx *routing.Context) error {
	urlStr := string(ctx.Path())
	remote := getRemote2(ctx)
	arr := strings.Split(urlStr, "/")
	if (len(arr) != 4 && len(arr) != 5) || arr[1] != "rpcapi" {
		JsonResponse2(ctx, map[string]interface{}{
			"code":    1,
			"codeMsg": "请求地址格式错误",
		})
		// utils.Info("登陆请求地址格式错误", "url", urlStr)
		return nil
	}
	var nodeName, serviceName, methonName string
	if len(arr) == 4 {
		serviceName = arr[2]
		methonName = arr[3]
	} else {
		nodeName = arr[2]
		serviceName = arr[3]
		methonName = arr[4]
	}

	reqBuff := ctx.PostBody()
	if len(reqBuff) == 0 {
		reqBuff = []byte("{}")
	}

	context := &Context{Remote: remote}

	// begin := time.Now()
	serviceMethon := utils.Upper(serviceName, 1) + "." + utils.Upper(methonName, 1)
	var (
		rspBuff []byte
		err     error
	)
	utils.ProtectCall(func() {
		if nodeName != "" {
			_, rspBuff, err = NodeJsonCallWithConn(context, nodeName, serviceMethon, reqBuff)
		} else {
			_, rspBuff, err = JsonCall(context, serviceMethon, reqBuff)
		}
	}, func() {
		err = fmt.Errorf("server internal error")
	})
	if err != nil {
		JsonResponse2(ctx, map[string]interface{}{
			"code":    1,
			"codeMsg": err.Error(),
		})
		utils.Info("JsonCall失败", "err", err)
		return nil
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(rspBuff)
	return nil
}

func JsonResponse2(ctx *routing.Context, response interface{}) {
	json, err := json.Marshal(response)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(json)
}
