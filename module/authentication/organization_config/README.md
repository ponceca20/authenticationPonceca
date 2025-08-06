# Organization Module Configuration

This module implements dynamic configuration management for organization-specific modules, allowing each organization to have customized RBAC permissions and module settings.

## Features

- **Dynamic Module Configuration**: Each organization can have different modules enabled/disabled
- **Per-Organization RBAC**: Customizable permissions for each organization
- **Clean Architecture**: Service, Repository, Handler, DTO pattern
- **RESTful API**: Complete CRUD operations for module configurations

## API Endpoints

### Module Configuration Management
- `GET /api/organization-config/modules` - List all module configurations for an organization
- `POST /api/organization-config/modules/configure` - Configure modules for an organization
- `GET /api/organization-config/modules/:id` - Get specific module configuration
- `PUT /api/organization-config/modules/:id` - Update module configuration
- `DELETE /api/organization-config/modules/:id` - Delete module configuration

### Module Operations
- `POST /api/organization-config/modules/:id/enable` - Enable a module
- `POST /api/organization-config/modules/:id/disable` - Disable a module
- `GET /api/organization-config/modules/:id/permissions` - Get module permissions

### Organization Settings
- `GET /api/organization-config/organization/:orgId/modules` - Get modules for organization
- `POST /api/organization-config/organization/:orgId/setup-expenses` - Setup expenses module

## Models

### OrganizationModuleConfig
```go
type OrganizationModuleConfig struct {
    ID              uint                   `gorm:"primaryKey"`
    OrganizationID  uint                  `gorm:"not null;index"`
    ModuleName      string                `gorm:"not null;size:50"`
    IsEnabled       bool                  `gorm:"default:false"`
    Permissions     datatypes.JSON        `gorm:"type:json"`
    Configuration   datatypes.JSON        `gorm:"type:json"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Organization    Organization          `gorm:"foreignKey:OrganizationID"`
}
```

## Usage Example

```bash
# Configure expenses module for an organization
curl -X POST "http://localhost:8080/api/organization-config/modules/configure" \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": 1,
    "modules": [
      {
        "module_name": "expenses",
        "is_enabled": true,
        "permissions": {
          "expense": ["read", "create", "update", "delete"],
          "category": ["read", "create"]
        },
        "configuration": {
          "max_expense_amount": 10000,
          "require_approval": true
        }
      }
    ]
  }'
```

## Best Practices Implemented

- **Clean Code**: Well-structured, readable code with proper naming conventions
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Security**: Authentication middleware for all endpoints
- **Separation of Concerns**: Clear separation between layers (handler, service, repository)
- **Validation**: Input validation using DTOs
- **Documentation**: Comprehensive comments and documentation

## Integration

This module integrates with:
- **RBAC System**: For dynamic permission management
- **Organization System**: For multi-tenant support
- **Authentication Middleware**: For security
- **Database**: Using GORM for data persistence

## Testing

Run tests with:
```bash
go test ./module/authentication/organization_config/...
```
