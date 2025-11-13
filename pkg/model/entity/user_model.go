package entity

import (
	"learn/fiber/pkg/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       string     `gorm:"primary_key" json:"id"`
	Email    string     `gorm:"type:varchar(255); not null; unique" json:"email"`
	Username string     `gorm:"type:varchar(255); not null;" json:"username"`
	Role     enum.ERole `gorm:"type:varchar(255); not null;" json:"role"`
	Password string     `gorm:"type:varchar(255); not null;" json:"password"`
	Blogs    []Blog     `gorm:"foreignKey:UserId" json:"blogs,omitempty"`
}

func (user *User) BeforeCreate(db *gorm.DB) error {
	user.Id = "user-" + uuid.New().String()
	return nil
}

type UserRegisterRequest struct {
	Email           string     `validate:"required,email" json:"email"`
	Username        string     `validate:"required" json:"username"`
	Password        string     `validate:"required" json:"password"`
	ConfirmPassword string     `validate:"required" json:"confirmPassword"`
	Role            enum.ERole `validate:"required,oneof=admin user" json:"role"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,email" json:"email" example:"G2G5e@example.com"`
	Password string `validate:"required" json:"password" example:"P@ssw0rd!"`
}

type UserUpdateRequest struct {
	Email    string     `validate:"omitempty,email" json:"email"`
	Username string     `validate:"omitempty" json:"username"`
	Role     enum.ERole `validate:"omitempty,oneof=admin user" json:"role"`
}

type UserResponse struct {
	Id        string     `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Role      enum.ERole `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
