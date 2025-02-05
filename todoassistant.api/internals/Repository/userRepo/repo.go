package userRepo

import (
	"test-va/internals/entity/userEntity"
)

type UserRepository interface {
	GetUsers(page int) ([]*userEntity.UsersRes, error)
	Persist(req *userEntity.CreateUserReq) error
	GetByEmail(email string) (*userEntity.GetByEmailRes, error)
	GetById(user_id string) (*userEntity.GetByIdRes, error)
	UpdateUser(req *userEntity.UpdateUserReq, userId string) error
	UpdateImage(userId, fileName string) error
	DeleteUser(user_id string) error
	ChangePassword(userId, newPassword string) error
	AddToken(req *userEntity.ResetPasswordRes) error
	GetTokenById(token, userId string) (*userEntity.ResetPasswordWithTokenRes, error)
	DeleteToken(tokenId string) error
	AssignVAToUser(user_id, token_id string) error
	//user settings functions
	GetNotificationSettingsById(userId string) (*userEntity.NotificationSettingsRes, error)
	GetProductEmailSettingsById(userId string) (*userEntity.ProductEmailSettingsRes, error)

	//set reminder settings
	SetReminderSettings(req *userEntity.ReminderSettingsReq, userId string) error
	GetReminderSettings(userId string) (*userEntity.ReminderSettingsRes, error)
	UpdateReminderSettings(req *userEntity.ReminderSettingsReq, userId string) error
	UpdateProductEmailSettings(req *userEntity.ProductEmailSettingsReq, userId string) error
	UpdateNotificationSettings(req *userEntity.NotificationSettingsReq, userId string) error
	// GetUserSettings(userId string) (*userEntity.UserSettingsRes, error)
}
