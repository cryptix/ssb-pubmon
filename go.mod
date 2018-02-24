module "github.com/cryptix/ssb-pubmon"

// vgo buggs
replace (
	// I had to checkout and set a tag myself to bypass the resolving of its go.mod file
	"github.com/gorilla/sessions" v0.0.0-20180220211045-7087b4d669d1bc3da42fb4e2eda73ae139a24439 => "../../gorilla/sessions"
)

// doesn't help, see above
require (
	//"github.com/gorilla/securecookie" v0.0.0-20170224193804-e59506cc896a
	//"github.com/gorilla/context" v0.0.0-20160817204632-08b5f424b927
)

// overwrite
require (
	// ours
	"cryptoscope.co/go/hexagen" v0.1.0
	"cryptoscope.co/go/secretstream" v0.9.1

	// defaults
	"github.com/BurntSushi/toml" v0.3.0

	// wrong/old version detect
	"github.com/asaskevich/govalidator" v0.0.0-20180115112450-4b3d68f87f176641ffc147420296013aff66ea32
	"github.com/dgrijalva/jwt-go" v0.0.0-20171019145719-dbeaa9332f19a944acb5736b4456cfcc02140e29
	"github.com/cryptix/go-muxrpc" v1.0.1

	// no go.mod
	"github.com/disintegration/imaging" v1.3.0
	"github.com/go-chi/chi" v0.0.0-20180202194135-e223a795a06a
	"github.com/gorilla/sessions" v0.0.0-20180220211045-7087b4d669d1bc3da42fb4e2eda73ae139a24439
	"github.com/gosimple/slug" v1.1.1
	"github.com/jinzhu/configor" v0.0.0-20171024081003-6ecfe629230f
	"github.com/jinzhu/copier" v0.0.0-20170922082739-db4671f3a9b8
	"github.com/jinzhu/gorm" v0.0.0-20180222050412-48a20a6e9f3f4d26095df82c3337efec6db0a6fc
	"github.com/jinzhu/inflection" v0.0.0-20170102125226-1c35d901db3d
	"github.com/jinzhu/now" v0.0.0-20170212112655-d939ba741945
	"github.com/jteeuwen/go-bindata" v0.0.0-20151023091102-a0ff2567cfb7
	"github.com/lib/pq" v0.0.0-20180201184707-88edab080323
	"github.com/mattn/go-sqlite3" v1.6.0
	"github.com/microcosm-cc/bluemonday" v0.0.0-20171222152607-542fd4642604
	"github.com/rainycape/unidecode" v0.0.0-20150907023854-cb7f23ec59be
	"github.com/theplant/cldr" v0.0.0-20170713054817-b97a4c2a76cf
	"golang.org/x/image" v0.0.0-20171214225156-12117c17ca67
	"golang.org/x/net" v0.0.0-20180218175443-cbe0f9307d01
	"gopkg.in/gomail.v2" v0.0.0-20160411212932-81ebce5c23df
	"gopkg.in/yaml.v2" v1.1.1-gopkgin-v2.1.1

	// QOR
	"github.com/qor/action_bar" v0.0.0-20171116074904-e1b463078c28
	"github.com/qor/activity" v0.0.0-20171031093715-06b5e02b7fa92664750be4526644080806febfd4
	"github.com/qor/admin" v0.0.0-20180211171430-1a23d1757f50b38d86f1a23e5da65df179f1d323
	"github.com/qor/assetfs" v0.0.0-20170713023933-ff57fdc13a14
	"github.com/qor/audited" v0.0.0-20171031102608-b7dc1e8737980ab1a3dd41851a764a983f247960
	"github.com/qor/auth" v0.0.0-20171011070053-d0a475c04fd9
	"github.com/qor/auth_themes" v0.0.0-20171205021512-d9aac462ad96
	"github.com/qor/cache" v0.0.0-20171031041927-c9d48d1f13ba2e1ad06a3f8d55be7ea3edf2d0c4
	"github.com/qor/filebox" v0.0.0-20171031092928-e1210ba127af
	"github.com/qor/help" v0.0.0-20171031093129-202965d1156a
	"github.com/qor/i18n" v0.0.0-20180202103326-ec0ba4042f729b4a0d0026191e1cef46013d0234
	"github.com/qor/l10n" v0.0.0-20171228132013-3ffa522dc571a1c8875465d175049b21a4689ffa
	"github.com/qor/mailer" v0.0.0-20170814094430-1e6ac7106955
	"github.com/qor/media" v0.0.0-20180206082634-c87696a0ab10
	"github.com/qor/media_library" v0.0.0-20171016111834-7c6ac542ddccad82fc4a552333d1af6c20c217c6
	"github.com/qor/middlewares" v0.0.0-20170822143614-781378b69454
	"github.com/qor/notification" v0.0.0-20171108104944-2a71c942089c
	"github.com/qor/oss" v0.0.0-20180122071320-9b2c0a096a42
	"github.com/qor/publish" v0.0.0-20171031104744-36587c0844c8828339e5bb3cbf9d0b22bce3652b
	"github.com/qor/qor" v0.0.0-20180212035102-2d6dc0813f93982324177788466955288aaffe95
	"github.com/qor/redirect_back" v0.0.0-20170907030740-b4161ed6f848
	"github.com/qor/render" v0.0.0-20171201033449-63566e46f01b
	"github.com/qor/responder" v0.0.0-20171031042654-b6def473574f621fee316696ad120d4fbf470826
	"github.com/qor/roles" v0.0.0-20171127045124-d6375609fe3e5da46ad3a574fae244fb633e79c1
	"github.com/qor/serializable_meta" v0.0.0-20171031110819-b432456ad58b075008d8b5633e16bd796926b5f4
	"github.com/qor/session" v0.0.0-20170907035918-8206b0adab70
	"github.com/qor/sorting" v0.0.0-20180111075739-8ebc1045295271c47b2562a13992c22599442ff7
	"github.com/qor/transition" v0.0.0-20171031111107-4015a3eee19c49a63b1d22beab1c0c084e72c53b
	"github.com/qor/validations" v0.0.0-20171228132639-f364bca61b46bd48a5e32552a37758864fdf005d
	"github.com/qor/wildcard_router" v0.0.0-20171031035524-56710e5bb5a4
	"github.com/qor/worker" v0.0.0-20180202090225-61267fc978aa9177939c5d68d2ba87028d89191b
)

