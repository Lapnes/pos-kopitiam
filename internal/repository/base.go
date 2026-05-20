package repository

import "gorm.io/gorm"

type BaseRepo struct {
	DB *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return BaseRepo{DB: db}
}