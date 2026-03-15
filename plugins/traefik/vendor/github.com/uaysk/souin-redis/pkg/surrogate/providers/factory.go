package providers

import (
	"github.com/uaysk/souin-redis/configurationtypes"
)

// SurrogateFactory generate a SurrogateInterface instance
func SurrogateFactory(config configurationtypes.AbstractConfigurationInterface, defaultStorerName string) SurrogateInterface {
	cdn := config.GetDefaultCache().GetCDN()

	switch cdn.Provider {
	case "akamai":
		return generateAkamaiInstance(config, defaultStorerName)
	case "cloudflare":
		return generateCloudflareInstance(config, defaultStorerName)
	case "fastly":
		return generateFastlyInstance(config, defaultStorerName)
	default:
		return generateSouinInstance(config, defaultStorerName)
	}
}
