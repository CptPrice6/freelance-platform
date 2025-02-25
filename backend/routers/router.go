package routers

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/register", &controllers.AuthController{}, "post:RegisterHandler")
	web.Router("/login", &controllers.AuthController{}, "post:LoginHandler")

	web.InsertFilter("/random", web.BeforeRouter, middleware.JWTMiddleware)

	web.Router("/random", &controllers.NumberController{})

	web.InsertFilter("/user", web.BeforeRouter, middleware.AdminMiddleware)

	web.Router("/user", &controllers.UserController{})

}
