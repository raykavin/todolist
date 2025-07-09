package entity

import (
	"errors"
	"todolist/internal/domain/person/valueobject"
	"todolist/internal/domain/shared"
	sharedvo "todolist/internal/domain/shared/valueobject"
)

var (
	ErrInvalidPersonName = errors.New("person name cannot be empty")
	ErrInvalidPhone      = errors.New("invalid phone number")
)

// Person represents a person in the system
type Person struct {
	shared.Entity
	name      string
	phone     string
	taxID     valueobject.TaxID
	email     valueobject.Email
	birthDate sharedvo.Date
}

// NewPerson creates a new Person entity
func NewPerson(
	id uint64,
	name string,
	phone string,
	taxID valueobject.TaxID,
	email valueobject.Email,
	birthDate sharedvo.Date,
) (*Person, error) {
	p := &Person{
		Entity:    shared.NewEntity(id),
		name:      name,
		taxID:     taxID,
		email:     email,
		phone:     phone,
		birthDate: birthDate,
	}

	if err := p.validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p Person) validate() error {
	if p.name == "" {
		return ErrInvalidPersonName
	}

	if p.phone == "" {
		return ErrInvalidPhone
	}

	return nil
}

// Getter methods

// Name returns the name of the person
func (p Person) Name() string { return p.name }

func (p Person) TaxID() valueobject.TaxID { return p.taxID }

// Email returns the email of the person
func (p Person) Email() valueobject.Email { return p.email }

// Phone returns the phone number of the person
func (p Person) Phone() string { return p.phone }

// BirthDate returns the birth date of the person
func (p Person) BirthDate() sharedvo.Date { return p.birthDate }

// Update methods

// UpdateName updates the name of the person
func (p *Person) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidPersonName
	}
	p.name = name
	p.SetAsModified()
	return nil
}

// UpdateEmail updates the email of the person
func (p *Person) UpdateEmail(email valueobject.Email) {
	p.email = email
	p.SetAsModified()
}

// UpdatePhone updates the phone number of the person
func (p *Person) UpdatePhone(phone string) error {
	if phone == "" {
		return ErrInvalidPhone
	}
	p.phone = phone
	p.SetAsModified()
	return nil
}
