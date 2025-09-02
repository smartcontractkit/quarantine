package models

import (
	"errors"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	InStock     bool      `json:"in_stock"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductRepository struct {
	products map[int]*Product
	nextID   int
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[int]*Product),
		nextID:   1,
	}
}

func (pr *ProductRepository) Create(name, category string, price float64) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}
	
	product := &Product{
		ID:        pr.nextID,
		Name:      name,
		Price:     price,
		Category:  category,
		InStock:   true,
		CreatedAt: time.Now(),
	}
	
	pr.products[product.ID] = product
	pr.nextID++
	return product, nil
}

func (pr *ProductRepository) GetByID(id int) (*Product, error) {
	product, exists := pr.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (pr *ProductRepository) GetByCategory(category string) []*Product {
	var products []*Product
	for _, product := range pr.products {
		if product.Category == category {
			products = append(products, product)
		}
	}
	return products
}

func (pr *ProductRepository) UpdateStock(id int, inStock bool) error {
	product, exists := pr.products[id]
	if !exists {
		return errors.New("product not found")
	}
	product.InStock = inStock
	return nil
}

func (pr *ProductRepository) Delete(id int) error {
	if _, exists := pr.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(pr.products, id)
	return nil
}
