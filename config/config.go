package config

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/cryptix/go/goutils"
	"github.com/cryptix/go/logging"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/configor"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qor/mailer"
	"github.com/qor/mailer/logger"
	"github.com/qor/middlewares"
	"github.com/qor/redirect_back"
	"github.com/qor/render"
	"github.com/qor/session"
	"github.com/qor/session/gorilla"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type MainConfig struct {
	CookieSecret string `env:"COOKIESecret"` // dd if=/dev/urandom bs=1 count=64 | base64 -w0
	Locale       string `default:"de-DE"`
	HTTPHost     string `default:":7000"`
	DB           struct {
		Name    string `env:"DBName" default:"ssbpub.db"`
		Adapter string
	}
	SMTP SMTPConfig
}

var (
	Config         *MainConfig
	Root           string
	View           *render.Render
	Mailer         *mailer.Mailer
	SessionManager session.ManagerInterface
	RedirectBack   *redirect_back.RedirectBack
	Middlewares    *middlewares.MiddlewareStack

	check = logging.CheckFatal
)

func Init(log logging.Interface) {
	if Config != nil {
		return
	}
	Config = new(MainConfig)
	err := configor.Load(Config, "config/database.yml", "config/smtp.yml")
	check(err)

	Root, err = goutils.LocatePackage("github.com/cryptix/ssb-pubmon")
	//check(err)
	if err != nil {
		log.Log("event", "LocatePackage error", "err", err, "msg", "falling back to '.'")
		Root = "."
	}

	Middlewares = new(middlewares.MiddlewareStack)
	Middlewares.Use(middlewares.Middleware{
		Name:    "recovery",
		Handler: logging.RecoveryHandler(),
	})

	cookieSecret, err := base64.StdEncoding.DecodeString(Config.CookieSecret)
	check(err)

	if len(cookieSecret) != 64 {
		panic("cookie secret too short")
	}

	SessionManager = gorilla.New("_session", sessions.NewCookieStore(cookieSecret))
	Middlewares.Use(middlewares.Middleware{
		Name: "session",
		Handler: func(next http.Handler) http.Handler {
			return SessionManager.Middleware(next)
		},
	})

	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
	Middlewares.Use(middlewares.Middleware{
		Name:        "redirect_back",
		InsertAfter: []string{"session"},
		Handler: func(next http.Handler) http.Handler {
			return RedirectBack.Middleware(next)
		},
	})

	View = render.New(&render.Config{
		Layout:    "application",
		ViewPaths: []string{filepath.Join(Root, "app", "views")},
	})

	htmlSanitizer := bluemonday.UGCPolicy()
	View.RegisterFuncMap("raw", func(str string) template.HTML {
		return template.HTML(htmlSanitizer.Sanitize(str))
	})

	/*
		fromAddr, err := mail.ParseAddress("do-not-reply <mail@z>")
		if err != nil {
			check(errors.Wrap(err, "failed to parse from addr"))
		}
	*/
	Mailer, err = mailer.New(&mailer.Config{
		DefaultEmailTemplate: &mailer.Email{},
		Sender:               logger.New(nil),
	})
	log.Log("event", "mailer.New", "err", err)
	Mailer.RegisterViewPath(filepath.Join(Root, "app", "views"))
}
