package models

import "fmt"

type Record struct {
	ID          int    `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
	Date        string // format YYYY-MM-DD
	Description string
	Value       float64
	Category    string
}

func (r Record) Log() {
	fmt.Printf(`	ID ----------- %v
	Date --------- %v
	Description -- %v
	Value -------- %v
	Category ----- %v`,
		r.ID,
		r.Date,
		r.Description,
		r.Value,
		r.Category,
	)
}
