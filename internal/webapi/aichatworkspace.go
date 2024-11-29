package webapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/database"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/models"
)

func aichatWorkspace(c *gin.Context) {
	// var k8sClient client.Client
	var user models.User
	// ctx := context.Background()

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

	// resourceName := fmt.Sprintf("%s-%s", "aichatworkspace", user.Name)
	// resource := &appsv1alpha1.AIChatWorkspace{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      resourceName,
	// 		Namespace: "default",
	// 	},
	// 	// TODO(user): Specify other spec details if needed.
	// 	Spec: appsv1alpha1.AIChatWorkspaceSpec{
	// 		WorkspaceName: user.WorkspaceName,
	// 		WorkspaceEnv:  "dev",
	// 		Models:        []string{"gemma2:2b", "qwen2.5-coder:3b"},
	// 	},
	// }

	// err := k8sClient.Create(ctx, resource)
	// if err != nil {
	// 	c.JSON(500, gin.H{
	// 		"Error": "Could Not Get User Profile",
	// 	})
	// 	c.Abort()
	// 	return
	// }

	user.Password = ""

	c.JSON(200, user)
}
