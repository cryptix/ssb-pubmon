package config

import (
	"html/template"
	"path/filepath"

	"github.com/cryptix/go/goutils"
	"github.com/jinzhu/configor"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qor/auth/providers/github"
	"github.com/qor/auth/providers/twitter"
	"github.com/qor/mailer"
	"github.com/qor/mailer/logger"
	"github.com/qor/redirect_back"
	"github.com/qor/render"
	"github.com/qor/session/manager"
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
		Name string `env:"DBName" default:"qor_example"`
	}
	SMTP    SMTPConfig
	Github  github.Config
	Twitter twitter.Config
}{}

var (
	Root         = goutils.MustLocatePackage("github.com/cryptix/synchrotron")
	View         *render.Render
	Mailer       *mailer.Mailer
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

func init() {
	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		panic(err)
	}

	View = render.New(&render.Config{
		Layout:    "application",
		ViewPaths: []string{filepath.Join(Root, "app", "views")},
	})

	htmlSanitizer := bluemonday.UGCPolicy()
	View.RegisterFuncMap("raw", func(str string) template.HTML {
		return template.HTML(htmlSanitizer.Sanitize(str))
	})

	// dialer := gomail.NewDialer(Config.SMTP.Host, Config.SMTP.Port, Config.SMTP.User, Config.SMTP.Password)
	// sender, err := dialer.Dial()

	// Mailer = mailer.New(&mailer.Config{
	// 	Sender: gomailer.New(&gomailer.Config{Sender: sender}),
	// })
	Mailer = mailer.New(&mailer.Config{
		Sender: logger.New(&logger.Config{}),
	})
}
