package webapi

import (
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/database"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func profile(c *gin.Context) {
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
