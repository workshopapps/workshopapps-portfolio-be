package routes

import (
	"test-va/cmd/handlers/userHandler"
	"test-va/cmd/middlewares"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/internals/service/userService"

	"github.com/gin-gonic/gin"
)

func UserRoutes(v1 *gin.RouterGroup, srv userService.UserSrv, tokenSrv tokenservice.TokenSrv) {
	userHandler := userHandler.NewUserHandler(srv)
	jwtMWare := middlewares.NewJWTMiddleWare(tokenSrv)

	// Register a user

	v1.POST("/user", userHandler.CreateUser)
	// Login into the user account
	v1.POST("/user/login", userHandler.Login)
	// Get a reset password token
	v1.POST("/user/reset-password", userHandler.ResetPassword)
	// Reset password with token id
	v1.POST("/user/reset-password-token", userHandler.ResetPasswordWithToken)

	users := v1.Group("/user")
	settings := users.Group("/settings")
	users.Use(jwtMWare.ValidateJWT())
	{
		// Get all users
		users.GET("", userHandler.GetUsers)
		// Get a specific user
		users.GET("/:user_id", userHandler.GetUser)
		// Update a specific user
		users.PATCH("/:user_id", userHandler.UpdateUser)
		// Update user image
		users.POST("/upload", userHandler.UploadImage)
		// Change user password
		users.PUT("/change-password", userHandler.ChangePassword)
		// Delete a user
		users.DELETE("/:user_id", userHandler.DeleteUser)
		// Assign VA to User
		users.POST("/assign-va/:va_id", userHandler.AssignVAToUser)

	}
	settings.Use(jwtMWare.ValidateJWT())
	{
		//get settings
		settings.GET("/", userHandler.GetSettings)

		//update reminder settings
		settings.PATCH("/reminder-settings", userHandler.UpdateReminderSettings)
		settings.PATCH("/notification-settings", userHandler.UpdateNotificationSettings)
		settings.PATCH("/product-email-settings", userHandler.UpdateProductEmailSettings)

	}

}
