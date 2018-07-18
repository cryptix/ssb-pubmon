module "github.com/qor/auth"

require (
	"github.com/dgrijalva/jwt-go" v0.0.0-20171019145719-dbeaa9332f19a944acb5736b4456cfcc02140e29
	"github.com/gorilla/sessions" v0.0.0-20180220211045-7087b4d669d1bc3da42fb4e2eda73ae139a24439
	"github.com/jinzhu/copier" v0.0.0-20170922082739-db4671f3a9b8
	"github.com/jinzhu/gorm" v0.0.0-20180222050412-48a20a6e9f3f4d26095df82c3337efec6db0a6fc
	"github.com/qor/mailer" v0.0.0-20170814094430-1e6ac7106955
	"github.com/qor/middlewares" v0.0.0-20170822143614-781378b69454
	"github.com/qor/qor" v0.0.0-20180212035102-2d6dc0813f93982324177788466955288aaffe95
	"github.com/qor/redirect_back" v0.0.0-20170907030740-b4161ed6f848
	"github.com/qor/responder" v0.0.0-20171031042654-b6def473574f621fee316696ad120d4fbf470826
	"github.com/qor/session" v0.0.0-20170907035918-8206b0adab70
	"gopkg.in/gomail.v2" v0.0.0-20160411212932-81ebce5c23df
)

// forks
replace (
	"github.com/gorilla/sessions" v0.0.0-20180220211045-7087b4d669d1bc3da42fb4e2eda73ae139a24439 => "github.com/cryptix/gorilla_sessions" v1.1.1-vfork
	"github.com/qor/admin" v0.0.0-20180211171430-1a23d1757f50b38d86f1a23e5da65df179f1d323 => "github.com/cryptix/qor_admin" v0.0.0-cryptix
	"github.com/qor/assetfs" v0.0.0-20170713023933-ff57fdc13a14 => "github.com/cryptix/qor_assetfs" v0.0.1-cryptix
	"github.com/qor/mailer" v0.0.0-20170814094430-1e6ac7106955 => "github.com/cryptix/qor_mailer" v1.1.1-cryptix
	"github.com/qor/middlewares" v0.0.0-20170822143614-781378b69454 => "github.com/cryptix/qor_middlewares" v0.0.0-20180106162028-18127115fc91fd180a122c3a8e6c85078b0624e6
	"github.com/qor/redirect_back" v0.0.0-20170907030740-b4161ed6f848 => "github.com/cryptix/qor_redirect_back" v0.0.0-20180106162833-275539f5ccf24815c8ac88ec932c1f49fd3f8c83
	"github.com/qor/render" v0.0.0-20171201033449-63566e46f01b => "github.com/cryptix/qor_render" v1.1.1-cryptix
	"github.com/qor/session" v0.0.0-20170907035918-8206b0adab70 => "github.com/cryptix/qor_session" v0.0.0-20180106161732-33f14a897a1f40d1ba043b67e49b4070948f10c0
)
