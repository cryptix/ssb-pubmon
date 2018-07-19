module github.com/cryptix/ssb-pubmon

require (
	github.com/BurntSushi/toml v0.3.0 // indirect
	github.com/PuerkitoBio/goquery v1.4.1 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/andybalholm/cascadia v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20180315120708-ccb8e960c48f // indirect
	github.com/aws/aws-sdk-go v1.14.29 // indirect
	github.com/cryptix/go v1.3.0
	github.com/dgrijalva/jwt-go v0.0.0-20180308231308-06ea1031745c
	github.com/disintegration/imaging v1.4.2 // indirect
	github.com/dustin/go-humanize v0.0.0-20180713052910-9f541cc9db5d
	github.com/fatih/color v1.7.0 // indirect
	github.com/gin-gonic/gin v1.1.4 // indirect
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/go-ini/ini v1.38.1 // indirect
	github.com/go-kit/kit v0.7.0
	github.com/gopherjs/gopherjs v0.0.0-20180628210949-0892b62f0d9f // indirect
	github.com/gorilla/sessions v1.1.1
	github.com/headzoo/surf v1.0.0 // indirect
	github.com/headzoo/ut v0.0.0-20140828025907-5657594ccf1d // indirect
	github.com/jinzhu/configor v0.0.0-20180614024415-4edaf76fe188
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3 // indirect
	github.com/jinzhu/gorm v1.9.1
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/jteeuwen/go-bindata v0.0.0-20180305030458-6025e8de665b
	github.com/jtolds/gls v0.0.0-20170503224851-77f18212c9c7 // indirect
	github.com/manucorporat/sse v0.0.0-20160126180136-ee05b128a739 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.3 // indirect
	github.com/microcosm-cc/bluemonday v1.0.0
	github.com/miolini/datacounter v0.0.0-20171104152933-fd4e42a1d5e0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/qor/action_bar v0.0.0-20171116074904-e1b463078c28
	github.com/qor/activity v0.0.0-20171031083715-06b5e02b7fa9
	github.com/qor/admin v0.0.0-20180706055925-6ac5bf51bf9c
	github.com/qor/assetfs v0.0.0-20170713023933-ff57fdc13a14
	github.com/qor/audited v0.0.0-20171228121055-b52c9c2f0571
	github.com/qor/auth v0.0.0-20180620080635-2bfb79d96185
	github.com/qor/auth_themes v0.0.0-20171205021512-d9aac462ad96
	github.com/qor/cache v0.0.0-20171031031927-c9d48d1f13ba // indirect
	github.com/qor/filebox v0.0.0-20171031092928-e1210ba127af
	github.com/qor/help v0.0.0-20171031093129-202965d1156a
	github.com/qor/i18n v0.0.0-20180329080744-05c1410442e9
	github.com/qor/l10n v0.0.0-20180330031315-ace6ea40bbe9
	github.com/qor/mailer v0.0.0-20180329083248-0555e49f99ac
	github.com/qor/media v0.0.0-20180514091353-03c15fd0c7b0
	github.com/qor/middlewares v0.0.0-20170822143614-781378b69454
	github.com/qor/notification v0.0.0-20171108104944-2a71c942089c
	github.com/qor/oss v0.0.0-20180613011527-e3d52e4600e3 // indirect
	github.com/qor/publish v0.0.0-20171031094744-36587c0844c8 // indirect
	github.com/qor/qor v0.0.0-20180607095400-a9b667cbbf18
	github.com/qor/redirect_back v0.0.0-20170907030740-b4161ed6f848
	github.com/qor/render v0.0.0-20171201033449-63566e46f01b
	github.com/qor/roles v0.0.0-20171127035124-d6375609fe3e
	github.com/qor/serializable_meta v0.0.0-20180510060738-5fd8542db417 // indirect
	github.com/qor/session v0.0.0-20170907035918-8206b0adab70
	github.com/qor/sorting v0.0.0-20180228073813-8308410fca54
	github.com/qor/transition v0.0.0-20171031101107-4015a3eee19c
	github.com/qor/validations v0.0.0-20171228122639-f364bca61b46
	github.com/qor/wildcard_router v0.0.0-20171031035524-56710e5bb5a4
	github.com/qor/worker v0.0.0-20180524074358-bfb40cda691f
	github.com/smartystreets/assertions v0.0.0-20180607162144-eb5b59917fa2 // indirect
	github.com/smartystreets/goconvey v0.0.0-20180222194500-ef6db91d284a // indirect
	github.com/smartystreets/gunit v0.0.0-20180314194857-6f0d6275bdcd // indirect
	github.com/theplant/htmltestingutils v0.0.0-20171010073838-5313e1d8c06a // indirect
	github.com/theplant/testingutils v0.0.0-20180706140413-80c19d2087bf // indirect
	github.com/yosssi/gohtml v0.0.0-20180130040904-97fbf36f4aa8 // indirect
	go.cryptoscope.co/hexagen v1.0.0
	go.cryptoscope.co/muxrpc v1.0.0
	go.cryptoscope.co/netwrap v0.0.0-20180427130219-dae5b5bc35c3
	go.cryptoscope.co/secretstream v0.0.0-20180615152413-def1de7d6a8b
	golang.org/x/crypto v0.0.0-20180718160520-a2144134853f
	golang.org/x/image v0.0.0-20180708004352-c73c2afc3b81 // indirect
	golang.org/x/sys v0.0.0-20180715085529-ac767d655b30 // indirect
	golang.org/x/text v0.3.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/ini.v1 v1.38.1 // indirect
	gopkg.in/yaml.v2 v2.2.1 // indirect
)

replace (
	github.com/qor/admin v0.0.0-20180706055925-6ac5bf51bf9c => github.com/cryptix/qor_admin v0.0.3-cryptix
	github.com/qor/assetfs v0.0.0-20170713023933-ff57fdc13a14 => github.com/cryptix/qor_assetfs v0.0.1-cryptix
	github.com/qor/auth v0.0.0-20180620080635-2bfb79d96185 => github.com/cryptix/qor_auth v0.0.3-cryptix
	github.com/qor/auth_themes v0.0.0-20171205021512-d9aac462ad96 => github.com/cryptix/qor_auth_themes v0.0.0-cryptix
	github.com/qor/mailer v0.0.0-20180329083248-0555e49f99ac => github.com/cryptix/qor_mailer v1.2.0-cryptix
	github.com/qor/middlewares v0.0.0-20170822143614-781378b69454 => github.com/cryptix/qor_middlewares v0.0.0-cryptix
	github.com/qor/redirect_back v0.0.0-20170907030740-b4161ed6f848 => github.com/cryptix/qor_redirect_back v0.0.0-cryptix
	github.com/qor/render v0.0.0-20171201033449-63566e46f01b => github.com/cryptix/qor_render v1.1.3-cryptix
	github.com/qor/session v0.0.0-20170907035918-8206b0adab70 => github.com/cryptix/qor_session v0.0.1-cryptix
)
