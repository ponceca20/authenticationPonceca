package categorias

type CategoryService struct {
	Repo CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		Repo: NewCategoryRepository(),
	}
}
func (s *CategoryService) GetAllCategories() ([]Categoria, error) {
	return s.Repo.GetAllCategories()
}
func (s *CategoryService) GetCategory(id int) (Categoria, error) {
	return s.Repo.GetCategory(id)
}
func (s *CategoryService) NewCategory(category *Categoria) error {
	return s.Repo.CreateCategory(category)
}
func (s *CategoryService) UpdateCategory(id int, updatedCategory *Categoria) (Categoria, error) {
	return s.Repo.UpdateCategory(id, updatedCategory)
}
func (s *CategoryService) DeleteCategory(id int) error {
	return s.Repo.DeleteCategory(id)
}
func (s *CategoryService) SeedCategories(categories []Categoria) error {
	return s.Repo.SeedCategories(categories)
}
