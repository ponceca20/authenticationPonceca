package product

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"practicev2/database"
	"practicev2/registry"
	"practicev2/registry/cache"

	"unicode"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Error si falta la URL de Redis en las variables de entorno
var (
	ErrMissingRedisURL = errors.New("missing Redis URL in environment variables")
)

// Esquema de los campos que se indexan en RedisSearch
var productSchema = map[string]string{
	"producto_id":               "NUMERIC",
	"empresa_id":                "NUMERIC",
	"producto_nombre":           "TEXT",
	"producto_nombre_fold":      "TEXT",
	"producto_descripcion":      "TEXT",
	"producto_descripcion_fold": "TEXT",
	"producto_precio":           "NUMERIC",
	"producto_codigo":           "TEXT",
	"producto_codigo_barras":    "TEXT",
	"tpv_visible":               "TAG",
	"almacen_visible":           "TAG",
	"es_inactivo":               "TAG",
	"tipo_producto":             "TEXT",
	"subcategoria":              "TEXT",
	"categoria":                 "TEXT",
	"clase":                     "TEXT",
	"imagen_url":                "TEXT",
}

// removeDiacritics normalizes a string by removing diacritical marks
func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// Estructura que representa una fila de producto desde la base de datos
type productoRow struct {
	ProductoID           uint64  `gorm:"column:producto_id"`
	EmpresaID            uint64  `gorm:"column:empresa_id"`
	ProductoNombre       string  `gorm:"column:producto_nombre"`
	ProductoDescripcion  string  `gorm:"column:producto_descripcion"`
	ProductoPrecio       float64 `gorm:"column:producto_precio"`
	ProductoCodigo       string  `gorm:"column:producto_codigo"`
	ProductoCodigoBarras string  `gorm:"column:producto_codigo_barras"`
	TpvVisible           bool    `gorm:"column:tpv_visible"`
	AlmacenVisible       bool    `gorm:"column:almacen_visible"`
	EsInactivo           bool    `gorm:"column:es_inactivo"`
	TipoProducto         string  `gorm:"column:tipo_producto"`
	Subcategoria         string  `gorm:"column:subcategoria"`
	Categoria            string  `gorm:"column:categoria"`
	Clase                string  `gorm:"column:clase"`
	ImagenURL            string  `gorm:"column:imagen_url"`
}

// Convierte la estructura productoRow a un mapa para indexar en RedisSearch
func (r productoRow) toFields() map[string]interface{} {
	return map[string]interface{}{
		"producto_id":               r.ProductoID,
		"empresa_id":                r.EmpresaID,
		"producto_nombre":           r.ProductoNombre,
		"producto_nombre_fold":      strings.ToLower(removeDiacritics(r.ProductoNombre)),
		"producto_descripcion":      r.ProductoDescripcion,
		"producto_descripcion_fold": strings.ToLower(removeDiacritics(r.ProductoDescripcion)),
		"producto_precio":           r.ProductoPrecio,
		"producto_codigo":           r.ProductoCodigo,
		"producto_codigo_barras":    r.ProductoCodigoBarras,
		"tpv_visible":               strings.ToLower(strconv.FormatBool(r.TpvVisible)),
		"almacen_visible":           strings.ToLower(strconv.FormatBool(r.AlmacenVisible)),
		"es_inactivo":               strings.ToLower(strconv.FormatBool(r.EsInactivo)),
		"tipo_producto":             r.TipoProducto,
		"subcategoria":              r.Subcategoria,
		"categoria":                 r.Categoria,
		"clase":                     r.Clase,
		"imagen_url":                r.ImagenURL,
	}
}

// Servicio para manejar la búsqueda e indexación de productos
type ProductoSearchService struct {
	ctx    context.Context
	client cache.RediSearchClient
}

// Crea una nueva instancia del servicio de búsqueda de productos
func NewProductoSearchService(ctx context.Context, recreateIndex ...bool) (*ProductoSearchService, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, ErrMissingRedisURL
	}
	c, err := cache.NewRediSearchClient(url)
	if err != nil {
		return nil, err
	}
	// Si se indica, recrea el índice en RedisSearch
	if len(recreateIndex) > 0 && recreateIndex[0] {
		if err := c.CreateIndex(ctx, cache.ProductIndexName, productSchema); err != nil {
			return nil, err
		}
	}
	return &ProductoSearchService{ctx: ctx, client: c}, nil
}

