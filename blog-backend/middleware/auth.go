package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/task/go_learn_task/blog-backend/config"
	"github.com/task/go_learn_task/blog-backend/database"
	"github.com/task/go_learn_task/blog-backend/models"
	"github.com/task/go_learn_task/blog-backend/utils"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			utils.ErrorResponse(c, 404, "User not found", err)
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Next()
	}
}
