package main

type ProductDAO struct {
	id          int
	title       string
	description string
	amount      int
	is_gift     bool
}

type InMemoryRepository struct {
	Products []ProductDAO
}

func NewInMemoryRepository(jsonFilePath string) InMemoryRepository {
	// open json file
	return InMemoryRepository{}
}

func (m InMemoryRepository) FindById(id int) (ProductDAO, error) {
	return ProductDAO{}, nil
}

type Repository interface {
	Find(id int) (ProductDAO, error)
}
