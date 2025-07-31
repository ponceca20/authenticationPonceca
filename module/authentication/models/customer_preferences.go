package models

import "time"

// CustomerPreferences stores various preferences for a customer.
type CustomerPreferences struct {
	ID                 string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CustomerProfileID  string          `json:"customer_profile_id" gorm:"uniqueIndex;type:varchar(36);not null"`
	CustomerProfile    CustomerProfile `json:"customer_profile" gorm:"foreignKey:CustomerProfileID"`
	Theme              string          `json:"theme,omitempty" gorm:"default:'light';size:50"`
	Language           string          `json:"language,omitempty" gorm:"default:'es';size:10"`
	TimeZone           string          `json:"time_zone,omitempty" gorm:"default:'UTC';size:50"`
	EmailNotifications bool            `json:"email_notifications" gorm:"default:true"`
	SmsNotifications   bool            `json:"sms_notifications" gorm:"default:false"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// TableName overrides the table name used by CustomerPreferences to ensure singular naming
func (CustomerPreferences) TableName() string {
	return "customer_preference"
}
