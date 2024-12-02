package webapi

import (
	"fmt"

	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/database"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/models"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/ollama"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Payload struct {
	Model string `json:"model" binding:"required"`
}

func generate(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func chat(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func pull(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func push(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func create(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func list(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	ollamaHost := setOllamaHost(user.WorkspaceName)
	models, err := ollama.ListModels(ollamaHost)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	c.JSON(200, models)
}

func listRunning(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	ollamaHost := setOllamaHost(user.WorkspaceName)
	models, err := ollama.ListRunningModels(ollamaHost)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	c.JSON(200, models)
}

func copy(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func delete(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func show(c *gin.Context) {
	var user models.User
	var payload Payload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs",
		})
		c.Abort()
		return
	}

	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	// fmt.Println("payload: ", payload)
	// ollamaHost := setOllamaHost(user.WorkspaceName)
	// model := ollama.ShowModel(payload.Model, ollamaHost)
	// fmt.Println("model: ", model)
	user.Password = ""

	c.JSON(200, user)
}

func heartBeat(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func embed(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func embeddings(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func createBlob(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func apiVersion(c *gin.Context) {
	var user models.User
	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get User Profile",
		})
		c.Abort()
		return
	}

	user.Password = ""

	c.JSON(200, user)
}

func setOllamaHost(workspace string) string {
	serviceName := fmt.Sprintf("%s-ollama", workspace)
	ollamaPort := int64(11434)
	ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, workspace, ollamaPort)
	return ollamaServerURI
}
