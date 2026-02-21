package migrations

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"web-porto-backend/internal/domain/models"
)

// Seed creates a default admin user if the users table is empty or if the admin email does not exist.
// Admin credentials are read from env ADMIN_EMAIL and ADMIN_PASSWORD for development convenience.
func Seed(db *gorm.DB) error {
	// Seed admin user
	if err := SeedAdminUser(db); err != nil {
		return err
	}

	// Seed categories
	if err := SeedCategories(db); err != nil {
		return err
	}

	// Seed analytics data
	if err := SeedAnalytics(db); err != nil {
		return err
	}

	return nil
}

// SeedAdminUser creates a default admin user if not exists
func SeedAdminUser(db *gorm.DB) error {
	// Ensure users table exists before seeding
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	// Check if admin user already exists (by email OR username)
	var count int64
	if err := db.Model(&models.User{}).Where("email = ? OR username = ?", adminEmail, "admin").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Admin user already exists, skipping seed")
		return nil
	}

	// Create admin user
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Username:     "admin",
		Email:        adminEmail,
		PasswordHash: string(hash),
		Role:         "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	log.Printf("Seeded default admin user: %s", adminEmail)
	return nil
}
