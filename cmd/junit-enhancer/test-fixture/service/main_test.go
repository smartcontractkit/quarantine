package main

import (
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	service := NewUserService()
	
	user := service.CreateUser("John Doe", "john@example.com")
	
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %s", user.Email)
	}
}

func TestUserService_GetUser(t *testing.T) {
	service := NewUserService()
	
	// Test getting non-existent user
	t.Run("non_existent_user", func(t *testing.T) {
		_, exists := service.GetUser(999)
		if exists {
			t.Error("Expected user to not exist")
		}
	})
	
	// Test getting existing user
	t.Run("existing_user", func(t *testing.T) {
		created := service.CreateUser("Jane Doe", "jane@example.com")
		
		user, exists := service.GetUser(created.ID)
		if !exists {
			t.Error("Expected user to exist")
		}
		if user.Name != "Jane Doe" {
			t.Errorf("Expected name 'Jane Doe', got %s", user.Name)
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	service := NewUserService()
	
	// Test deleting non-existent user
	t.Run("non_existent_user", func(t *testing.T) {
		deleted := service.DeleteUser(999)
		if deleted {
			t.Error("Expected deletion to fail for non-existent user")
		}
	})
	
	// Test deleting existing user
	t.Run("existing_user", func(t *testing.T) {
		user := service.CreateUser("Bob Smith", "bob@example.com")
		
		deleted := service.DeleteUser(user.ID)
		if !deleted {
			t.Error("Expected deletion to succeed")
		}
		
		// Verify user is gone
		_, exists := service.GetUser(user.ID)
		if exists {
			t.Error("Expected user to be deleted")
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	service := NewUserService()
	
	// Test empty list
	t.Run("empty_list", func(t *testing.T) {
		users := service.ListUsers()
		if len(users) != 0 {
			t.Errorf("Expected 0 users, got %d", len(users))
		}
	})
	
	// Test list with users
	t.Run("with_users", func(t *testing.T) {
		service.CreateUser("User 1", "user1@example.com")
		service.CreateUser("User 2", "user2@example.com")
		service.CreateUser("User 3", "user3@example.com")
		
		users := service.ListUsers()
		if len(users) != 3 {
			t.Errorf("Expected 3 users, got %d", len(users))
		}
	})
}

func FuzzCreateUser(f *testing.F) {
	f.Add("John Doe", "john@example.com")
	f.Add("", "")
	f.Add("Alice", "alice@test.com")
	
	f.Fuzz(func(t *testing.T, name, email string) {
		service := NewUserService()
		user := service.CreateUser(name, email)
		
		if user.Name != name {
			t.Errorf("Expected name %q, got %q", name, user.Name)
		}
		if user.Email != email {
			t.Errorf("Expected email %q, got %q", email, user.Email)
		}
		if user.ID != 1 {
			t.Errorf("Expected first user to have ID 1, got %d", user.ID)
		}
	})
}
