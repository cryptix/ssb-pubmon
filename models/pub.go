package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Pub struct {
	gorm.Model
	Key         string `gorm:"unique_index"`
	Addresses   []Address
	LastSuccess time.Time
}

// Composite primary keys can't be created on SQLite tables
// https://github.com/jinzhu/gorm/issues/1037
type Address struct {
	Host  string `gorm:"primary_key"`
	Port  int    `gorm:"primary_key"`
	PubID uint
}
