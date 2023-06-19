package routers

import (
	"encoding/json"
	"fmt"
	"github.com/CloudNativeGame/aigc-gateway/pkg/resources"
	mem "github.com/CloudNativeGame/aigc-gateway/pkg/session"
	"github.com/CloudNativeGame/aigc-gateway/pkg/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/logto-io/go/client"
	"github.com/logto-io/go/core"
)

func RegisterResourceRouters(router *gin.Engine, logtoConfig *client.LogtoConfig) {
	// Add a link to perform a sign-in request on the home page
	router.GET("/resources", func(ctx *gin.Context) {
		rm := resources.NewResourceManager()
		resources, err := rm.ListResources(nil, nil)
		if err != nil {
			ctx.String(400, err.Error())
		}
		ctx.JSON(200, resources)
	})

	router.GET("/resource/:namespace/:name/:id", func(ctx *gin.Context) {
		name := ctx.Param("name")
		namespace := ctx.Param("namespace")
		id := ctx.Param("id")

		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		userInfo, err := logtoClient.FetchUserInfo()
		if err != nil {
			ctx.Error(err)
		}
		cm := userInfo.CustomData
		key := fmt.Sprintf("%s-%s", namespace, name)

		value := cm[key]
		if value == nil {
			ctx.String(403, "have no privilege to get the instance with ID %d", id)
			return
		}

		rm := &resources.ResourceMeta{}
		valueBytes, err := interfaceToBytes(value)
		if err != nil {
			ctx.Error(err)
		}
		err = json.Unmarshal(valueBytes, rm)
		if err != nil {
			ctx.Error(err)
		}

		if id != rm.ID || name != rm.Name || namespace != rm.Namespace {
			ctx.String(403, "have no privilege to get the instance with ID %d", id)
			return
		}

		// add json wrapper
		resourceManager := resources.NewResourceManager()
		resource, err := resourceManager.GetResource(rm)

		if err != nil {
			if resources.GetErrorReason(err) == resources.PauseReason {
				ctx.Status(423)
				return
			}
			ctx.String(400, err.Error())
		}

		ctx.JSON(200, resource)
	})

	// create a new one
	router.PUT("/resource/:namespace/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		namespace := ctx.Param("namespace")

		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		userInfo, err := logtoClient.FetchUserInfo()

		if err != nil {
			ctx.Error(err)
		}

		cm := userInfo.CustomData

		if cm == nil {
			cm = make(map[string]interface{})
		}

		key := fmt.Sprintf("%s-%s", namespace, name)

		// check if already installed
		if cm[key] != nil {
			ctx.String(400, "already installed")
			return
		}

		rm := &resources.ResourceMeta{
			Name:      name,
			Namespace: namespace,
		}

		// add json wrapper
		resourceManager := resources.NewResourceManager()
		meta, err := resourceManager.CreateResource(rm)
		if err != nil {
			ctx.Error(err)
		}

		cm[key] = meta

		err = user.UpdateUserMetaData(userInfo.Sub, cm)
		if err != nil {
			ctx.Error(err)
		}

		ctx.JSON(200, meta)

	})

	router.POST("/resource/:namespace/:name/pause", func(ctx *gin.Context) {
		name := ctx.Param("name")
		namespace := ctx.Param("namespace")

		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		userInfo, err := logtoClient.FetchUserInfo()

		if err != nil {
			ctx.Error(err)
		}

		cm := userInfo.CustomData

		key := fmt.Sprintf("%s-%s", namespace, name)
		value := cm[key]

		valueBytes, err := interfaceToBytes(value)

		if err != nil {
			ctx.Error(err)
		}

		rm := &resources.ResourceMeta{}

		err = json.Unmarshal(valueBytes, rm)

		if err != nil {
			ctx.Error(err)
		}

		// add json wrapper
		resourceManager := resources.NewResourceManager()
		err = resourceManager.PauseResource(rm)

		if err != nil {
			ctx.String(400, err.Error())
		}
		ctx.Status(200)
	})

	router.POST("/resource/:namespace/:name/recover", func(ctx *gin.Context) {
		name := ctx.Param("name")
		namespace := ctx.Param("namespace")

		session := sessions.Default(ctx)
		logtoClient := client.NewLogtoClient(
			logtoConfig,
			&mem.SessionStorage{Session: session},
		)

		userInfo, err := logtoClient.FetchUserInfo()

		if err != nil {
			ctx.Error(err)
		}

		cm := userInfo.CustomData

		key := fmt.Sprintf("%s-%s", namespace, name)
		value := cm[key]

		if value == nil {
			ctx.Status(200)
			return
		}

		valueBytes, err := interfaceToBytes(value)

		if err != nil {
			ctx.Error(err)
		}

		rm := &resources.ResourceMeta{}

		json.Unmarshal(valueBytes, rm)

		// add json wrapper
		resourceManager := resources.NewResourceManager()
		_, err = resourceManager.RecoverResource(rm)

		if err != nil {
			ctx.String(400, err.Error())
		}
		ctx.Status(200)
		return
	})
}

func interfaceToBytes(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func getUserInfoFromSession(logtoConfig *client.LogtoConfig, ctx *gin.Context) (core.IdTokenClaims, error) {
	// Init LogtoClient
	session := sessions.Default(ctx)
	logtoClient := client.NewLogtoClient(
		logtoConfig,
		&mem.SessionStorage{Session: session},
	)

	tc, err := logtoClient.GetIdTokenClaims()

	if err != nil {
		return core.IdTokenClaims{}, err
	}

	return tc, nil
}
