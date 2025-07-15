package dto

// CreatePersonRequest represents the request to create a person
type CreatePersonRequest struct {
	Name      string `json:"name"       validate:"required,min=2,max=100"`
	Email     string `json:"email"      validate:"required,email"`
	TaxID     string `json:"tax_id"     validate:"required,min=11,max=14"`
	Phone     string `json:"phone"      validate:"omitempty,min=11,max=20"`
	BirthDate string `json:"birth_date" validate:"omitempty,datetime=2006-01-02"`
}

// UpdatePersonRequest represents the request to update a person
type UpdatePersonRequest struct {
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=11,max=20"`
}

// PersonResponse represents a person in API responses
type PersonResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	TaxID     string `json:"tax_id"`
	Phone     string `json:"phone,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
