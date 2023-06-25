package db

import (
	"bb-parse/internal/models"
	"bb-parse/utils"
	"log"
	"strconv"
	"strings"

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

func (c *Client) GetRecordByDesc(desc string) (*models.Record, error) {
	out := &models.Record{}
	tx := c.DB.First(out, "description = ?", desc)
	if tx.Error != nil {
		return nil, errors.Wrapf(tx.Error, "c.DB.First description = %s", desc)
	}
	return out, nil
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

func (c *Client) ParseAndCompare(input string) (*models.Record, error) {
	out := &models.Record{}
	splittedDate := strings.Split(input[:10], ".")
	out.Date = strings.Join([]string{splittedDate[2], splittedDate[1], splittedDate[0]}, "-")
	out.Description = strings.ToLower(utils.TrimSpaces(input[10:47]))
	v, _ := strconv.ParseFloat(strings.ReplaceAll(utils.TrimSpaces(input[49:69]), ",", "."), 64)
	out.Value = v

	dbRecord, err := c.GetRecordByDesc(out.Description)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return out, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "c.GetRecordByDesc %s", out.Description)
	}
	out.Category = dbRecord.Category
	return out, nil
}
