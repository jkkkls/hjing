package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jkkkls/hjing/rpc"
	"github.com/jkkkls/hjing/utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type MiddleFunc func(*routing.Context) error

var (
	ApiMiddleFunc MiddleFunc
)

func RegisterApiMiddleFunc(f MiddleFunc) {
	ApiMiddleFunc = f
}

func getFastRemote(ctx *routing.Context) string {
	ip := string(ctx.Request.Header.Peek("x-forwarded-for"))
	if ip != "" && ip != "unknown" {
		ip = strings.Split(ip, ",")[0]
	}
	if ip == "" {
		ip = strings.Split(ctx.RemoteIP().String(), ":")[0]
	}
	return ip
}

type HttpParam func(*routing.Router)

func RunApiHttp(port int, params ...HttpParam) error {
	router := routing.New()

	for _, v := range params {
		v(router)
	}

	router.Options("/*", func(ctx *routing.Context) error {
		return nil
	})

	//api
	api := router.Group("/api")
	api.Use(func(ctx *routing.Context) error {
		if ApiMiddleFunc != nil {
			return ApiMiddleFunc(ctx)
		}
		return nil
	})
	api.Post("/*", ProtectedHandler)

	utils.Info("启动api2服务器", "port", port)

	return fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), router.HandleRequest)
}

// ProtectedHandler2 入口
func ProtectedHandler(ctx *routing.Context) error {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")

	urlStr := string(ctx.Path())
	remote := getFastRemote(ctx)
	arr := strings.Split(urlStr, "/")
	if len(arr) != 4 || arr[1] != "api" {
		JsonResponse2(ctx, map[string]interface{}{
			"code":    1,
			"codeMsg": "请求地址格式错误",
		})
		return nil
	}

	reqBuff := ctx.PostBody()
	if len(reqBuff) == 0 {
		reqBuff = []byte("{}")
	}

	context := &rpc.Context{Remote: remote}

	// begin := time.Now()
	serviceMethon := utils.Upper(arr[2], 1) + "." + utils.Upper(arr[3], 1)
	var (
		rspBuff []byte
		err     error
	)
	utils.ProtectCall(func() {
		_, rspBuff, err = rpc.JsonCall(context, serviceMethon, reqBuff)
	}, func() {
		err = fmt.Errorf("server internal error")
	})
	// cost := time.Since(begin).String()
	// utils.Debug("请求信息", "api", serviceMethon, "remote", remote, "url", urlStr, "reqBuff", string(reqBuff), "rspBuff", string(rspBuff), "err", err, "cost", cost)
	if err != nil {
		JsonResponse2(ctx, map[string]interface{}{
			"code":    1,
			"codeMsg": err.Error(),
		})
		utils.Info("JsonCall失败", "api", serviceMethon, "remote", remote, "url", urlStr, "err", err)
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
