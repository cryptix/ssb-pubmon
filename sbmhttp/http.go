package sbmhttp

import (
	"fmt"
	"net/http"

	"github.com/cryptix/go/logging"
	kitlog "github.com/go-kit/kit/log"

	"github.com/cryptix/ssb-pubmon/config"
	"github.com/cryptix/ssb-pubmon/config/admin"
	"github.com/cryptix/ssb-pubmon/config/auth"
	"github.com/cryptix/ssb-pubmon/config/i18n"
	"github.com/cryptix/ssb-pubmon/config/routes"
	"github.com/cryptix/ssb-pubmon/db"
	"github.com/cryptix/ssb-pubmon/db/migrations"
	"github.com/cryptix/ssb-pubmon/models"
)

func loggingHandler(l logging.Interface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			l = kitlog.With(l, "urlPath", req.URL.Path)
			if u := auth.Auth.GetCurrentUser(req); u != nil {
				l = kitlog.With(l, "user", u.(*models.User).ID)
			}
			ctx = logging.NewContext(ctx, l)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func InitServ(l kitlog.Logger, ver string) http.Handler {
	config.Init(kitlog.With(l, "unit", "config"))
	logging.CheckFatal(db.Open()) // todo pass back error for cleaner shutdown
	migrations.Migrate()
	i18n.Init()
	auth.Init()
	admin.Init(kitlog.With(l, "unit", "admin"))
	mux := http.NewServeMux()
	httpLog := kitlog.With(l, "unit", "http")
	mux.Handle("/", routes.Router(httpLog))
	mux.HandleFunc("/version", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Version: %s", ver)
	})
	admin.Admin.MountTo("/admin", mux)
	admin.Filebox.MountTo("/downloads", mux)
	config.View.FuncMapMaker = fmm
	return loggingHandler(httpLog)(config.Middlewares.Apply(mux))
}
