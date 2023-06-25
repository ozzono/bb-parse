package db

import (
	"bb-parse/internal/models"
	"log"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

type Client struct {
	DB *gorm.DB
}

func NewDB() (*Client, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open gorm.db")
	}
	err = db.AutoMigrate(&models.Record{})
	if err != nil {
		return nil, errors.Wrap(err, "db.AutoMigrate record")
	}
	return &Client{db}, nil
}

func (c *Client) AddRecords(r []*models.Record) error {
	tx := c.DB.Create(r)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "c.DB.Create records")
	}
	log.Printf("affected rows: %d", tx.RowsAffected)
	return nil
}

func (c *Client) GetRecords() ([]*models.Record, error) {
	out := []*models.Record{}
	tx := c.DB.Find(&out)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "c.DB.Find all records")
	}
	if len(out) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	log.Printf("records found: %d", len(out))
	return out, nil
}

func (c *Client) DelRecord(r *models.Record) error {
	tx := c.DB.Delete(&models.Record{}, r.ID)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "c.DB.Delete record")
	}
	log.Printf("affected rows: %d", tx.RowsAffected)
	return nil
}
