package gin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"itmo-profile/config"
	"itmo-profile/internal/gin/controllers"
)

func Create() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	var secret = config.GetEnv("API_SECRET", "supersecret")
	store := cookie.NewStore([]byte(secret))

	r.Use(sessions.Sessions("EASYSESSIONID1", store))

	controllers.SetupControllers(r)

	return r
}
