package webapi

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/dushu1105/jt808/jtnet"
)

const SERVER_KEY = "808server"
func addTcpServerInContext(s *jtnet.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(SERVER_KEY, s)
		c.Next()
	}
}

func RunWebServer(s *jtnet.Server) {
	fmt.Println("Starting web api system...")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "x-requested-with", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(addTcpServerInContext(s))
	//r.POST("/create_rtc", CreateTalk)
	r.POST("/jt808_8107", Jt808_8107)
	r.POST("/jt808_8300", Jt808_8300)
	r.POST("/jt808_8702", Jt808_8702)
	r.POST("/jt808_9101", Jt808_9101)
	r.POST("/jt808_9201", Jt808_9201)
	r.POST("/jt808_9206", Jt808_9206)
	//r.Static("/static/", "./static/")
	r.Run(":8089")
}
