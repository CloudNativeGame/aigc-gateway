package routers

import (
	"encoding/json"
	"github.com/CloudNativeGame/aigc-gateway/pkg/resources"
	mem "github.com/CloudNativeGame/aigc-gateway/pkg/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var redirectURL = os.Getenv("Redirect_Url")

func RegisterSignRouters(router *gin.Engine, logtoConfig *client.LogtoConfig) {
	// Add a link to perform a sign-in request on the home page
	router.GET("/", func(ctx *gin.Context) {
		// Init LogtoClient
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		resp, err := logtoClient.FetchUserInfo()

		authState := false

		if err == nil {
			authState = true
		}

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title":     "Main website",
			"authState": authState,
			"userInfo":  resp,
		})
	})

	router.GET("/auth", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		userInfo, err := logtoClient.FetchUserInfo()
		if err != nil {
			ctx.String(http.StatusUnauthorized, err.Error())
			return
		}

		originUrl := ctx.GetHeader("X-Original-Url")

		resourceManager := resources.NewResourceManager()

		for _, info := range userInfo.CustomData {
			if info == nil {
				continue
			}
			valueBytes, err := interfaceToBytes(info)
			if err != nil {
				ctx.Error(err)
			}

			rm := &resources.ResourceMeta{}

			err = json.Unmarshal(valueBytes, rm)

			if err != nil {
				ctx.Error(err)
			}

			endpoint, err := resourceManager.GetResourceEndpoint(rm)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
			}

			if endpoint != "" && strings.Contains(originUrl, endpoint) {
				ctx.String(http.StatusOK, "")
				return
			}
		}

		ctx.String(http.StatusUnauthorized, "User and Access Endpoint Mismatch")
		return
	})

	// Add a route for handling sign-in requests
	router.GET("/sign-in", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		signInUri, _ := url.JoinPath(redirectURL, "sign-in-callback")

		// The sign-in request is handled by Logto.
		// The user will be redirected to the Redirect URI on signed in.
		signInUri, err := logtoClient.SignIn(signInUri)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect the user to the Logto sign-in page.
		ctx.Redirect(http.StatusTemporaryRedirect, signInUri)
	})

	// Add a route for handling signing out requests
	router.GET("/sign-out", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		// The sign-out request is handled by Logto.
		// The user will be redirected to the Post Sign-out Redirect URI on signed out.
		signOutUri, signOutErr := logtoClient.SignOut(redirectURL)

		if signOutErr != nil {
			ctx.String(http.StatusOK, signOutErr.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, signOutUri)
	})

	// Add a route for handling sign-in callback requests
	router.GET("/sign-in-callback", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		// The sign-in callback request is handled by Logto
		err := logtoClient.HandleSignInCallback(ctx.Request)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Jump to the page specified by the developer.
		// This example takes the user back to the home page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	})
}
