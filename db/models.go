package db

import (
	"omsms/util/enums"

	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name      string
	Java      uint
	Backup    enums.BackupStrat
	HostNames []string `gorm:"type:json"`
}

func RegisterModels(db *gorm.DB) {
	db.AutoMigrate(&Server{})
}
