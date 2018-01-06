package clean

import (
	"github.com/pkg/errors"
	"github.com/qor/auth"
	"github.com/qor/auth/claims"
	"github.com/qor/auth/providers/password"
)

// ErrPasswordConfirmationNotMatch password confirmation not match error
var ErrPasswordConfirmationNotMatch = errors.New("password confirmation doesn't match password")

// New initialize clean theme
func New(config *auth.Config) (*auth.Auth, error) {
	if config == nil {
		config = &auth.Config{}
	}
	config.ViewPaths = append(config.ViewPaths, "github.com/qor/auth_themes/clean/views")

	if config.DB == nil {
		return nil, errors.New("cleanTheme: Please configure *gorm.DB")
	}

	if config.Render == nil {
		return nil, errors.New("cleanTheme: Please configure Render")
	}

	a, err := auth.New(config)
	if err != nil {
		return nil, errors.Wrap(err, "cleanTheme: failed to create auth system")
	}

	a.RegisterProvider(password.New(&password.Config{
		Confirmable: true,
		RegisterHandler: func(context *auth.Context) (*claims.Claims, error) {
			context.Request.ParseForm()

			if context.Request.Form.Get("confirm_password") != context.Request.Form.Get("password") {
				return nil, ErrPasswordConfirmationNotMatch
			}

			return password.DefaultRegisterHandler(context)
		},
	}))

	if a.Config.DB != nil {
		// Migrate Auth Identity model
		a.Config.DB.AutoMigrate(a.Config.AuthIdentityModel)
	}
	return a, nil
}
