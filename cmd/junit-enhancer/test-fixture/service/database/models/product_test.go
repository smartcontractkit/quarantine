package models

import (
	"testing"
	"time"
)

func TestProductRepository_Create(t *testing.T) {
	repo := NewProductRepository()
	
	t.Run("valid_product", func(t *testing.T) {
		product, err := repo.Create("Laptop", "Electronics", 999.99)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if product.ID != 1 {
			t.Errorf("Expected ID 1, got %d", product.ID)
		}
		if product.Name != "Laptop" {
			t.Errorf("Expected name 'Laptop', got %s", product.Name)
		}
		if product.Price != 999.99 {
			t.Errorf("Expected price 999.99, got %f", product.Price)
		}
		if !product.InStock {
			t.Error("Expected product to be in stock by default")
		}
	})
	
	t.Run("empty_name", func(t *testing.T) {
		_, err := repo.Create("", "Electronics", 100.0)
		if err == nil {
			t.Error("Expected error for empty product name")
		}
	})
	
	t.Run("negative_price", func(t *testing.T) {
		_, err := repo.Create("Test Product", "Test", -10.0)
		if err == nil {
			t.Error("Expected error for negative price")
		}
	})
}

func TestProductRepository_GetByID(t *testing.T) {
	repo := NewProductRepository()
	
	t.Run("existing_product", func(t *testing.T) {
		created, err := repo.Create("Mouse", "Electronics", 29.99)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
		
		retrieved, err := repo.GetByID(created.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve product: %v", err)
		}
		
		if retrieved.Name != "Mouse" {
			t.Errorf("Expected name 'Mouse', got %s", retrieved.Name)
		}
	})
	
	t.Run("non_existent_product", func(t *testing.T) {
		_, err := repo.GetByID(999)
		if err == nil {
			t.Error("Expected error for non-existent product")
		}
	})
}

func TestProductRepository_GetByCategory(t *testing.T) {
	repo := NewProductRepository()
	
	// Create products in different categories
	repo.Create("Laptop", "Electronics", 999.99)
	repo.Create("Mouse", "Electronics", 29.99)
	repo.Create("Desk", "Furniture", 199.99)
	repo.Create("Chair", "Furniture", 149.99)
	repo.Create("Book", "Education", 19.99)
	
	t.Run("electronics_category", func(t *testing.T) {
		products := repo.GetByCategory("Electronics")
		if len(products) != 2 {
			t.Errorf("Expected 2 electronics products, got %d", len(products))
		}
	})
	
	t.Run("furniture_category", func(t *testing.T) {
		products := repo.GetByCategory("Furniture")
		if len(products) != 2 {
			t.Errorf("Expected 2 furniture products, got %d", len(products))
		}
	})
	
	t.Run("non_existent_category", func(t *testing.T) {
		products := repo.GetByCategory("NonExistent")
		if len(products) != 0 {
			t.Errorf("Expected 0 products for non-existent category, got %d", len(products))
		}
	})
}

func TestProductRepository_UpdateStock(t *testing.T) {
	repo := NewProductRepository()
	
	product, err := repo.Create("Test Product", "Test", 50.0)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	t.Run("update_existing_product", func(t *testing.T) {
		err := repo.UpdateStock(product.ID, false)
		if err != nil {
			t.Fatalf("Failed to update stock: %v", err)
		}
		
		updated, err := repo.GetByID(product.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve updated product: %v", err)
		}
		
		if updated.InStock {
			t.Error("Expected product to be out of stock")
		}
	})
	
	t.Run("update_non_existent_product", func(t *testing.T) {
		err := repo.UpdateStock(999, true)
		if err == nil {
			t.Error("Expected error for non-existent product")
		}
	})
}

func TestProductRepository_Delete(t *testing.T) {
	repo := NewProductRepository()
	
	product, err := repo.Create("To Delete", "Test", 25.0)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	
	t.Run("delete_existing_product", func(t *testing.T) {
		err := repo.Delete(product.ID)
		if err != nil {
			t.Fatalf("Failed to delete product: %v", err)
		}
		
		_, err = repo.GetByID(product.ID)
		if err == nil {
			t.Error("Expected error when retrieving deleted product")
		}
	})
	
	t.Run("delete_non_existent_product", func(t *testing.T) {
		err := repo.Delete(999)
		if err == nil {
			t.Error("Expected error when deleting non-existent product")
		}
	})
}

func FuzzProductRepository_Create(f *testing.F) {
	f.Add("Laptop", "Electronics", 999.99)
	f.Add("", "Test", 100.0)
	f.Add("Test Product", "Category", -50.0)
	f.Add("Valid Product", "", 25.0)
	
	f.Fuzz(func(t *testing.T, name, category string, price float64) {
		repo := NewProductRepository()
		
		product, err := repo.Create(name, category, price)
		
		if name == "" || price < 0 {
			// Should return error
			if err == nil {
				t.Errorf("Expected error for invalid input: name=%q, price=%f", name, price)
			}
			return
		}
		
		// Should succeed
		if err != nil {
			t.Errorf("Unexpected error for valid input: %v", err)
			return
		}
		
		if product.Name != name {
			t.Errorf("Expected name %q, got %q", name, product.Name)
		}
		if product.Category != category {
			t.Errorf("Expected category %q, got %q", category, product.Category)
		}
		if product.Price != price {
			t.Errorf("Expected price %f, got %f", price, product.Price)
		}
		if !product.InStock {
			t.Error("Expected product to be in stock by default")
		}
		if product.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}
		if time.Since(product.CreatedAt) > time.Second {
			t.Error("CreatedAt should be recent")
		}
	})
}
