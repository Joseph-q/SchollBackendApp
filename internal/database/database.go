package database

import (
	"log"
	"sync"

	"github.com/juseph-q/SchoolPr/internal/config"
	"github.com/juseph-q/SchoolPr/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDataBase(conf *config.Config) (*gorm.DB, error) {
	var (
		db   *gorm.DB
		once sync.Once
	)

	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(conf.Database.Url), &gorm.Config{})

		db.AutoMigrate(models.Courses{}, models.Students{}, models.Assistance{}, models.AssistanceHistorial{})
		if err != nil {
			log.Fatalf("Error al conectar a la base de datos: %v", err)
		}
	})

	return db, nil
}
