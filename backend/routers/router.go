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

	// user logic
	web.InsertFilter("/user/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/user", &controllers.UserController{}, "get:GetUserHandler")
	web.Router("/user", &controllers.UserController{}, "put:UpdateUserHandler")
	web.Router("/user", &controllers.UserController{}, "delete:DeleteUserHandler")
	web.Router("/user/logout", &controllers.AuthController{}, "post:LogoutHandler")
	web.Router("/user/auth", &controllers.AuthController{}, "get:AuthHandler")

	web.Router("/user/freelancer", &controllers.FreelancerController{}, "put:UpdateFreelancerDataHandler")
	web.Router("/user/freelancer/skills", &controllers.SkillController{}, "post:AddFreelancerSkillHandler")
	web.Router("/user/freelancer/skills", &controllers.SkillController{}, "delete:DeleteFreelancerSkillHandler")

	web.InsertFilter("/skills/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/skills", &controllers.SkillController{}, "get:GetSkillsHandler")

	web.Router("/user/client", &controllers.ClientController{}, "put:UpdateClientDataHandler")

	web.InsertFilter("/admin/*", web.BeforeRouter, middleware.AdminAuthMiddleware)
	web.Router("/admin/users", &controllers.AdminController{}, "get:GetUsersHandler")
	web.Router("/admin/users/:id", &controllers.AdminController{}, "delete:DeleteUserHandler")
	web.Router("/admin/users/:id", &controllers.AdminController{}, "put:UpdateUserHandler")

	web.Router("/admin/skills", &controllers.SkillController{}, "post:AddSkillHandler")
	web.Router("/admin/skills/:id", &controllers.SkillController{}, "delete:DeleteSkillHandler")
	web.Router("/admin/skills/:id", &controllers.SkillController{}, "put:UpdateSkillHandler")

	web.InsertFilter("/freelancers/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/freelancers", &controllers.FreelancerController{}, "get:GetFreelancersHandler")
	web.Router("/freelancers/:id", &controllers.FreelancerController{}, "get:GetFreelancerHandler")

}
