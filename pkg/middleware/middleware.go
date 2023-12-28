package middleware

import "gorm.io/gorm"

type Middleware struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Middleware {
	return &Middleware{
		DB: db,
	}
}
