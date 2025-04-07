package routers

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/beego/beego/v2/server/web"
)

func init() {

	// auth logic
	web.Router("/register", &controllers.AuthController{}, "post:RegisterHandler")
	web.Router("/login", &controllers.AuthController{}, "post:LoginHandler")
	web.Router("/refresh", &controllers.AuthController{}, "post:RefreshTokenHandler")

	// user logic
	web.InsertFilter("/user/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/user", &controllers.UserController{}, "get:GetUserHandler")
	web.Router("/user", &controllers.UserController{}, "put:UpdateUserHandler")
	web.Router("/user", &controllers.UserController{}, "delete:DeleteUserHandler")
	web.Router("/user/auth", &controllers.AuthController{}, "get:AuthHandler")

	web.Router("/user/attachments/:id", &controllers.AttachmentController{}, "get:DownloadAttachment")

	// freelancer role-specific logic
	web.Router("/user/freelancer", &controllers.FreelancerController{}, "put:UpdateFreelancerDataHandler")

	web.Router("/user/freelancer/skills", &controllers.SkillController{}, "post:AddFreelancerSkillHandler")
	web.Router("/user/freelancer/skills", &controllers.SkillController{}, "delete:DeleteFreelancerSkillHandler")

	web.Router("/user/freelancer/jobs", &controllers.JobController{}, "get:GetFreelancerJobsHandler")
	web.Router("/user/freelancer/jobs/:id", &controllers.JobController{}, "get:GetFreelancerJobHandler")

	web.Router("/user/freelancer/applications", &controllers.ApplicationController{}, "post:SubmitApplication")
	web.Router("/user/freelancer/applications", &controllers.ApplicationController{}, "get:GetFreelancerApplications")
	web.Router("/user/freelancer/applications/:id", &controllers.ApplicationController{}, "get:GetFreelancerApplication")
	web.Router("/user/freelancer/applications/:id", &controllers.ApplicationController{}, "delete:DeleteApplication")
	web.Router("/user/freelancer/applications/:id", &controllers.ApplicationController{}, "put:UpdateApplication")

	// skill logic
	web.InsertFilter("/skills/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/skills", &controllers.SkillController{}, "get:GetSkillsHandler")

	// client role-specific logic
	web.Router("/user/client", &controllers.ClientController{}, "put:UpdateClientDataHandler")

	web.Router("/user/client/jobs", &controllers.JobController{}, "post:CreateJobHandler")
	web.Router("/user/client/jobs", &controllers.JobController{}, "get:GetClientJobsHandler")
	web.Router("/user/client/jobs/:id", &controllers.JobController{}, "get:GetClientJobHandler")
	web.Router("/user/client/jobs/:id", &controllers.JobController{}, "delete:DeleteClientJobHandler")
	web.Router("/user/client/jobs/:id", &controllers.JobController{}, "put:UpdateClientJobHandler")
	web.Router("/user/client/jobs/:id/complete", &controllers.JobController{}, "post:CompleteJobHandler")

	web.Router("/user/client/jobs/applications/:id", &controllers.ApplicationController{}, "post:ChangeApplicationStatus")

	// admin role-specific logic
	web.InsertFilter("/admin/*", web.BeforeRouter, middleware.AdminAuthMiddleware)
	web.Router("/admin/users", &controllers.AdminController{}, "get:GetUsersHandler")
	web.Router("/admin/users/:id", &controllers.AdminController{}, "delete:DeleteUserHandler")
	web.Router("/admin/users/:id", &controllers.AdminController{}, "put:UpdateUserHandler")

	web.Router("/admin/skills", &controllers.SkillController{}, "post:AddSkillHandler")
	web.Router("/admin/skills/:id", &controllers.SkillController{}, "delete:DeleteSkillHandler")
	web.Router("/admin/skills/:id", &controllers.SkillController{}, "put:UpdateSkillHandler")

	//TODO:
	web.Router("/admin/jobs/:id", &controllers.JobController{}, "delete:DeleteJobHandler")

	// public freelancer logic
	web.InsertFilter("/freelancers/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/freelancers", &controllers.FreelancerController{}, "get:GetFreelancersHandler")
	web.Router("/freelancers/:id", &controllers.FreelancerController{}, "get:GetFreelancerHandler")

	// public client logic
	web.InsertFilter("/clients/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/clients", &controllers.ClientController{}, "get:GetClientsHandler")
	web.Router("/clients/:id", &controllers.ClientController{}, "get:GetClientHandler")

	//public job logic
	web.InsertFilter("/jobs/*", web.BeforeRouter, middleware.UserAuthMiddleware)
	web.Router("/jobs", &controllers.JobController{}, "get:GetJobsHandler")
	web.Router("/jobs/:id", &controllers.JobController{}, "get:GetJobHandler")

}
