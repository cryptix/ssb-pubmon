package notify

import (
	"time"

	"github.com/cryptix/ssb-pubmon/db"
	"github.com/qor/admin"
	"github.com/qor/notification"
	"github.com/qor/notification/channels/database"
)

var Sender *notification.Notification

func Init() {
	Sender = notification.New(&notification.Config{})
	Sender.RegisterChannel(database.New(&database.Config{DB: db.GetBase()}))
	Sender.Action(&notification.Action{
		Name:         "Dismiss",
		MessageTypes: []string{"connection_try"},
		Visible: func(data *notification.QorNotification, context *admin.Context) bool {
			return data.ResolvedAt == nil
		},
		Handler: func(argument *notification.ActionArgument) error {
			return argument.Context.GetDB().Model(argument.Message).Update("resolved_at", time.Now()).Error
		},
		Undo: func(argument *notification.ActionArgument) error {
			return argument.Context.GetDB().Model(argument.Message).Update("resolved_at", nil).Error
		},
	})
}
