package webapi

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/database"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/middlewares"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/models"

	"github.com/gin-gonic/gin"
)

func StartWebAPI() {
	ctx := context.Background()
	logger := log.FromContext(ctx)

	err := database.InitDatabase()
	if err != nil {
		logger.Error(err, "could not create database", "error", err)
	}

	database.GlobalDB.AutoMigrate(&models.User{})
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome To AIChat Workspace")
	})

	api := r.Group("/api")
	{
		public := api.Group("/public")
		{
			public.POST("/login", login)
			public.POST("/signup", signup)
		}

		protected := api.Group("/protected").Use(middlewares.Authz())
		{
			protected.GET("/profile", profile)
			// protected.POST("/workspace", aichatWorkspace)
		}

		ollama := api.Group("/ollama").Use(middlewares.Authz())
		{
			ollama.GET("/version", apiVersion)

			// POST Generate
			ollama.POST("/generate", generate)

			// POST Chat
			ollama.POST("/chat", chat)

			// POST Pull
			ollama.POST("/pull", pull)

			// POST Push
			ollama.POST("/push", push)

			// POST Create
			ollama.POST("/create", create)

			// GET List
			ollama.GET("/list", list)

			// GET ListRunning
			ollama.GET("/listRunning", listRunning)

			// POST Copy
			ollama.POST("/copy", copy)
			// DELETE Delete

			// POST Show
			ollama.POST("/show", show)
			// HEAD Heartbeat

			// POST Embed
			ollama.POST("/embed", embed)

			// POST Embeddings
			ollama.POST("/embeddings", embeddings)

			// POST CreateBlob
			ollama.POST("/createBlob", createBlob)
		}
	}

	return r
}
