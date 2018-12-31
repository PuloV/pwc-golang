package model

import (
	"fmt"
	"time"
)

type KeyWord struct {
	ID        uint `gorm:"primary_key"`
	KeyWord   string
	PageID    uint
	Weight    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (kw *KeyWord) Save() error {
	if dbConnection.NewRecord(kw) {
		dbConnection.Create(kw)

		if dbConnection.NewRecord(kw) {
			return fmt.Errorf("Failed to create KeyWord record %v", kw)
		}
	} else {
		dbConnection.Save(kw)
	}
	return nil
}

func (kw *KeyWord) Delete() {
	dbConnection.Delete(kw)
}
