package mainad

import (
	"eleliafrika.com/backend/users"
	"github.com/gin-gonic/gin"
)

func Mainadsroutes(router *gin.Engine) {
	mainadsroutes := router.Group("/mainads")
	{
		mainadsroutes.POST("/create", users.JWTAuthMiddleWare(), CreateMainAd)
		mainadsroutes.GET("/getmainads", GetAllMainAds)
		mainadsroutes.GET("/getsinglemainad", GetSingleMainAd)
		mainadsroutes.POST("/update", users.JWTAuthMiddleWare(), UpdateMainAd)
		mainadsroutes.POST("/delete", users.JWTAuthMiddleWare(), DeleteMainAd)
		mainadsroutes.POST("/restore", users.JWTAuthMiddleWare(), RestoreMainAd)
		mainadsroutes.POST("/activate", users.JWTAuthMiddleWare(), ActivateMainAd)
		mainadsroutes.POST("/deactivate", users.JWTAuthMiddleWare(), DeactivateMainAd)
	}
}
