package main

import (
	"os"

	"github.com/CloudNativeGame/aigc-gateway/pkg/middleware"
	"github.com/CloudNativeGame/aigc-gateway/pkg/routers"
	"github.com/CloudNativeGame/aigc-gateway/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
)

func main() {
	router := gin.Default()

	// load templates
	router.Delims("{[{", "}]}")
	router.LoadHTMLGlob("aigc-dashboard/dist/*.html")
	router.Use(static.Serve("/assets", static.LocalFile("aigc-dashboard/dist/assets", true)))

	endpoint := os.Getenv("Endpoint")
	logtoConfig := &client.LogtoConfig{

		Endpoint:  endpoint,
		AppId:     os.Getenv("App_Id"),
		AppSecret: os.Getenv("App_Secret"),
		Scopes:    []string{"email", "custom_data"},
	}
	// We use memory-based session in this example
	store := memstore.NewStore([]byte("your session secret"))
	store.Options(sessions.Options{
		Domain: utils.GetDomainFromEndpoint(endpoint),
		Path:   "/",
		MaxAge: 604800,
	})
	router.Use(sessions.Sessions("logto-session", store))
	// middleware
	router.Use(middleware.Cors())
	routers.RegisterSignRouters(router, logtoConfig)
	routers.RegisterResourceRouters(router, logtoConfig)

	router.Run(":8090")
}
