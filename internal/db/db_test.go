package db

import (
	"bb-parse/internal/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var (
	record = &models.Record{
		Date:        "test-date",
		Description: "test-description",
		Category:    "test-category",
		Value:       10,
	}
)

type dbSuite struct {
	suite.Suite
	c *Client
}

func TestDB(t *testing.T) {
	c, err := NewDB()
	assert.NoError(t, err)

	s := new(dbSuite)
	s.c = c

	suite.Run(t, s)
}

func (dbs dbSuite) Test1Add() {
	err := dbs.c.AddRecords([]*models.Record{record})
	dbs.Require().NoError(err)
}

func (dbs dbSuite) Test2Get() {
	records, err := dbs.c.GetRecords()
	dbs.Require().NoError(err)
	dbs.Require().Equal(1, len(records))
	records[0].Log()
	record.ID = records[0].ID
}

func (dbs dbSuite) Test3Del() {
	err := dbs.c.DelRecord(record)
	dbs.Require().NoError(err)

	_, err = dbs.c.GetRecordByDesc(record.Description)
	dbs.Require().True(errors.Is(err, gorm.ErrRecordNotFound))
}
