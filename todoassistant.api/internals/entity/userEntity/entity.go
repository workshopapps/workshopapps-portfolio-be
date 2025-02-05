package userEntity

type CreateUserReq struct {
	UserId        string `json:"user_id"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name"  validate:"required"`
	Email         string `json:"email" validate:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password" validate:"required,min=6"`
	Gender        string `json:"gender"`
	DateOfBirth   string `json:"date_of_birth"`
	AccountStatus string `json:"account_status"`
	PaymentStatus string `json:"payment_status"`
	DateCreated   string `json:"date_created"`
}

type CreateUserRes struct {
	UserId       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type NotificationSettingsRes struct {
	NewComments     bool `json:"new_comment"`
	ExpiredTasks    bool `json:"expired_tasks"`
	ReminderTasks   bool `json:"reminder_tasks"`
	VaAcceptingTask bool `json:"va_accepting_task"`
	TaskAssingnedVa bool `json:"task_assigned"`
	Subscribtion    bool `json:"subscription"`
}

type ProductEmailSettingsRes struct {
	NewProducts        bool `json:"new_product"`
	LoginAlert         bool `json:"login_alert"`
	PromotionAndOffers bool `json:"promotions_and_offers"`
	TipsDailyDigest    bool `json:"tips_daily_digest"`
}

type NotificationSettingsReq struct {
	UserId          string `json:"user_id"`
	NewComments     bool   `json:"new_comment"`
	ExpiredTasks    bool   `json:"expired_tasks"`
	ReminderTasks   bool   `json:"reminder_tasks"`
	VaAcceptingTask bool   `json:"va_accepting_task"`
	TaskAssingnedVa bool   `json:"task_assigned"`
	Subscribtion    bool   `json:"subscription"`
}

type ProductEmailSettingsReq struct {
	NewProducts        bool `json:"new_product"`
	LoginAlert         bool `json:"login_alert"`
	PromotionAndOffers bool `json:"promotions_and_offers"`
	TipsDailyDigest    bool `json:"tips_daily_digest"`
}
type ReminderSettingsReq struct {
	RemindMeVia  string `json:"remind_me_via"`
	WhenSnooze   string `json:"when_snooze"`
	AutoReminder string `json:"auto_reminder"`
	ReminderTime string `json:"reminder_time"`
	Refresh      string `json:"refresh"`
}

type UserSettingsRes struct {
	ReminderSettings     ReminderSettingsRes     `json:"reminder_settings"`
	NotificationSettings NotificationSettingsRes `json:"notification_settings"`
	ProductEmailSettings ProductEmailSettingsRes `json:"product_email_settings"`
}
type ReminderSettingsRes struct {
	RemindMeVia  string `json:"remind_me_via"`
	WhenSnooze   string `json:"when_snooze"`
	AutoReminder string `json:"auto_reminder"`
	ReminderTime string `json:"reminder_time"`
	Refresh      string `json:"refresh"`
}

type LoginRes struct {
	UserId               string                  `json:"user_id"`
	Email                string                  `json:"email"`
	FirstName            string                  `json:"first_name"`
	LastName             string                  `json:"last_name"`
	Phone                string                  `json:"phone"`
	Gender               string                  `json:"gender"`
	Avatar               string                  `json:"avatar"`
	CountryId            int                     `json:"country_id"`
	Occupation           string                  `json:"occupation"`
	NotificationSettings NotificationSettingsRes `json:"notification_settings"`
	ProductEmailSettings ProductEmailSettingsRes `json:"product_email_settings"`
	Token                string                  `json:"access_token"`
	RefreshToken         string                  `json:"refresh_token"`
}

type GetByEmailRes struct {
	UserId     string `json:"user_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
	Avatar     string `json:"avatar"`
	CountryId  int    `json:"country_id"`
	Occupation string `json:"occupation"`
}

type GetByIdRes struct {
	UserId      string `json:"user_id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"date_of_birth"`
}

type UpdateUserReq struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Gender        string `json:"gender"`
	DateOfBirth   string `json:"date_of_birth"`
	Avatar        string `json:"avatar"`
	AccountStatus string `json:"account_status"`
	PaymentStatus string `json:"payment_status"`
	CountryId     int    `json:"country_id"`
	Occupation    string `json:"occupation"`
}

type UpdateUserRes struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
	Avatar     string `json:"avatar"`
	CountryId  int    `json:"country_id"`
	Occupation string `json:"occupation"`
}

type UsersRes struct {
	UserId      string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
	DateCreated string `json:"date_created"`
}

type ChangePasswordReq struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordReq struct {
	Email string `json:"email" validate:"email,required"`
}

type ResetPasswordRes struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	Token   string `json:"token"`
	Expiry  string `json:"expiry"`
}

type ResetPasswordWithTokenReq struct {
	Password string `json:"password" validate:"required"`
}

type ResetPasswordWithTokenRes struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	Token   string `json:"token"`
	Expiry  string `json:"expiry"`
}

type GoogleLoginReq struct {
	Id        string `json:"googleId"`
	FirstName string `json:"givenName" binding:"required"`
	LastName  string `json:"familyName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Profile   string `json:"imageUrl"`
	Name      string `json:"name" binding:"required"`
}

type FacebookLoginReq struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

type ProfileImageRes struct {
	Image    string `json:"avatar"`
	Size     int64  `json:"size"`
	FileType string `json:"fileType"`
}