const (
	defaultTimeout = 5 * time.Second // Tiempo máximo para operaciones
	batchSize      = 2000            // Tamaño de lote para indexación masiva
)

// Consulta base para obtener productos y sus relaciones
const baseProductQuery = `
SELECT
  p.id AS producto_id,
  p.empresa_id,
  p.nombre AS producto_nombre,
  p.descripcion AS producto_descripcion,
  p.precio AS producto_precio,
  p.codigo AS producto_codigo,
  p.codigo_barras AS producto_codigo_barras,
  p.tpv_visible,
  p.almacen_visible,
  p.es_inactivo,
  tp.nombre AS tipo_producto,
  sc.nombre AS subcategoria,
  cat.nombre AS categoria,
  c.nombre AS clase,
  (
    SELECT pm.url
    FROM producto_media pm
    WHERE pm.producto_id = p.id
    ORDER BY pm.es_principal DESC, pm.orden ASC
    LIMIT 1
  ) AS imagen_url
FROM producto p
LEFT JOIN tipo_producto tp ON p.tipo_producto_id = tp.id
LEFT JOIN sub_categoria sc ON p.sub_categoria_id = sc.id
LEFT JOIN categoria cat ON sc.categoria_id = cat.id
LEFT JOIN clase c ON p.clase_id = c.id
`

// Consulta para obtener productos por lotes usando un cursor
const indexSQLCursor = baseProductQuery + `
WHERE p.id > ?
ORDER BY p.id
LIMIT ?;
`

// Consulta para obtener un solo producto por ID
const singleIndexSQL = baseProductQuery + `
WHERE p.id = ?
LIMIT 1;
`

// Agrega un documento a RedisSearch con timeout
func (s *ProductoSearchService) addWithTimeout(id string, fields map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(s.ctx, defaultTimeout)
	defer cancel()
	return s.client.AddDocument(ctx, cache.ProductIndexName, id, fields)
}

// Elimina un documento de RedisSearch con timeout
func (s *ProductoSearchService) deleteWithTimeout(id string) error {
	ctx, cancel := context.WithTimeout(s.ctx, defaultTimeout)
	defer cancel()
	return s.client.DeleteDocument(ctx, cache.ProductIndexName, id)
}

// Indexa todos los productos de la base de datos en RedisSearch
func (s *ProductoSearchService) IndexAll() error {
	var lastID uint64
	var total, failed int
	for {
		var rows []productoRow
		if err := database.DBconn.Raw(indexSQLCursor, lastID, batchSize).Scan(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			break
		}
		for _, r := range rows {
			id := strconv.FormatUint(r.ProductoID, 10)
			if err := s.addWithTimeout(id, r.toFields()); err != nil {
				log.Printf("IndexAll: failed to index product %s: %v", id, err)
				failed++
				continue
			}
			if r.ProductoID > lastID {
				lastID = r.ProductoID
			}
			total++
		}
		if len(rows) < batchSize {
			break
		}
	}
	log.Printf("IndexAll: indexed %d products, failed %d", total, failed)
	if failed > 0 {
		return errors.New("IndexAll: some products failed to index")
	}
	return nil
}

// Indexa un solo producto por su ID
func (s *ProductoSearchService) IndexProductByID(id uint64) error {
	var r productoRow
	if err := database.DBconn.Raw(singleIndexSQL, id).Scan(&r).Error; err != nil {
		return err
	}
	docID := strconv.FormatUint(id, 10)
	return s.addWithTimeout(docID, r.toFields())
}

// Elimina un producto del índice por su ID
func (s *ProductoSearchService) DeleteProductByID(id uint64) error {
	docID := strconv.FormatUint(id, 10)
	return s.deleteWithTimeout(docID)
}

// Realiza una búsqueda en RedisSearch
func (s *ProductoSearchService) Search(q string, offset, limit int) ([]map[string]string, int64, error) {
	// normalize query to remove diacritics for accent-insensitive search
	normQ := strings.ToLower(removeDiacritics(q))
	ctx, cancel := context.WithTimeout(s.ctx, defaultTimeout)
	defer cancel()
	return s.client.Search(ctx, cache.ProductIndexName, normQ, offset, limit)
}

