package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"github.com/qor/wildcard_router"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/ssb-pubmon/config/admin/bindatafs"
	"github.com/cryptix/ssb-pubmon/config/auth"
	"github.com/cryptix/ssb-pubmon/controllers"
	"github.com/cryptix/ssb-pubmon/db"
)

var rootMux *http.ServeMux
var WildcardRouter *wildcard_router.WildcardRouter

func Router(l logging.Interface) *http.ServeMux {
	if rootMux == nil {
		router := chi.NewRouter()

		router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				var (
					tx         = db.GetBase()
					qorContext = &qor.Context{Request: req, Writer: w}
				)

				if locale := utils.GetLocale(qorContext); locale != "" {
					tx = tx.Set("l10n:locale", locale)
				}

				ctx := context.WithValue(req.Context(), utils.ContextDBName, tx)
				next.ServeHTTP(w, req.WithContext(ctx))
			})
		})

		router.Get("/", wwe(l, controllers.Index))
		router.Get("/last", wwe(l, controllers.LastChecks))
		router.Get("/alive", wwe(l, controllers.Alive))
		router.Get("/hexagen", wwe(l, controllers.Hexagen))
		router.Get("/switch_locale", controllers.SwitchLocale)

		router.With(auth.Authority.Authorize()).Route("/account", func(r chi.Router) {
			r.Get("/", controllers.AccountShow)
			//r.Post("/profile", controllers.SetUserProfile)
		})

		rootMux = http.NewServeMux()

		rootMux.Handle("/auth/", auth.Auth.NewServeMux())
		//rootMux.Handle("/system/", utils.FileServer(http.Dir(filepath.Join(config.Root, "public"))))
		assetFS := bindatafs.AssetFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images", "dist", "fonts", "vendors")
		for _, path := range []string{"javascripts", "stylesheets", "images", "dist", "fonts", "vendors"} {
			rootMux.Handle(fmt.Sprintf("/%s/", path), assetFS)
		}

		WildcardRouter = wildcard_router.New()
		WildcardRouter.MountTo("/", rootMux)
		WildcardRouter.AddHandler(router)
	}
	return rootMux
}

type myHandler func(w http.ResponseWriter, req *http.Request) error

func wwe(log logging.Interface, fn myHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := fn(w, req)
		if err != nil {
			if _, ok := err.(*controllers.BadReqError); ok {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			if l := logging.FromContext(req.Context()); l != nil {
				log = l
			}
			log.Log("event", "error", "msg", "view execution failed",
				"path", req.URL.Path,
				"err", err)
			fmt.Fprint(w, err)
		}
	}
}
