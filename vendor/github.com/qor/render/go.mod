module "github.com/qor/render"

require (
	"github.com/gosimple/slug" v1.1.1
	"github.com/jinzhu/gorm" v0.0.0-20180222040412-48a20a6e9f3f
	"github.com/jinzhu/inflection" v0.0.0-20170102125226-1c35d901db3d
	"github.com/jinzhu/now" v0.0.0-20170212112655-d939ba741945
	"github.com/microcosm-cc/bluemonday" v0.0.0-20171222152607-542fd4642604
	"github.com/qor/assetfs" v0.0.0-20170713023933-ff57fdc13a14
	"github.com/qor/qor" v0.0.0-20180212025102-2d6dc0813f93
	"github.com/rainycape/unidecode" v0.0.0-20150907023854-cb7f23ec59be
	"golang.org/x/net" v0.0.0-20180218175443-cbe0f9307d01
)

replace (
	"github.com/qor/assetfs" v0.0.0-20170713023933-ff57fdc13a14 => "github.com/cryptix/qor_assetfs" v0.0.1-cryptix
)
