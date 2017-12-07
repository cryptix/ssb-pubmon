package migrations

import (
	"github.com/qor/activity"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/help"
	"github.com/qor/media/asset_manager"
	"github.com/qor/transition"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/ssb-pubmon/db"
	"github.com/cryptix/ssb-pubmon/models"
)

func init() {
	AutoMigrate(&asset_manager.AssetManager{})

	AutoMigrate(&models.Setting{})

	AutoMigrate(&models.User{})

	AutoMigrate(&models.Pub{}, &models.Address{})

	AutoMigrate(&transition.StateChangeLog{})

	AutoMigrate(&activity.QorActivity{})

	AutoMigrate(&models.MediaLibrary{})

	AutoMigrate(&help.QorHelpEntry{})

	AutoMigrate(&auth_identity.AuthIdentity{})
}

var check = logging.CheckFatal

func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		check(db.DB.AutoMigrate(value).Error)
	}
}
