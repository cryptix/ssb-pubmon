module github.com/qor/auth

require (
	github.com/cryptix/go v1.1.0
	github.com/dgrijalva/jwt-go v0.0.0-20171019145719-dbeaa9332f19a944acb5736b4456cfcc02140e29
	github.com/golang/protobuf v1.1.0 // indirect
	github.com/google/go-github v0.0.0-20180716180158-c0b63e2f9bb1
	github.com/google/go-querystring v0.0.0-20170111101155-53e6ce116135 // indirect
	github.com/jinzhu/copier v0.0.0-20170922082739-db4671f3a9b8
	github.com/jinzhu/gorm v1.9.1
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/mrjones/oauth v0.0.0-20180629183705-f4e24b6d100c
	github.com/pkg/errors v0.8.0
	github.com/qor/mailer v0.0.0-20170814094430-1e6ac7106955
	github.com/qor/qor v0.0.0-20180607095400-a9b667cbbf18
	github.com/qor/redirect_back v0.0.0-20170907030740-b4161ed6f848
	github.com/qor/render v0.0.0-20171201033449-63566e46f01b
	github.com/qor/responder v0.0.0-20171031042654-b6def473574f621fee316696ad120d4fbf470826
	github.com/qor/roles v0.0.0-20171127035124-d6375609fe3e
	github.com/qor/session v0.0.0-20170907035918-8206b0adab70
	golang.org/x/crypto v0.0.0-20180718160520-a2144134853f
	golang.org/x/oauth2 v0.0.0-20180620175406-ef147856a6dd
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f // indirect
)

// forks
replace (
	github.com/qor/admin v0.0.0-20180211171430-1a23d1757f50b38d86f1a23e5da65df179f1d323 => github.com/cryptix/qor_admin v0.0.0-cryptix
	github.com/qor/assetfs v0.0.0-20170713023933-ff57fdc13a14 => github.com/cryptix/qor_assetfs v0.0.1-cryptix
	github.com/qor/mailer v0.0.0-20170814094430-1e6ac7106955 => github.com/cryptix/qor_mailer v1.1.1-cryptix
	github.com/qor/middlewares v0.0.0-20170822143614-781378b69454 => github.com/cryptix/qor_middlewares v0.0.0-cryptix
	github.com/qor/redirect_back v0.0.0-20170907030740-b4161ed6f848 => github.com/cryptix/qor_redirect_back v0.0.0-cryptix
	github.com/qor/render v0.0.0-20171201033449-63566e46f01b => github.com/cryptix/qor_render v1.1.1-cryptix
	github.com/qor/session v0.0.0-20170907035918-8206b0adab70 => github.com/cryptix/qor_session v0.0.1-cryptix
)
