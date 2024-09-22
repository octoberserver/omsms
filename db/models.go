package db

import (
	"omsms/util/enums"

	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name   string
	Java   uint
	Backup enums.BackupStrat
}

func RegisterModels(db *gorm.DB) {
	db.AutoMigrate(&Server{})
}
