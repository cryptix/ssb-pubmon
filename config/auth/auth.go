package auth

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/authority"
	"github.com/qor/auth_themes/clean"
	"github.com/qor/middlewares"
	"github.com/qor/redirect_back"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/ssb-pubmon/config"
	"github.com/cryptix/ssb-pubmon/db"
	"github.com/cryptix/ssb-pubmon/models"
)

var (
	// Auth initialize Auth for Authentication
	Auth *auth.Auth
	// Authority initialize Authority for Authorization
	Authority *authority.Authority
)

func Init() {
	var err error
	Auth = clean.New(&auth.Config{
		DB:         db.GetBase(),
		Render:     config.View,
		ViewPaths:  []string{filepath.Join(config.Root, "app", "views")},
		Mailer:     config.Mailer,
		UserModel:  models.User{},
		Redirector: Redirector{RedirectBack: config.RedirectBack},
		SessionStorer: &auth.SessionStorer{
			SessionName:    "_auth_session",
			SessionManager: config.SessionManager,
			SigningMethod:  jwt.SigningMethodHS256,
		},
	})
	logging.CheckFatal(err)

	Authority = authority.New(&authority.Config{
		Auth:                Auth,
		AccessDeniedHandler: authority.NewAccessDeniedHandler(Auth, "/auth/login"),
	})
	config.Middlewares.Use(middlewares.Middleware{
		Name:        "authority",
		InsertAfter: []string{"session"},
		Handler: func(handler http.Handler) http.Handler {
			return Authority.Middleware(handler)
		},
	})
	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})
}

// Redirector our custom redirector for /checkin
type Redirector struct {
	*redirect_back.RedirectBack
}

// Redirect redirect back after action
func (redirector Redirector) Redirect(w http.ResponseWriter, req *http.Request, action string) {
	u, userOk := Auth.GetCurrentUser(req).(*models.User)
	switch action {
	case "login":
		if userOk {
			var authInfo auth_identity.Basic

			authInfo.Provider = "password"
			authInfo.UID = u.Email

			tx := db.GetBase()

			if tx.Table("auth_identities").Where(authInfo).Scan(&authInfo).RecordNotFound() {
				log.Print("ERROR: could not find authInfo for logged in user!?")
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			// confirmed but not checked in?
			if authInfo.ConfirmedAt != nil && u.Role == "" {
				http.Redirect(w, req, "/account/checkin", http.StatusSeeOther)
				return
			}
		}
	case "reset_password":
		fallthrough
	case "send_reset_mail":
		fallthrough
	case "confirmation_send":
		http.Redirect(w, req, "/auth/login", http.StatusSeeOther)
		return
	case "confirm":
		http.Redirect(w, req, "/account/checkin", http.StatusSeeOther)
		return
	}
	redirector.RedirectBack.RedirectBack(w, req)
}
