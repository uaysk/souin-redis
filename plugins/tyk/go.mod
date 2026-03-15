module github.com/uaysk/souin-redis/plugins/tyk

go 1.25

require (
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/uaysk/souin-redis v1.7.8
	github.com/uaysk/souin-redis/plugins/souin v1.7.8
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pquerna/cachecontrol v0.2.0
	go.uber.org/zap v1.27.0
)

replace (
	github.com/uaysk/souin-redis v1.7.8 => ../..
	github.com/uaysk/souin-redis/plugins/souin => ../souin
	github.com/uaysk/souin-redis/plugins/souin/storages => ../souin/storages
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 => github.com/alecthomas/kingpin/v2 v2.3.2
)
