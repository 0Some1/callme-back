package models

type User struct {
	Id        uint   `json:"id"`
	Username  string `validate:"required,min=3,max=32" json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique" validate:"required,email,max=32" json:"email"`
	Password  []byte `validate:"required" json:"-"`
}
