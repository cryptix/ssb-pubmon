package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Pub struct {
	gorm.Model
	Host, Key   string
	Port        int
	LastSuccess time.Time
	Failures    int
}
