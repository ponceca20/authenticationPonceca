package http_tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestHTTPIntegrationSuite ejecuta todos los tests de integración HTTP usando testify suite
func TestHTTPIntegrationSuite(t *testing.T) {
	suite.Run(t, new(HTTPIntegrationTestSuite))
}

// TestCompleteCompanyLifecycleHTTP ejecuta específicamente el test de empresa
func TestCompleteCompanyLifecycleHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

	testSuite.TestCompleteCompanyLifecycleHTTP()
}

// TestCompleteSchoolLifecycleHTTP ejecuta específicamente el test de colegio (placeholder)
func TestCompleteSchoolLifecycleHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

	// TODO: Implementar test de ciclo de vida de colegio
	// testSuite.TestCompleteSchoolLifecycleHTTP()
	t.Skip("School lifecycle test not implemented yet")
}

// TestEcommerceFlowHTTP ejecuta específicamente el test de e-commerce (placeholder)
func TestEcommerceFlowHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

	// TODO: Implementar test de flujo de e-commerce
	// testSuite.TestEcommerceFlowHTTP()
	t.Skip("E-commerce flow test not implemented yet")
}

// TestSecurityAndPermissionsHTTP ejecuta específicamente el test de seguridad (placeholder)
func TestSecurityAndPermissionsHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

	// TODO: Implementar test de seguridad y permisos
	// testSuite.TestSecurityAndPermissionsHTTP()
	t.Skip("Security and permissions test not implemented yet")
}
