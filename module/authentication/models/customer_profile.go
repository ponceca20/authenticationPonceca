package models

import "time"

// CustomerProfile holds e-commerce specific data for an Identity.
type CustomerProfile struct {
	ID             string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID     string   `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity       Identity `json:"identity" gorm:"foreignKey:IdentityID"`

	// E-commerce specific data
	CustomerNumber         string     `json:"customer_number" gorm:"uniqueIndex;size:191"`
	PreferredPaymentMethod string     `json:"preferred_payment_method,omitempty"`
	CreditLimit            float64    `json:"credit_limit" gorm:"default:0"`
	TotalSpent             float64    `json:"total_spent" gorm:"default:0"`
	LoyaltyPoints          int        `json:"loyalty_points" gorm:"default:0"`

	// Marketing preferences
	AcceptsMarketing  bool   `json:"accepts_marketing" gorm:"default:false"`
	PreferredLanguage string `json:"preferred_language" gorm:"default:'es'"`

	// Standard model fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