// forks
replace (
	"github.com/qor/admin" v1.0.0 => "github.com/cryptix/qor_admin" v0.0.0-20180109004425-e73aa4543e56cafc7a46027d3353f69b64ca91f4
	"github.com/qor/assetfs" v0.0.0-20170713023933-ff57fdc13a14 => "github.com/cryptix/qor_assetfs" v0.0.0-20171203171418-0690e5426d85d229d6e9c919d9637120347058d8
	"github.com/qor/auth" v0.0.0-20171011070053-d0a475c04fd9 => "github.com/cryptix/qor_auth" v0.0.0-20180106160426-35f1d6dea9a2cc13b571e6f6c557437f55655a08
	"github.com/qor/auth_themes" v0.0.0-20171205021512-d9aac462ad96 => "github.com/cryptix/qor_auth_themes" v0.0.0-20180106155717-45a77b9ec82d221cb11936397565ed00150d1c7a
	"github.com/qor/mailer" v0.0.0-20170814094430-1e6ac7106955 => "github.com/cryptix/qor_mailer" v0.0.0-20171203163834-2d5b212956cc0e7c81aa0e4834bfffad8c4ab6ac
	"github.com/qor/middlewares" v0.0.0-20170822143614-781378b69454 => "github.com/cryptix/qor_middlewares" v0.0.0-20180106162028-18127115fc91fd180a122c3a8e6c85078b0624e6
	"github.com/qor/redirect_back" v0.0.0-20170907030740-b4161ed6f848 => "github.com/cryptix/qor_redirect_back" v0.0.0-20180106162833-275539f5ccf24815c8ac88ec932c1f49fd3f8c83
	"github.com/qor/render" v0.0.0-20171201033449-63566e46f01b => "github.com/cryptix/qor_render" v0.0.0-20180106154953-e7333e0f2548
	"github.com/qor/session" v0.0.0-20170907035918-8206b0adab70 => "github.com/cryptix/qor_session" v0.0.0-20180106161732-33f14a897a1f40d1ba043b67e49b4070948f10c0
)
