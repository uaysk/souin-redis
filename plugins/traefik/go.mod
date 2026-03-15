module github.com/uaysk/souin-redis/plugins/traefik

go 1.24.5

require (
	cel.dev/expr v0.19.1
	dario.cat/mergo v1.0.1
	filippo.io/edwards25519 v1.1.0
	github.com/KimMachineGun/automemlimit v0.7.1
	github.com/Masterminds/goutils v1.1.1
	github.com/Masterminds/semver/v3 v3.3.0
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/Microsoft/go-winio v0.6.0
	github.com/akyoto/cache v1.0.6
	github.com/antlr4-go/antlr/v4 v4.13.0
	github.com/aryann/difflib v0.0.0-20210328193216-ff5ff6dc229b
	github.com/beorn7/perks v1.0.1
	github.com/caddyserver/caddy/v2 v2.10.0
	github.com/caddyserver/certmagic v0.23.0
	github.com/caddyserver/zerossl v0.1.3
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/chzyer/readline v1.5.1
	github.com/cloudflare/circl v1.6.0
	github.com/cpuguy83/go-md2man/v2 v2.0.6
	github.com/darkweak/go-esi v0.0.5
	github.com/darkweak/storages/core v0.0.18
	github.com/darkweak/storages/redis v0.0.13
	github.com/francoispqt/gojay v1.2.13
	github.com/go-jose/go-jose/v3 v3.0.4
	github.com/go-kit/kit v0.13.0
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572
	github.com/google/cel-go v0.24.1
	github.com/google/pprof v0.0.0-20231212022811-ec68065c825e
	github.com/google/uuid v1.6.0
	github.com/huandu/xstrings v1.5.0
	github.com/inconshreveable/mousetrap v1.1.0
	github.com/klauspost/cpuid/v2 v2.2.10
	github.com/libdns/libdns v1.0.0-beta.1
	github.com/manifoldco/promptui v0.9.0
	github.com/mholt/acmez/v3 v3.1.2
	github.com/miekg/dns v1.1.63
	github.com/mitchellh/copystructure v1.2.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/onsi/ginkgo/v2 v2.15.0
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/pierrec/lz4/v4 v4.1.23
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.2.0
	github.com/prometheus/client_golang v1.22.0
	github.com/prometheus/client_model v0.6.2
	github.com/prometheus/common v0.62.0
	github.com/prometheus/procfs v0.15.1
	github.com/quic-go/qpack v0.5.1
	github.com/quic-go/quic-go v0.50.1
	github.com/redis/rueidis v1.0.54
	github.com/rs/xid v1.5.0
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/shopspring/decimal v1.4.0
	github.com/slackhq/nebula v1.6.1
	github.com/smallstep/certificates v0.26.1
	github.com/smallstep/nosql v0.6.1
	github.com/smallstep/pkcs7 v0.0.0-20231024181729-3b98ecc1ca81
	github.com/smallstep/scep v0.0.0-20231024192529-aee96d7ad34d
	github.com/smallstep/truststore v0.13.0
	github.com/spf13/cast v1.7.0
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.6
	github.com/stoewer/go-strcase v1.2.0
	github.com/tailscale/tscert v0.0.0-20240608151842-d3f834017e53
	github.com/uaysk/souin-redis v1.7.8
	github.com/urfave/cli v1.22.14
	github.com/zeebo/blake3 v0.2.4
	go.step.sm/cli-utils v0.9.0
	go.step.sm/crypto v0.45.0
	go.step.sm/linkedca v0.20.1
	go.uber.org/automaxprocs v1.6.0
	go.uber.org/mock v0.5.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	go.uber.org/zap/exp v0.3.0
	golang.org/x/crypto v0.36.0
	golang.org/x/crypto/x509roots/fallback v0.0.0-20250305170421-49bf5b80c810
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842
	golang.org/x/mod v0.24.0
	golang.org/x/net v0.38.0
	golang.org/x/sync v0.14.0
	golang.org/x/sys v0.31.0
	golang.org/x/term v0.30.0
	golang.org/x/text v0.23.0
	golang.org/x/time v0.11.0
	golang.org/x/tools v0.31.0
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 v3.0.1
	howett.net/plist v1.0.0
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4 // indirect
	github.com/dgraph-io/ristretto v0.2.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	go.etcd.io/bbolt v1.3.9 // indirect
)

replace (
	github.com/uaysk/souin-redis v1.7.8 => ../..
	go.uber.org/zap v1.26.0 => go.uber.org/zap v1.21.0
)
