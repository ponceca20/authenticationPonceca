package customer

import (
	"practicev2/module/authentication/models"
)

// CustomerProfileDTO is a public representation of a customer's full profile.
type CustomerProfileDTO struct {
	ID               string                   `json:"id"`
	IdentityID       string                   `json:"identity_id"`
	Email            string                   `json:"email"`
	FirstName        string                   `json:"first_name"`
	LastName         string                   `json:"last_name"`
	Phone            string                   `json:"phone"`
	CustomerNumber   string                   `json:"customer_number"`
	AcceptsMarketing bool                     `json:"accepts_marketing"`
	TotalSpent       float64                  `json:"total_spent"`
	LoyaltyPoints    int                      `json:"loyalty_points"`
	Addresses        []models.ShippingAddress `json:"addresses,omitempty"`
	// Preferences    models.CustomerPreferences `json:"preferences,omitempty"`
}

// CustomerRegistrationDTO defines the data for registering a new e-commerce customer.
type CustomerRegistrationDTO struct {
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	FirstName        string `json:"first_name" validate:"required"`
	LastName         string `json:"last_name" validate:"required"`
	Phone            string `json:"phone,omitempty"`
	AcceptsMarketing bool   `json:"accepts_marketing"`
}

// AddressDTO defines the data for creating or updating a shipping address.
type AddressDTO struct {
	AddressLine1 string `json:"address_line_1" validate:"required"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	PostalCode   string `json:"postal_code" validate:"required"`
	Country      string `json:"country" validate:"required"`
	IsDefault    bool   `json:"is_default"`
}

// ToCustomerProfileDTO converts a CustomerProfile and its associated Identity to a public DTO.
func ToCustomerProfileDTO(identity *models.Identity, profile *models.CustomerProfile) CustomerProfileDTO {
	return CustomerProfileDTO{
		ID:               profile.ID,
		IdentityID:       identity.ID,
		Email:            identity.Email,
		FirstName:        identity.FirstName,
		LastName:         identity.LastName,
		Phone:            identity.Phone,
		CustomerNumber:   profile.CustomerNumber,
		AcceptsMarketing: profile.AcceptsMarketing,
		TotalSpent:       profile.TotalSpent,
		LoyaltyPoints:    profile.LoyaltyPoints,
	}
}

// PreferencesDTO is used for updating customer preferences.
type PreferencesDTO struct {
	Theme              string `json:"theme" validate:"oneof=light dark"`
	Language           string `json:"language" validate:"oneof=en es"`
	EmailNotifications bool   `json:"email_notifications"`
	SmsNotifications   bool   `json:"sms_notifications"`
}
