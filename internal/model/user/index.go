package user

import (
	"github.com/yockii/qscore/pkg/common"
	"github.com/yockii/qscore/pkg/database"
)

// User 用户表
type User struct {
	common.BaseModel
	Email                  string `json:"email" gorm:"uniqueIndex:uni_user_email_site;size:100;comment:邮箱"`
	SiteID                 uint64 `json:"siteId,string" gorm:"uniqueIndex:uni_user_email_site;comment:站点ID"`
	CreatedAt              int64  `json:"createdAt" gorm:"autoCreateTime:milli"`
	UpdatedAt              int64  `json:"updatedAt" gorm:"autoUpdateTime:milli"`
	EncryptedPassword      string `json:"encryptedPassword" gorm:"size:100;comment:加密密码"`
	ResetPasswordToken     string `json:"resetPasswordToken" gorm:"uniqueIndex;size:100;comment:重置密码token"`
	ResetPasswordSentAt    int64  `json:"resetPasswordSentAt" gorm:"comment:重置密码发送时间"`
	RememberCreatedAt      int64  `json:"rememberCreatedAt" gorm:"comment:记住创建时间"`
	SignInCount            int64  `json:"signInCount" gorm:"comment:登录次数"`
	CurrentSignInAt        int64  `json:"currentSignInAt" gorm:"comment:当前登录时间"`
	LastSignInAt           int64  `json:"lastSignInAt" gorm:"comment:上次登录时间"`
	CurrentSignInIp        string `json:"currentSignInIp" gorm:"size:100;comment:当前登录IP"`
	Admin                  *bool  `json:"admin" gorm:"comment:是否管理员"`
	ConfirmationToken      string `json:"confirmationToken" gorm:"uniqueIndex;size:100;comment:确认token"`
	ConfirmedAt            int64  `json:"confirmedAt" gorm:"comment:确认时间"`
	ConfirmationSentAt     int64  `json:"confirmationSentAt" gorm:"comment:确认发送时间"`
	UnconfirmedEmail       string `json:"unconfirmedEmail" gorm:"size:100;comment:未确认邮箱"`
	Locale                 string `json:"locale" gorm:"size:100;comment:语言"`
	EncryptedOtpSecret     string `json:"encryptedOtpSecret" gorm:"size:100;comment:加密OTP密钥"`
	EncryptedOtpSecretIv   string `json:"encryptedOtpSecretIv" gorm:"size:100;comment:加密OTP密钥IV"`
	EncryptedOtpSecretSalt string `json:"encryptedOtpSecretSalt" gorm:"size:100;comment:加密OTP密钥盐"`
	ConsumedTimestep       int64  `json:"consumedTimestep" gorm:"comment:消耗时间步长"`
	OtpRequiredForLogin    *bool  `json:"otpRequiredForLogin" gorm:"comment:登录是否需要OTP"`
	LastEmailedAt          int64  `json:"lastEmailedAt" gorm:"comment:上次发送邮件时间"`
	OtpBackupCodes         string `json:"otpBackupCodes" gorm:"size:1000;comment:OTP备份码"`
	FilteredLanguages      string `json:"filteredLanguages" gorm:"size:1000;comment:过滤语言"`
	AccountID              uint64 `json:"accountId" gorm:"index;comment:账户ID"`
	Disabled               *bool  `json:"disabled" gorm:"comment:是否禁用"`
	Moderator              *bool  `json:"moderator" gorm:"comment:是否是版主"`
	InviteID               uint64 `json:"inviteId" gorm:"comment:邀请ID"`
	RememberToken          string `json:"rememberToken" gorm:"uniqueIndex;size:100;comment:记住token"`
	ChosenLanguages        string `json:"chosenLanguages" gorm:"size:1000;comment:选择语言"`
	CreatedByApplicationID uint64 `json:"createdByApplicationId" gorm:"index;comment:创建应用ID"`
	Approved               *bool  `json:"approved" gorm:"comment:是否已批准"`
}

func (*User) TableComment() string {
	return `本应用用户表`
}

func (u *User) AddRequired() string {
	if u.Email == "" || u.EncryptedPassword == "" {
		return "email, encryptedPassword"
	}
	return ""
}

func (u *User) CheckDuplicatedModel() database.Model {
	return &User{
		Email: u.Email,
	}
}

func (u *User) UpdateModel() database.Model {
	return &User{
		Email:                  u.Email,
		EncryptedPassword:      u.EncryptedPassword,
		ResetPasswordToken:     u.ResetPasswordToken,
		ResetPasswordSentAt:    u.ResetPasswordSentAt,
		RememberCreatedAt:      u.RememberCreatedAt,
		SignInCount:            u.SignInCount,
		CurrentSignInAt:        u.CurrentSignInAt,
		LastSignInAt:           u.LastSignInAt,
		CurrentSignInIp:        u.CurrentSignInIp,
		Admin:                  u.Admin,
		ConfirmationToken:      u.ConfirmationToken,
		ConfirmedAt:            u.ConfirmedAt,
		ConfirmationSentAt:     u.ConfirmationSentAt,
		UnconfirmedEmail:       u.UnconfirmedEmail,
		Locale:                 u.Locale,
		EncryptedOtpSecret:     u.EncryptedOtpSecret,
		EncryptedOtpSecretIv:   u.EncryptedOtpSecretIv,
		EncryptedOtpSecretSalt: u.EncryptedOtpSecretSalt,
		ConsumedTimestep:       u.ConsumedTimestep,
		OtpRequiredForLogin:    u.OtpRequiredForLogin,
		LastEmailedAt:          u.LastEmailedAt,
		OtpBackupCodes:         u.OtpBackupCodes,
		FilteredLanguages:      u.FilteredLanguages,
		AccountID:              u.AccountID,
		Disabled:               u.Disabled,
		Moderator:              u.Moderator,
		InviteID:               u.InviteID,
		RememberToken:          u.RememberToken,
		ChosenLanguages:        u.ChosenLanguages,
		CreatedByApplicationID: u.CreatedByApplicationID,
		Approved:               u.Approved,
	}
}

func (u *User) FuzzyQueryMap() map[string]string {
	result := make(map[string]string)
	if u.Email != "" {
		result["email"] = "%" + u.Email + "%"
	}
	return result
}

func (u *User) ListOmits() string {
	return "encryptedPassword,resetPasswordToken,resetPasswordSentAt,rememberCreatedAt,signInCount,currentSignInAt,lastSignInAt,currentSignInIp,confirmationToken,confirmedAt,confirmationSentAt,unconfirmedEmail,locale,encryptedOtpSecret,encryptedOtpSecretIv,encryptedOtpSecretSalt,consumedTimestep,otpRequiredForLogin,lastEmailedAt,otpBackupCodes,filteredLanguages,disabled,inviteId,rememberToken,chosenLanguages,createdByApplicationId,approved"
}

func init() {
	database.Models = append(database.Models, &User{})
}
