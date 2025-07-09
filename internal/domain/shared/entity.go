package shared

import "time"

// Entity is a base structure for all domain entities
type Entity struct {
	id        int64
	createdAt time.Time
	updateAt  time.Time
}

// NewEntity creates a base entity with a specific ID
func NewEntity(id int64) Entity {
	now := time.Now()
	return Entity{
		id:        id,
		createdAt: now,
		updateAt:  now,
	}
}

// ID returns the entity's ID
func (e Entity) ID() int64 { return e.id }

// CreatedAt returns the entity's creation timestamp
func (e Entity) CreatedAt() time.Time { return e.createdAt }

// UpdateAt returns the entity's last update timestamp
func (e Entity) UpdatedAt() time.Time { return e.updateAt }

// IsModifiedAfter checks if the entity was modified after the given time
func (e Entity) IsModifiedAfter(t time.Time) bool { return e.updateAt.After(t) }

// Age returns the duration since the entity was created
func (e Entity) Age() time.Duration { return time.Since(e.createdAt) }

// SetAsModified updates the modification timestamp
func (e *Entity) SetAsModified() { e.updateAt = time.Now() }
