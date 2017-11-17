package migrations

import (
	"github.com/qor/activity"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/banner_editor"
	"github.com/qor/help"
	"github.com/qor/media/asset_manager"
	"github.com/qor/transition"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/synchrotron/app/models"
	"github.com/cryptix/synchrotron/config/admin"
	"github.com/cryptix/synchrotron/db"
)

func init() {
	AutoMigrate(&asset_manager.AssetManager{})

	AutoMigrate(&models.Product{}, &models.ProductVariation{}, &models.ProductImage{}, &models.ColorVariation{}, &models.ColorVariationImage{}, &models.SizeVariation{})
	AutoMigrate(&models.Color{}, &models.Size{}, &models.Material{}, &models.Category{}, &models.Collection{})

	AutoMigrate(&models.Address{})

	AutoMigrate(&models.Order{}, &models.OrderItem{})

	AutoMigrate(&models.Store{})

	AutoMigrate(&models.Setting{})

	AutoMigrate(&models.User{})

	AutoMigrate(&transition.StateChangeLog{})

	AutoMigrate(&activity.QorActivity{})

	AutoMigrate(&admin.QorWidgetSetting{})

	AutoMigrate(&models.Page{})

	AutoMigrate(&models.MediaLibrary{})

	AutoMigrate(&models.Article{})

	AutoMigrate(&help.QorHelpEntry{})

	AutoMigrate(&auth_identity.AuthIdentity{})

	AutoMigrate(&banner_editor.QorBannerEditorSetting{})
}

var check = logging.CheckFatal

func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		check(db.DB.AutoMigrate(value).Error)
	}
}
