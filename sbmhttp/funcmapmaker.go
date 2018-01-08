package sbmhttp

import (
	"html/template"
	"net/http"

	"github.com/qor/i18n/inline_edit"
	"github.com/qor/render"
	"github.com/qor/session"

	"github.com/cryptix/ssb-pubmon/config"
	"github.com/cryptix/ssb-pubmon/config/admin"
	"github.com/cryptix/ssb-pubmon/config/i18n"
	"github.com/cryptix/ssb-pubmon/config/utils"
	"github.com/cryptix/ssb-pubmon/models"
)

func fmm(render *render.Render, req *http.Request, w http.ResponseWriter) template.FuncMap {
	funcMap := template.FuncMap{}

	// Add `t` method
	for key, fc := range inline_edit.FuncMap(i18n.I18n, utils.GetCurrentLocale(req), utils.GetEditMode(w, req)) {
		funcMap[key] = fc
	}

	for key, value := range admin.ActionBar.FuncMap(w, req) {
		funcMap[key] = value
	}

	funcMap["flashes"] = func() []session.Message {
		return config.SessionManager.Flashes(w, req)
	}

	// Add `action_bar` method
	funcMap["render_action_bar"] = func() template.HTML {
		return admin.ActionBar.Render(w, req)
	}

	funcMap["current_locale"] = func() string {
		return utils.GetCurrentLocale(req)
	}

	funcMap["current_user"] = func() *models.User {
		return utils.GetCurrentUser(req)
	}

	return funcMap
}
