package main

import (
	"bland/controller"
	_ "bland/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey bearerToken
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	{
        // Define the route for sending calls 
		v1.POST("/call", controller.SendCall)
		// Define the route for analyzing call
		v1.POST("call/:call_id/analyze", controller.AnalyzeCall)
		// Define the route for getting call details
	    v1.GET("/calls/:call_id", controller.GetCallDetails)
		// Define the route for creating a folder
	   v1.POST("/folders", controller.CreateFolder)
	   // Define the route that creates the pathway and move to specfic folder
	   v1.POST("/pathways/create-and-move", controller.CreateAndMovePathway)
	   // Define the route for creating a chat to test AI bots
	   v1.POST("/pathways/chat/create", controller.CreateChat)
	   v1.GET("/convo_pathway/:pathway_id", controller.GetPathwayInfo)
	   v1.POST("/pathway/update/:pathway_id", controller.UpdatePathway)
	   v1.DELETE("/delete/convo_pathway/:pathway_id", controller.DeletePathway)
	   v1.POST("/pathways/chat/:chat_id/send", controller.SendMessageToChat)
    }

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
