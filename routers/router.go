package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"ihome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	// http://localhost/api/v1.0/areas
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetArea")
	// http://localhost/api/v1.0/houses/index
	beego.Router("/api/v1.0/houses/index", &controllers.HousesIndexController{}, "get:HousesIndex")
	// http://localhost/api/v1.0/session
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:GetSessionData")
	// http://localhost/api/v1.0/users
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")
	// http://localhost/api/v1.0/session    delete
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "delete:DelSessionData")
	// http://localhost/api/v1.0/sessions  post   登录
	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	// http://localhost/api/v1.0/user/avatar
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:Postavatar")
}
