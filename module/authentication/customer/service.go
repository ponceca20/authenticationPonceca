package customer

import (
	"errors"
	"fmt"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerService defines the interface for customer-related business logic.
type CustomerService interface {
	RegisterCustomer(dto *CustomerRegistrationDTO) (*models.CustomerProfile, error)
	GetCustomerProfile(identityID string) (*CustomerProfileDTO, error)
	AddAddress(identityID string, dto *AddressDTO) (*models.ShippingAddress, error)
	ListAddresses(identityID string) ([]models.ShippingAddress, error)
	UpdateAddress(identityID, addressID string, dto *AddressDTO) (*models.ShippingAddress, error)
	DeleteAddress(identityID, addressID string) error
	GetPreferences(identityID string) (*models.CustomerPreferences, error)
	UpdatePreferences(identityID string, dto *PreferencesDTO) (*models.CustomerPreferences, error)
}

type customerService struct {
	customerRepo CustomerRepository
	authRepo     auth.AuthRepository // To check for existing identities
}

// NewCustomerService creates a new instance of CustomerService.
func NewCustomerService(customerRepo CustomerRepository, authRepo auth.AuthRepository) CustomerService {
	return &customerService{customerRepo: customerRepo, authRepo: authRepo}
}

// RegisterCustomer handles the logic for creating a new e-commerce customer.
func (s *customerService) RegisterCustomer(dto *CustomerRegistrationDTO) (*models.CustomerProfile, error) {
	// Validate DTO (basic validation)
	if err := utils.ValidateStruct(dto); err != nil {
		// In a real app, you'd format this error nicely.
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// Check if identity already exists
	if _, err := s.authRepo.FindIdentityByEmail(dto.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("email is already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Prepare models
	identity := &models.Identity{
		ID:           uuid.New().String(),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		PasswordHash: hashedPassword,
	}

	profile := &models.CustomerProfile{
		ID:               uuid.New().String(),
		CustomerNumber:   fmt.Sprintf("CUS-%s", uuid.New().String()[:8]),
		AcceptsMarketing: dto.AcceptsMarketing,
	}

	// Create customer in a transaction
	if err := s.customerRepo.CreateCustomer(identity, profile); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return profile, nil
}

// GetCustomerProfile retrieves the full profile for a given identity ID.
func (s *customerService) GetCustomerProfile(identityID string) (*CustomerProfileDTO, error) {
	identity, err := s.authRepo.FindIdentityByID(identityID)
	if err != nil {
		return nil, errors.New("identity not found")
	}

	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	// Optionally, load addresses and preferences here
	// addresses, _ := s.customerRepo.ListShippingAddresses(profile.ID)
	// preferences, _ := s.customerRepo.GetPreferences(profile.ID)

	// Load addresses and add them to the DTO
	addresses, _ := s.ListAddresses(identityID)
	dto := ToCustomerProfileDTO(identity, profile)
	dto.Addresses = addresses

	return &dto, nil
}

// AddAddress adds a new shipping address for a customer.
func (s *customerService) AddAddress(identityID string, dto *AddressDTO) (*models.ShippingAddress, error) {
	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	address := &models.ShippingAddress{
		ID:              uuid.New().String(),
		CustomerProfileID: profile.ID,
		AddressLine1:    dto.AddressLine1,
		AddressLine2:    dto.AddressLine2,
		City:            dto.City,
		State:           dto.State,
		PostalCode:      dto.PostalCode,
		Country:         dto.Country,
		IsDefault:       dto.IsDefault,
	}

	if err := s.customerRepo.AddShippingAddress(address); err != nil {
		return nil, fmt.Errorf("failed to add address: %w", err)
	}
	return address, nil
}

// ListAddresses lists all shipping addresses for a customer.
func (s *customerService) ListAddresses(identityID string) ([]models.ShippingAddress, error) {
	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}
	return s.customerRepo.ListShippingAddresses(profile.ID)
}

// UpdateAddress updates a shipping address.
func (s *customerService) UpdateAddress(identityID, addressID string, dto *AddressDTO) (*models.ShippingAddress, error) {
	// First, ensure the address belongs to the customer making the request.
	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	address, err := s.customerRepo.FindAddressByID(addressID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	if address.CustomerProfileID != profile.ID {
		return nil, errors.New("address does not belong to this customer")
	}

	// Update fields
	address.AddressLine1 = dto.AddressLine1
	address.AddressLine2 = dto.AddressLine2
	address.City = dto.City
	address.State = dto.State
	address.PostalCode = dto.PostalCode
	address.Country = dto.Country
	address.IsDefault = dto.IsDefault

	if err := s.customerRepo.UpdateAddress(address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}
	return address, nil
}

// DeleteAddress deletes a shipping address.
func (s *customerService) DeleteAddress(identityID, addressID string) error {
	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return errors.New("customer profile not found")
	}

	address, err := s.customerRepo.FindAddressByID(addressID)
	if err != nil {
		// If not found, it's already gone.
		return nil
	}

	if address.CustomerProfileID != profile.ID {
		return errors.New("address does not belong to this customer")
	}

	return s.customerRepo.DeleteAddress(addressID)
}

// GetPreferences retrieves a customer's preferences.
func (s *customerService) GetPreferences(identityID string) (*models.CustomerPreferences, error) {
	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	prefs := &models.CustomerPreferences{CustomerProfileID: profile.ID}
	if err := s.customerRepo.FindOrCreatePreferences(prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}

// UpdatePreferences updates a customer's preferences.
func (s *customerService) UpdatePreferences(identityID string, dto *PreferencesDTO) (*models.CustomerPreferences, error) {
	profile, err := s.customerRepo.FindProfileByIdentityID(identityID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	prefs := &models.CustomerPreferences{CustomerProfileID: profile.ID}
	if err := s.customerRepo.FindOrCreatePreferences(prefs); err != nil {
		return nil, err
	}

	// Update fields
	prefs.Theme = dto.Theme
	prefs.Language = dto.Language
	prefs.EmailNotifications = dto.EmailNotifications
	prefs.SmsNotifications = dto.SmsNotifications

	if err := s.customerRepo.UpdatePreferences(prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}
