package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/local-server-test1/cmd/api"
	"github.com/local-server-test1/cmd/database"
	// arm "github.com/synspective/syns-arms/pkg/sdk"
	// "github.com/synspective/syns-arms/pkg/sdk/customer"
)

// canAccessServiceUUID checks whether service/UUID is accessible for user or not through arm connect
func canAccessServiceUUID(service string, struuids []string, authtoken string) (canAccess bool) {
	// res, err := arm.Customer().GetUUIDs(arm.AuthContext(authtoken), &customer.UUIDRequest{
	// 	Service: wrapperspb.String(service),
	// })

	// if err != nil {
	// 	return false
	// }

	// for _, stru := range struuids {
	// 	uuid, err := strconv.ParseUint(stru, 10, 64)
	// 	if err != nil {
	// 		return false
	// 	}

	// 	for _, u := range res.Uuids {
	// 		if u == uuid {
	// 			return true
	// 		}
	// 	}
	// }

	// return false
	return true
}

func ValidateToken(c *gin.Context) {
	// authHeader := c.GetHeader("Authorization")
	// tokenExists, tokenHeader := lib.CheckForBearerToken(authHeader)
	// if tokenExists == true {
	// 	c.Set("auth_token", tokenHeader)
	// } else {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"message": "Unauthorized Access",
	// 	})
	// }
	c.Set("auth_token", true)
}

func canAccess(c *gin.Context) {
	// authHeader := c.GetHeader("Authorization")
	// _, tokenHeader := lib.CheckForBearerToken(authHeader)

	// uuid := c.Param("uuid")
	// service := c.Param("service")
	// if service == "" {
	// 	service = c.MustGet("service").(string)
	// }
	// canAccess := canAccessServiceUUID(service, []string{uuid}, tokenHeader)
	canAccess := true
	if !canAccess {
		c.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{
			"error": "You do not have permission to this UUID.",
		})
	}

}

func AddServiceParm(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("service", service)
		c.Next()
	}
}

func main() {
	// connect to arm
	// if err := arm.Connect(); err != nil {
	// 	panic(err)
	// }

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
	}))
	r.Use(gin.Recovery(), gzip.Gzip(gzip.DefaultCompression))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	v1 := r.Group("/v1")
	v1.Use(ValidateToken)
	v1.GET("/tiles/:service/:uuid/:timestamp/:z/:x/:y", canAccess, api.GetTiles)
	v1.GET("/summary/:service/:uuid/:timestamp", canAccess, api.GetSummary)
	v1.GET("/summaries", api.GetSummaries)

	ldm := v1.Group("/land-subsidence").Use(AddServiceParm("land-subsidence"))
	{
		ldm.GET("/summary/:uuid/:type/:timestamp/:lyrtype", canAccess, api.GetLDMSummary)
		ldm.GET("/tiles/:uuid/:type/:timestamp/:lyrtype/:z/:x/:y", canAccess, api.GetLDMTiles)
		ldm.GET("/summaries", api.GetSummaries)
	}

	database.Init()
	r.Run(":7778")
	database.Close()
}
