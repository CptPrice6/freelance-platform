package routers

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/register", &controllers.AuthController{}, "post:RegisterHandler")
	web.Router("/login", &controllers.AuthController{}, "post:LoginHandler")
	web.Router("/refresh", &controllers.AuthController{}, "post:RefreshTokenHandler")

	web.InsertFilter("/user/*", web.BeforeRouter, middleware.UserAuthMiddleware)

	web.Router("/user", &controllers.UserController{})
	web.Router("/user/random", &controllers.NumberController{})
	web.Router("/user/logout", &controllers.AuthController{}, "post:LogoutHandler")
	web.Router("/user/auth", &controllers.AuthController{}, "get:AuthHandler")

	web.InsertFilter("/admin/*", web.BeforeRouter, middleware.AdminAuthMiddleware)
	web.Router("/admin/logout-user", &controllers.AdminController{}, "post:LogoutUserHandler")

}
