package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	PostID uint
	Name   string `json:"name"`
	Path   string `json:"path"`
}
