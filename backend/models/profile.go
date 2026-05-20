package models

// ProfileInput is the request body for creating or updating a user profile.
type ProfileInput struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Image         string `json:"image_url"`
	EmailVerified bool   `json:"email_verified"`
}
