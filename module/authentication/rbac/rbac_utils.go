package rbac

import (
	"crypto/rand"
	"fmt"
)

// generateID genera un ID único usando crypto/rand
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
