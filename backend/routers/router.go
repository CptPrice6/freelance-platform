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

	web.InsertFilter("/user/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/user/random", &controllers.NumberController{})
	web.Router("/user/data", &controllers.UserController{})
	web.Router("/user/logout", &controllers.AuthController{}, "post:LogoutHandler")

	web.InsertFilter("/admin/*", web.BeforeRouter, middleware.AdminAuthMiddleware)
	web.Router("/admin/logout-user", &controllers.AdminController{}, "post:LogoutUserHandler")

}
