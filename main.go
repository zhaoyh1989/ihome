package main

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	_ "ihome/models"
	_ "ihome/routers"
	"net/http"
	"strings"
)

func main() {
	ignoreStaticPath()
	beego.Run()
}

func ignoreStaticPath() {
	beego.SetStaticPath("/group1/M00", "fdfs/storage_data/data")
	//透明static
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}
func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	logs.Debug("request url: ", orpath)
	//如果请求url还有api字段，说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
	//if orpath == "/" || strings.Index(orpath, "profile.html") >= 0 || strings.Index(orpath, "static") >= 0 {
	//	logs.Debug("start run http.ServeFile")
	//	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
	//} else {
	//	logs.Debug("start run ctx.Redirect")
	//	ctx.Redirect(302, "/")
	//}
}
