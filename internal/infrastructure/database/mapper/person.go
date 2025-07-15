package mapper

import (
	"todolist/internal/domain/person/entity"
	vo "todolist/internal/domain/person/valueobject"
	sharedvo "todolist/internal/domain/shared/valueobject"
	"todolist/internal/infrastructure/database/model"
)

// PersonMapper handles conversion between domain entity and database model
type PersonMapper struct{}

// NewPersonMapper creates a new PersonMapper
func NewPersonMapper() *PersonMapper {
	return &PersonMapper{}
}

// ToModel converts domain entity to database model
func (m *PersonMapper) ToModel(person *entity.Person) *model.Person {
	p := &model.Person{
		ID:        person.ID(),
		Name:      person.Name(),
		Phone:     person.Phone(),
		Email:     person.Email().Value(),
		TaxID:     person.TaxID().String(),
		CreatedAt: person.CreatedAt(),
		UpdatedAt: person.UpdatedAt(),
	}

	if !person.BirthDate().IsZero() {
		birthDate := person.BirthDate().Time()
		p.BirthDate = &birthDate
	}

	return p
}

// ToDomain converts database model to domain entity
func (m *PersonMapper) ToDomain(model *model.Person) (*entity.Person, error) {
	email, err := vo.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	taxID, err := vo.NewTaxID(model.TaxID)
	if err != nil {
		return nil, err
	}

	birthDate := sharedvo.Date{}
	if model.BirthDate != nil {
		birthDate, err = sharedvo.NewDate(model.BirthDate)
		if err != nil {
			return nil, err
		}
	}

	person, err := entity.NewPerson(
		model.ID,
		model.Name,
		model.Phone,
		taxID,
		email,
		&birthDate,
	)
	if err != nil {
		return nil, err
	}

	// Set timestamps from database
	person.Entity.SetCreatedAt(model.CreatedAt)
	person.Entity.SetUpdatedAt(model.UpdatedAt)

	return person, nil
}

// ToDomainList converts a list of models to domain entities
func (m *PersonMapper) ToDomainList(models []*model.Person) ([]*entity.Person, error) {
	people := make([]*entity.Person, 0, len(models))

	for _, model := range models {
		person, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, nil
}
