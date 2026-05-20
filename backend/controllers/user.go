package controllers

import (
	"fmt"
	"slices"
	"strings"

	"library-management-system/backend/initializers"
	"library-management-system/backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UpsertUserProfile creates or updates the authenticated user's profile (SMS onboarding or edits).
func UpsertUserProfile(c *fiber.Ctx) error {
	sub, err := subjectFromAccessToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	var body models.ProfileInput
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := validateBackupFieldsByMethod(sub, body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	u, err := upsertProfile(initializers.DB, sub, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"user": u})
}

// ProfileRequirements returns the missing profile fields based on passwordless method.
func ProfileRequirements(c *fiber.Ctx) error {
	sub, err := subjectFromAccessToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	var u models.User
	if err := initializers.DB.Where("auth0_id = ?", sub).First(&u).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	required := []string{"first_name", "last_name"}
	if strings.HasPrefix(sub, "sms|") {
		required = append(required, "email")
	} else {
		required = append(required, "phone_number")
	}

	missing := make([]string, 0, len(required))
	for _, field := range required {
		if isUserFieldMissing(u, field) {
			missing = append(missing, field)
		}
	}
	slices.Sort(missing)

	return c.JSON(fiber.Map{"required_fields": missing})
}

func upsertProfile(db *gorm.DB, auth0ID string, in models.ProfileInput) (*models.User, error) {
	auth0ID = strings.TrimSpace(auth0ID)
	if auth0ID == "" {
		return nil, fmt.Errorf("missing auth0_id")
	}
	firstName := strings.TrimSpace(in.FirstName)
	if firstName == "" {
		return nil, fmt.Errorf("first_name is required")
	}
	lastName := strings.TrimSpace(in.LastName)
	if lastName == "" {
		return nil, fmt.Errorf("last_name is required")
	}
	email := strings.TrimSpace(in.Email)
	phone := strings.TrimSpace(in.PhoneNumber)

	row := models.User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		PhoneNumber:   phone,
		EmailVerified: in.EmailVerified,
		Image:         strings.TrimSpace(in.Image),
		Auth0ID:       auth0ID,
	}

	var existing models.User
	err := db.Where("auth0_id = ?", auth0ID).First(&existing).Error
	// if err == gorm.ErrRecordNotFound {
	// 	if email == "" {
	// 		return nil, fmt.Errorf("email is required")
	// 	}
	// 	if err := db.Create(&row).Error; err != nil {
	// 		return nil, err
	// 	}
	// 	return &row, nil
	// }
	if err != nil {
		return nil, err
	}

	existing.FirstName = row.FirstName
	existing.LastName = row.LastName
	if row.Email != "" {
		existing.Email = row.Email
	}
	if row.PhoneNumber != "" {
		existing.PhoneNumber = row.PhoneNumber
	}
	existing.EmailVerified = row.EmailVerified
	if row.Image != "" {
		existing.Image = row.Image
	}

	if err := db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func validateBackupFieldsByMethod(sub string, in models.ProfileInput) error {
	sub = strings.TrimSpace(sub)
	email := strings.TrimSpace(in.Email)
	phone := strings.TrimSpace(in.PhoneNumber)

	switch {
	case strings.HasPrefix(sub, "email|"):
		if phone == "" {
			return fmt.Errorf("phone_number is required as backup for email+otp")
		}
	case strings.HasPrefix(sub, "sms|"):
		if email == "" {
			return fmt.Errorf("email is required as backup for mobile+otp")
		}
	}
	return nil
}

func isUserFieldMissing(u models.User, field string) bool {
	switch field {
	case "first_name":
		return strings.TrimSpace(u.FirstName) == ""
	case "last_name":
		return strings.TrimSpace(u.LastName) == ""
	case "email":
		return strings.TrimSpace(u.Email) == ""
	case "phone_number":
		return strings.TrimSpace(u.PhoneNumber) == ""
	default:
		return false
	}
}
