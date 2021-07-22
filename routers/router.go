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
	// http://localhost/api/v1.0/user
	beego.Router("/api/v1.0/user", &controllers.UserController{}, "get:GetUser")
	// http://localhost/api/v1.0/user/name
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:UpName")
	// http://localhost/api/v1.0/user/auth
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{}, "get:GetUser;post:PostAuth")
	// http://localhost/api/v1.0/user/houses
	beego.Router("/api/v1.0/user/houses", &controllers.HouseController{}, "get:GetHouse")
	// http://localhost/api/v1.0/houses
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "post:PostHouse")
	// http://localhost/api/v1.0/houses/7/images
	beego.Router("/api/v1.0/houses/?:houseid/images", &controllers.HouseController{}, "post:PostHouseImages")
	// http://localhost/api/v1.0/houses/7
	beego.Router("/api/v1.0/houses/?:houseid", &controllers.HouseController{}, "get:GetDetailHouseData")
}
