package models

type User struct {
	Id        uint
	Username  string `validate:"required,min=3,max=32"`
	FirstName string `validate:"required,min=3,max=32"`
	LastName  string `validate:"required,min=3,max=32"`
	Email     string `validate:"required,email,max=32" ,gorm:"unique"`
	Password  []byte
}
