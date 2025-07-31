package customer

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// CustomerRepository defines the interface for database operations related to customers.
type CustomerRepository interface {
	CreateCustomer(identity *models.Identity, profile *models.CustomerProfile) error
	FindProfileByIdentityID(identityID string) (*models.CustomerProfile, error)
	AddShippingAddress(address *models.ShippingAddress) error
	ListShippingAddresses(profileID string) ([]models.ShippingAddress, error)
	FindAddressByID(addressID string) (*models.ShippingAddress, error)
	UpdateAddress(address *models.ShippingAddress) error
	DeleteAddress(addressID string) error
	FindOrCreatePreferences(prefs *models.CustomerPreferences) error
	UpdatePreferences(prefs *models.CustomerPreferences) error
}

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new instance of CustomerRepository.
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// CreateCustomer creates a new identity and a customer profile in a single transaction.
func (r *customerRepository) CreateCustomer(identity *models.Identity, profile *models.CustomerProfile) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(identity).Error; err != nil {
			return err
		}

		profile.IdentityID = identity.ID
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindProfileByIdentityID finds a customer profile using the associated identity's ID.
func (r *customerRepository) FindProfileByIdentityID(identityID string) (*models.CustomerProfile, error) {
	var profile models.CustomerProfile
	if err := r.db.Where("identity_id = ?", identityID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// AddShippingAddress adds a new shipping address for a customer.
func (r *customerRepository) AddShippingAddress(address *models.ShippingAddress) error {
	return r.db.Create(address).Error
}

// ListShippingAddresses retrieves all shipping addresses for a given customer profile ID.
func (r *customerRepository) ListShippingAddresses(profileID string) ([]models.ShippingAddress, error) {
	var addresses []models.ShippingAddress
	if err := r.db.Where("customer_profile_id = ?", profileID).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

// FindAddressByID finds a shipping address by its ID.
func (r *customerRepository) FindAddressByID(addressID string) (*models.ShippingAddress, error) {
	var address models.ShippingAddress
	if err := r.db.Where("id = ?", addressID).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// UpdateAddress saves changes to a shipping address.
func (r *customerRepository) UpdateAddress(address *models.ShippingAddress) error {
	return r.db.Save(address).Error
}

// DeleteAddress soft deletes a shipping address.
func (r *customerRepository) DeleteAddress(addressID string) error {
	return r.db.Delete(&models.ShippingAddress{}, "id = ?", addressID).Error
}

// FindOrCreatePreferences finds customer preferences or creates them if they don't exist.
func (r *customerRepository) FindOrCreatePreferences(prefs *models.CustomerPreferences) error {
	return r.db.Where(models.CustomerPreferences{CustomerProfileID: prefs.CustomerProfileID}).FirstOrCreate(prefs).Error
}

// UpdatePreferences saves changes to customer preferences.
func (r *customerRepository) UpdatePreferences(prefs *models.CustomerPreferences) error {
	return r.db.Save(prefs).Error
}
