package http_tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestHTTPIntegrationSuite ejecuta todos los tests de integración HTTP
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

// TestCompleteSchoolLifecycleHTTP ejecuta específicamente el test de colegio
func TestCompleteSchoolLifecycleHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

}

// TestEcommerceFlowHTTP ejecuta específicamente el test de e-commerce
func TestEcommerceFlowHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

}

// TestSecurityAndPermissionsHTTP ejecuta específicamente el test de seguridad
func TestSecurityAndPermissionsHTTP(t *testing.T) {
	testSuite := new(HTTPIntegrationTestSuite)
	testSuite.SetT(t)
	testSuite.SetupSuite()
	defer testSuite.TearDownSuite()

	testSuite.TestSecurityAndPermissionsHTTP()
}
