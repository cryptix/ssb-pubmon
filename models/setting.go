package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
)

type Setting struct {
	gorm.Model
	l10n.Locale
}
