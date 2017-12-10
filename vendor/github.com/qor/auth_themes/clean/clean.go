package clean

import (
	"errors"

	"github.com/qor/auth"
	"github.com/qor/auth/claims"
	"github.com/qor/auth/providers/password"
)

// ErrPasswordConfirmationNotMatch password confirmation not match error
var ErrPasswordConfirmationNotMatch = errors.New("password confirmation doesn't match password")

// New initialize clean theme
func New(config *auth.Config) *auth.Auth {
	if config == nil {
		config = &auth.Config{}
	}
	config.ViewPaths = append(config.ViewPaths, "github.com/qor/auth_themes/clean/views")

	if config.DB == nil {
		panic("Please configure *gorm.DB for Auth theme clean")
	}

	if config.Render == nil {
		panic("Please configure renderer for Auth theme clean")
	}

	Auth, err := auth.New(config)
	if err != nil {
		panic(err)
	}

	Auth.RegisterProvider(password.New(&password.Config{
		Confirmable: true,
		RegisterHandler: func(context *auth.Context) (*claims.Claims, error) {
			context.Request.ParseForm()

			if context.Request.Form.Get("confirm_password") != context.Request.Form.Get("password") {
				return nil, ErrPasswordConfirmationNotMatch
			}

			return password.DefaultRegisterHandler(context)
		},
	}))

	if Auth.Config.DB != nil {
		// Migrate Auth Identity model
		Auth.Config.DB.AutoMigrate(Auth.Config.AuthIdentityModel)
	}
	return Auth
}
