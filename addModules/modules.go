// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/modules/modules.go
package modules

import (
	// Los módulos se registran mediante sus funciones init.

	_ "practicev2/module/auth"
	_ "practicev2/module/auth/rutascompuestas"
	_ "practicev2/module/imagenes"
	_ "practicev2/module/middleware"
	_ "practicev2/module/product"
	// Agregar nuevos módulos acá
)
