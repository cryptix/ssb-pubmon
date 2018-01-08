package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/qor/audited"
	"github.com/qor/l10n"
	"github.com/qor/media"
	"github.com/qor/sorting"
	"github.com/qor/validations"

	"github.com/cryptix/ssb-pubmon/config"
)

var (
	db *gorm.DB
)

func GetBase() *gorm.DB {
	if db == nil {
		panic("sbm-db: db is still nil")
	}
	return db
}

func Open() error {
	if db != nil {
		return nil
	}
	var err error

	if config.Config == nil {
		return errors.New("sbm: config not loaded")
	}

	if config.Config.DB.Adapter == "postgres" {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	} else if config.Config.DB.Adapter == "sqlite" {
		db, err = gorm.Open("sqlite3", config.Config.DB.Name)
	} else {
		return errors.New("not supported database adapter")
	}
	if err != nil {
		return errors.Wrap(err, "db.Open failed")
	}

	if os.Getenv("DEBUG") != "" {
		db.LogMode(true)
	}

	l10n.RegisterCallbacks(db)
	sorting.RegisterCallbacks(db)
	validations.RegisterCallbacks(db)
	audited.RegisterCallbacks(db)
	media.RegisterCallbacks(db)
	return nil
}
