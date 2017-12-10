package admin

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/action_bar"
	"github.com/qor/activity"
	"github.com/qor/admin"
	"github.com/qor/help"
	"github.com/qor/i18n/exchange_actions"
	"github.com/qor/media/asset_manager"
	"github.com/qor/media/media_library"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"
	"golang.org/x/crypto/bcrypt"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/ssb-pubmon/config"
	"github.com/cryptix/ssb-pubmon/config/auth"
	"github.com/cryptix/ssb-pubmon/config/i18n"
	"github.com/cryptix/ssb-pubmon/config/notify"
	"github.com/cryptix/ssb-pubmon/db"
	"github.com/cryptix/ssb-pubmon/models"
)

var (
	l         logging.Interface
	inited    bool
	admindb   *gorm.DB
	Admin     *admin.Admin
	ActionBar *action_bar.ActionBar
)

func Init(log logging.Interface) {
	if inited {
		return
	}
	l = log
	admindb = db.GetBase()

	Admin = admin.New(&admin.AdminConfig{
		Auth:     auth.AdminAuth{},
		DB:       admindb,
		SiteName: "ssb-pubmon Admin",
		//AssetFS:  bindatafs.AssetFS.NameSpace("admin"),
		SessionManager: config.SessionManager,
	})

	notify.Init()
	Admin.NewResource(notify.Sender)

	// Add Dashboard
	Admin.AddMenu([]string{"admin"}, &admin.Menu{Name: "Dashboard", Link: "/admin"})

	// Add Media Library
	Admin.AddResource(&media_library.MediaLibrary{}, &admin.Config{Menu: []string{"Site Management"}})

	// Add Asset Manager, for rich editor
	assetManager := Admin.AddResource(&asset_manager.AssetManager{}, &admin.Config{Invisible: true})

	// Add Help
	Help := Admin.NewResource(&help.QorHelpEntry{})
	Help.GetMeta("Body").Config = &admin.RichEditorConfig{AssetManager: assetManager}

	// Add User
	user := Admin.AddResource(&models.User{}, &admin.Config{Menu: []string{"User Management"}})
	user.Meta(&admin.Meta{Name: "Role", Config: &admin.SelectOneConfig{Collection: []string{"Admin", "Maintainer", "Member"}}})
	user.Meta(&admin.Meta{Name: "Password",
		Type:   "password",
		Valuer: func(interface{}, *qor.Context) interface{} { return "" },
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			if newPassword := utils.ToString(metaValue.Value); newPassword != "" {
				bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
				if err != nil {
					context.DB.AddError(validations.NewError(user, "Password", "Can't encrpt password"))
					return
				}
				u := resource.(*models.User)
				u.Password = string(bcryptPassword)
			}
		},
	})
	user.Meta(&admin.Meta{Name: "Confirmed", Valuer: func(user interface{}, ctx *qor.Context) interface{} {
		if user.(*models.User).ID == 0 {
			return true
		}

		/* todo: load fromauthID
		ctx.GetDB()
		return user.(*models.User).Confirmed
		*/
		return false
	}})

	user.Filter(&admin.Filter{
		Name: "Role",
		Config: &admin.SelectOneConfig{
			Collection: []string{"Admin", "Maintainer", "Member"},
		},
	})

	user.IndexAttrs("ID", "Email", "Name", "Role")
	user.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name"},
				{"Email", "Password"},
				{"Avatar"},
				{"Role"},
				{"Confirmed"},
			},
		},
	)
	user.EditAttrs(user.ShowAttrs())

	// pub
	pubRes := Admin.AddResource(&models.Pub{})
	pubRes.Action(&admin.Action{
		Name: "Try",
		Handler: func(argument *admin.ActionArgument) error {
			for _, pub := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				if err := models.PubHealth.Trigger("try", pub.(*models.Pub), db); err != nil {
					return err
				}
				db.Select("state").Save(pub)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if pub, ok := record.(*models.Pub); ok {
				return pub.State != "success"
			}
			return true
		},
		Modes: []string{"index", "show", "menu_item"},
	})
	pubRes.Filter(&admin.Filter{
		Name: "State",
		//Operations: []string{"contains"},
		Config: &admin.SelectOneConfig{Collection: []string{"trying", "worked", "unchecked", "failed"}},
	})
	activity.Register(pubRes)

	pubAddr := Admin.AddResource(&models.Address{})
	pubAddr.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Addr", "Pub"},
				{"LastTry", "Failures", "Took"},
			},
		},
	)
	pubAddr.EditAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"LastTry"},
			},
		},
	)

	// Add Translations
	Admin.AddResource(i18n.I18n, &admin.Config{Menu: []string{"Site Management"}, Priority: 1})

	// Add Worker
	Worker := getWorker()
	exchange_actions.RegisterExchangeJobs(i18n.I18n, Worker)
	Admin.AddResource(Worker, &admin.Config{Menu: []string{"Site Management"}})

	// Add Setting
	Admin.AddResource(&models.Setting{}, &admin.Config{Name: "Site Managment", Singleton: true})

	// Add Search Center Resources
	Admin.AddSearchResource(user)

	// Add ActionBar
	ActionBar = action_bar.New(Admin)
	ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

}
