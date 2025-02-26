package routers

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/register", &controllers.AuthController{}, "post:RegisterHandler")
	web.Router("/login", &controllers.AuthController{}, "post:LoginHandler")
	web.Router("/refresh-token", &controllers.AuthController{}, "post:RefreshTokenHandler")

	web.InsertFilter("/random", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/random", &controllers.NumberController{})

	web.InsertFilter("/user", web.BeforeRouter, middleware.AdminAuthMiddleware)
	web.Router("/user", &controllers.UserController{})

}
