package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID   uint
	Follower User `gorm:"foreignkey:UserID"`
	state
}

type state int

const (
	declined state = iota
	accepted
	pending
)
