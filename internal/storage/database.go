package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(databaseURL string) *Database {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	return &Database{db: db}
}

// Migrate auto-migrates the tables for Audit and PageResult
func (d *Database) Migrate() error {
	return d.db.AutoMigrate(&Audit{}, &PageResult{})
}

// CreateAudit inserts a new audit record
func (d *Database) CreateAudit(audit *Audit) error {
	return d.db.Create(audit).Error
}

func (d *Database) GetAudit(id string) (*Audit, error) {
	var audit Audit
	err := d.db.Preload("PageResults").First(&audit, "id = ?", id).Error
	return &audit, err
}

// UpdateAudit saves changes to an existing audit
func (d *Database) UpdateAudit(audit *Audit) error {
	return d.db.Save(audit).Error
}

// CreatePageResult inserts a page result for a given audit
func (d *Database) CreatePageResult(result *PageResult) error {
	return d.db.Create(result).Error
}

func (d *Database) DB() *gorm.DB {
	return d.db
}
