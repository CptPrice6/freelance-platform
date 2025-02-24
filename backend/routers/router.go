package routers

import (
	"backend/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/random", &controllers.NumberController{})
}
