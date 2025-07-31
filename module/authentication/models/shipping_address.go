package models

import "time"

// ShippingAddress represents a shipping address for a customer.
type ShippingAddress struct {
	ID              string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CustomerProfileID string          `json:"customer_profile_id" gorm:"index;type:varchar(36)"`
	CustomerProfile CustomerProfile `json:"customer_profile" gorm:"foreignKey:CustomerProfileID"`
	AddressLine1    string          `json:"address_line_1" gorm:"not null"`
	AddressLine2    string          `json:"address_line_2,omitempty"`
	City            string          `json:"city" gorm:"not null"`
	State           string          `json:"state" gorm:"not null"`
	PostalCode      string          `json:"postal_code" gorm:"not null"`
	Country         string          `json:"country" gorm:"not null"`
	IsDefault       bool            `json:"is_default" gorm:"default:false"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty" gorm:"index"`
}
