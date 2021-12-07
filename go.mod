module wailik.com

go 1.17

require (
	github.com/deepmap/oapi-codegen v1.8.3
	github.com/getkin/kin-openapi v0.61.0
	github.com/gofiber/fiber/v2 v2.18.0
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.12
	github.com/labstack/echo/v4 v4.2.1
	github.com/open-policy-agent/opa v0.32.0
	github.com/syndtr/goleveldb v1.0.0
	go.uber.org/zap v1.19.1
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.16
)

require (
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/andybalholm/brotli v1.0.2 // indirect
	github.com/dgraph-io/ristretto v0.1.0
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/klauspost/compress v1.13.5 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.29.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.0.0-20180916065949-5c77d914dd0b // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	golang.org/x/text v0.3.6 // indirect
)

require (
	github.com/bluele/gcache v0.0.2
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/rs/xid v1.3.0
	github.com/steambap/captcha v1.4.1
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.274
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses v1.0.274
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.274
)

require (
	github.com/coreos/etcd v3.3.27+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/grpc v1.38.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/coreos/bbolt v1.3.6 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace google.golang.org/grpc v1.38.0 => google.golang.org/grpc v1.26.0

replace github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