// Actualiza (reindexa) productos por clasificación (tipo, subcategoría, etc.)
func (s *ProductoSearchService) UpdateProductsByClassification(classificationField, classificationValue string) error {
	var ids []uint64
	var tableAlias string
	switch classificationField {
	case "tipo_producto":
		tableAlias = "tp"
	case "subcategoria":
		tableAlias = "sc"
	case "categoria":
		tableAlias = "cat"
	case "clase":
		tableAlias = "c"
	default:
		return errors.New("invalid classification field")
	}
	if strings.TrimSpace(classificationValue) == "" {
		return errors.New("classification value cannot be empty")
	}
	queryIDs := `
SELECT p.id AS producto_id
FROM producto p
LEFT JOIN tipo_producto tp ON p.tipo_producto_id = tp.id
LEFT JOIN sub_categoria sc ON p.sub_categoria_id = sc.id
LEFT JOIN categoria cat ON sc.categoria_id = cat.id
LEFT JOIN clase c ON p.clase_id = c.id
WHERE ` + tableAlias + `.nombre = ?
ORDER BY p.id;
`
	rows, err := database.DBconn.Raw(queryIDs, classificationValue).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	var failed int
	for _, id := range ids {
		if err := s.IndexProductByID(id); err != nil {
			log.Printf("UpdateProductsByClassification: failed to reindex product %d: %v", id, err)
			failed++
		}
	}
	if failed > 0 {
		return errors.New("UpdateProductsByClassification: some products failed to reindex")
	}
	return nil
}

// RecreateIndex dropea y crea el índice con el esquema actual de producto
func (s *ProductoSearchService) RecreateIndex() error {
	return s.client.CreateIndex(s.ctx, cache.ProductIndexName, productSchema)
}

// Handler HTTP para exponer endpoints de búsqueda de productos
type ProductoSearchHandler struct {
	svc *ProductoSearchService
}

// Crea un nuevo handler HTTP para búsqueda de productos
func NewProductoSearchHandler(svc *ProductoSearchService) *ProductoSearchHandler {
	return &ProductoSearchHandler{svc: svc}
}

// Endpoint: /api/v1/product/search
// Realiza una búsqueda de productos usando los parámetros de la URL
func (h *ProductoSearchHandler) Search(c *fiber.Ctx) error {
	q := c.Query("q", "*")                         // Consulta de búsqueda, por defecto "*"
	off, _ := strconv.Atoi(c.Query("offset", "0")) // Offset para paginación
	lim, _ := strconv.Atoi(c.Query("limit", "10")) // Límite de resultados
	data, total, err := h.svc.Search(q, off, lim)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success", "query": q,
		"offset": off, "limit": lim,
		"total": total, "data": data,
	})
}

// Endpoint: /api/v1/product/health
// Verifica la salud del servicio de búsqueda y la base de datos
func (h *ProductoSearchHandler) Health(c *fiber.Ctx) error {
	if _, _, err := h.svc.Search("*", 0, 1); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "error", "message": "search unavailable",
		})
	}
	ctx, cancel := context.WithTimeout(h.svc.ctx, defaultTimeout)
	defer cancel()
	if err := database.DBconn.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "error", "message": "db unavailable",
		})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

// Endpoint: /api/v1/product/reindex
// Reindexa (recrea índice y vuelve a indexar) todos los productos en RedisSearch
func (h *ProductoSearchHandler) Reindex(c *fiber.Ctx) error {
	// 1) recrear índice según esquema
	if err := h.svc.RecreateIndex(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al recrear índice: " + err.Error(),
		})
	}
	// 2) indexar todos los productos
	if err := h.svc.IndexAll(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error reindexando productos: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{"status": "ok", "message": "Índice recreado y productos reindexados"})
}

// Registra las rutas HTTP para el módulo de búsqueda de productos
func RegisterRoutesProductoSearch(app *fiber.App) {
	ctx := context.Background()
	recreate := os.Getenv("RECREATE_PRODUCT_INDEX") == "1"
	svc, err := NewProductoSearchService(ctx, recreate)
	if err != nil {
		log.Printf("WARNING: Search service not available: %v", err)
		return
	}
	// Indexa todos los productos en un goroutine al iniciar
	go func() {
		if err := svc.IndexAll(); err != nil {
			log.Printf("WARNING: Error indexing products: %v", err)
		}
	}()
	h := NewProductoSearchHandler(svc)
	api := app.Group("/api/v1/product")
	api.Get("/search", h.Search)
	api.Get("/health", h.Health)
	api.Post("/reindex", h.Reindex)
}

// Registra el módulo en el sistema de módulos de la aplicación
func init() {
	registry.RegisterModule(RegisterRoutesProductoSearch)
}
