package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	PartnerID int64     `json:"partner_id"` // Referencia Partner (tenant)
	Name      string    `json:"name"`
	Mobile    string    `json:"mobile,omitempty"`
	Password  string    `json:"password"`
	Type      int       `json:"type"` // 1 = master, 2 = administrator, 3 = colaborator
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	List(ctx context.Context, partnerID int64, limit, offset int64) ([]*User, error)
	ListDeleted(ctx context.Context, partnerID int64, limit, offset int64) ([]*User, error)
	GetByID(ctx context.Context, partnerID, id int64) (*User, error)
	GetByMobile(ctx context.Context, partnerID int64, mobile string) (*User, error)
	GetByMobileGlobal(ctx context.Context, mobile string) (*User, error)
	Update(ctx context.Context, partnerID int64, user *User) error
	Delete(ctx context.Context, partnerID, id int64) error
	ToggleActive(ctx context.Context, partnerID, id int64, active bool) error
}

type UserCreatedHook interface {
	OnUserCreated(ctx context.Context, user *User, plainPassword string)
}

func (u *User) Validate() error {
	if u.PartnerID <= 0 {
		return errors.New("user partner_id is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Mobile == "" {
		return errors.New("mobile is required")
	}
	if u.Type < 1 || u.Type > 3 {
		return errors.New("type must be 1 (master) or 2 (administrator) or 3 (colaborator)")
	}
	return nil
}

func (u *User) Activate() {
	u.Active = true
}
