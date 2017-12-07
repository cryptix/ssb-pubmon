package config

import (
	"html/template"

	"github.com/cryptix/go/logging"
	"github.com/jinzhu/configor"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qor/mailer"
	"github.com/qor/mailer/logger"
	"github.com/qor/redirect_back"
	"github.com/qor/render"
	"github.com/qor/session/manager"

	"github.com/cryptix/ssb-pubmon/config/admin/bindatafs"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

var Config = struct {
	Port int `default:"7000" env:"PORT"`
	DB   struct {
		Name    string `env:"DBName" default:"qor_example"`
		Adapter string
	}
	SMTP SMTPConfig
}{}

var (
	View         *render.Render
	Mailer       *mailer.Mailer
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
	check = logging.CheckFatal
)

func init() {
	err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml")
	check(err)

	View = render.New(&render.Config{
		Layout:          "application",
		ViewPaths:       []string{"app/views"},
		AssetFileSystem: bindatafs.AssetFS,
	})

	htmlSanitizer := bluemonday.UGCPolicy()
	View.RegisterFuncMap("raw", func(str string) template.HTML {
		return template.HTML(htmlSanitizer.Sanitize(str))
	})

	Mailer, err = mailer.New(&mailer.Config{
		AssetFS: bindatafs.AssetFS,
		Sender:  logger.New(&logger.Config{}),
	})
	check(err)
}
