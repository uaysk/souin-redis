package roadrunner

import (
	"github.com/uaysk/souin-redis/pkg/middleware"
	"github.com/uaysk/souin-redis/plugins/souin/agnostic"
)

const (
	configurationKey = "http.cache"
)

// ParseConfiguration parse the Roadrunner configuration into a valid HTTP
// cache configuration object.
func parseConfiguration(cfg Configurer) middleware.BaseConfiguration {
	var configuration middleware.BaseConfiguration
	agnostic.ParseConfiguration(&configuration, cfg.Get(configurationKey).(map[string]interface{}))

	return configuration
}
