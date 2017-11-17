package l10n_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/l10n"

	_ "github.com/go-sql-driver/mysql"
)

func checkHasErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func checkHasProductInLocale(db *gorm.DB, locale string, t *testing.T) {
	var count int
	if db.Set("l10n:locale", locale).Count(&count); count != 1 {
		t.Errorf("should has only one product for locale %v, but found %v", locale, count)
	}
}

func checkHasProductInAllLocales(db *gorm.DB, t *testing.T) {
	checkHasProductInLocale(db, l10n.Global, t)
	checkHasProductInLocale(db, "zh", t)
	checkHasProductInLocale(db, "en", t)
}

func TestCreateWithCreate(t *testing.T) {
	product := Product{Code: "CreateWithCreate"}
	checkHasErr(t, dbGlobal.Create(&product).Error)
	checkHasErr(t, dbCN.Create(&product).Error)
	checkHasErr(t, dbEN.Create(&product).Error)

	checkHasProductInAllLocales(dbGlobal.Model(&Product{}).Where("id = ? AND code = ?", product.ID, "CreateWithCreate"), t)
}

func TestCreateWithSave(t *testing.T) {
	product := Product{Code: "CreateWithSave"}
	checkHasErr(t, dbGlobal.Create(&product).Error)
	checkHasErr(t, dbCN.Create(&product).Error)
	checkHasErr(t, dbEN.Create(&product).Error)

	checkHasProductInAllLocales(dbGlobal.Model(&Product{}).Where("id = ? AND code = ?", product.ID, "CreateWithSave"), t)
}

func TestUpdate(t *testing.T) {
	product := Product{Code: "Update", Name: "global"}
	checkHasErr(t, dbGlobal.Create(&product).Error)
	sharedDB := dbGlobal.Model(&Product{}).Where("id = ? AND code = ?", product.ID, "Update")

	product.Name = "中文名"
	checkHasErr(t, dbCN.Create(&product).Error)
	checkHasProductInLocale(sharedDB.Where("name = ?", "中文名"), "zh", t)

	product.Name = "English Name"
	checkHasErr(t, dbEN.Create(&product).Error)
	checkHasProductInLocale(sharedDB.Where("name = ?", "English Name"), "en", t)

	product.Name = "新的中文名"
	product.Code = "NewCode // should be ignored when update"
	dbCN.Save(&product)
	checkHasProductInLocale(sharedDB.Where("name = ?", "新的中文名"), "zh", t)

	product.Name = "New English Name"
	product.Code = "NewCode // should be ignored when update"
	dbEN.Save(&product)
	checkHasProductInLocale(sharedDB.Where("name = ?", "New English Name"), "en", t)

	// Check sync columns with UpdateColumn
	dbGlobal.Model(&Product{}).Where("id = ?", product.ID).UpdateColumns(map[string]interface{}{"quantity": gorm.Expr("quantity + ?", 2)})

	var newGlobalProduct Product
	var newENProduct Product
	dbGlobal.Find(&newGlobalProduct, product.ID)
	dbEN.Find(&newENProduct, product.ID)

	if newGlobalProduct.Quantity != product.Quantity+2 || newENProduct.Quantity != product.Quantity+2 {
		t.Errorf("should sync update columns results correctly")
	}

	// Check sync columns with Save
	newGlobalProduct.Quantity = 5
	dbGlobal.Save(&newGlobalProduct)

	var newGlobalProduct2 Product
	var newENProduct2 Product
	dbGlobal.Find(&newGlobalProduct2, product.ID)
	dbEN.Find(&newENProduct2, product.ID)
	if newGlobalProduct2.Quantity != 5 || newENProduct2.Quantity != 5 {
		t.Errorf("should sync update columns results correctly")
	}
}

func TestQuery(t *testing.T) {
	product := Product{Code: "Query", Name: "global"}
	dbGlobal.Create(&product)
	dbCN.Create(&product)

	var productCN Product
	dbCN.First(&productCN, product.ID)
	if productCN.LanguageCode != "zh" {
		t.Error("Should find localized zh product with unscoped mode")
	}

	var newProduct Product
	if dbCN.Set("l10n:mode", "locale").First(&newProduct, product.ID).RecordNotFound() {
		t.Error("Should find localized zh product with locale mode")
	}

	var newProduct2 Product
	if dbCN.Set("l10n:mode", "global").First(&newProduct2); newProduct2.LanguageCode != l10n.Global {
		t.Error("Should find global product with global mode")
	}

	var productEN Product
	dbEN.First(&productEN, product.ID)
	if productEN.LanguageCode != l10n.Global {
		t.Error("Should find global product for en with unscoped mode")
	}

	if !dbEN.Set("l10n:mode", "locale").First(&productEN, product.ID).RecordNotFound() {
		t.Error("Should find no record with locale mode")
	}

	if dbEN.Set("l10n:mode", "global").First(&productEN); productEN.LanguageCode != l10n.Global {
		t.Error("Should find global product with global mode")
	}

	if dbEN.Joins("LEFT JOIN brands ON products.brand_id = brands.id").First(&productEN).Error != nil {
		t.Error("Should handle queries with extra joins")
	}
}

func TestQueryWithPreload(t *testing.T) {
	product := Product{
		Code:       "Query",
		Name:       "global",
		Brand:      Brand{Name: "Brand"},
		Tags:       []Tag{{Name: "tag0"}, {Name: "tag2"}},
		Categories: []Category{{Name: "category1"}, {Name: "category2"}},
	}

	dbGlobal.Create(&product)
	dbCN.Create(&product)

	var productCN Product
	dbCN.Preload("Brand").Preload("Tags").Preload("Categories").First(&productCN, product.ID)

	if (productCN.Brand.LanguageCode != "zh") || len(productCN.Tags) != 2 || len(productCN.Categories) != 2 {
		t.Error("Failed to preload data relations")
	}
}

func TestDelete(t *testing.T) {
	product := Product{Code: "Delete", Name: "global"}
	dbGlobal.Create(&product)
	dbCN.Create(&product)

	if dbCN.Delete(&product).RowsAffected != 1 {
		t.Errorf("Should delete localized record")
	}

	if dbEN.Delete(&product).RowsAffected != 0 {
		t.Errorf("Should delete none record in unlocalized locale")
	}

	// relocalize deleted record
	dbCN.Save(&product)
	var count uint
	if dbCN.Model(&Product{}).Where("code = ? AND name = ?", "Delete", "global").Count(&count); count != 1 {
		t.Errorf("Should be able to relocalize deleted records, get record %v", count)
	}
}

func TestResetLanguageCodeWithGlobalDB(t *testing.T) {
	product := Product{Code: "Query", Name: "global"}
	product.LanguageCode = "test"
	dbGlobal.Save(&product)
	if product.LanguageCode != l10n.Global {
		t.Error("Should reset language code in global mode")
	}
}

func TestManyToManyRelations(t *testing.T) {
	product := Product{Code: "Delete", Name: "global", Tags: []Tag{{Name: "tag1"}, {Name: "tag2"}}}
	dbGlobal.Save(&product)
}
