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
	Name  *string `json:"name,omitempty"  validate:"omitempty,min=2,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=11,max=20"`
}
