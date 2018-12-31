package model

import (
	"fmt"
	"time"
)

type Page struct {
	ID           uint `gorm:"primary_key"`
	URL          string
	Domain       string
	ResponseCode uint
	LoadTime     uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (p *Page) Save() error {
	if dbConnection.NewRecord(p) {
		dbConnection.Create(p)

		if dbConnection.NewRecord(p) {
			return fmt.Errorf("Failed to create Page record %v", p)
		}
	} else {
		dbConnection.Save(p)
	}
	return nil
}

func (p *Page) Delete() {
	dbConnection.Delete(p)
}
