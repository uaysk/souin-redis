package traefik

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/uaysk/souin-redis/configurationtypes"
	"github.com/uaysk/souin-redis/pkg/middleware"
)

// SouinTraefikMiddleware declaration.
type SouinTraefikMiddleware struct {
	next http.Handler
	name string
	*middleware.SouinBaseHandler
}

// TestConfiguration is the temporary configuration for Træfik
type TestConfiguration map[string]interface{}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *TestConfiguration {
	return &TestConfiguration{}
}

func configCacheKey(keyConfiguration map[string]interface{}) configurationtypes.Key {
	key := configurationtypes.Key{}
	for keyK, keyV := range keyConfiguration {
		switch keyK {
		case "disable_body":
			key.DisableBody = cast.ToBool(keyV)
		case "disable_host":
			key.DisableHost = cast.ToBool(keyV)
		case "disable_method":
			key.DisableMethod = cast.ToBool(keyV)
		case "disable_query":
			key.DisableQuery = cast.ToBool(keyV)
		case "disable_scheme":
			key.DisableScheme = cast.ToBool(keyV)
		case "disable_vary":
			key.DisableVary = cast.ToBool(keyV)
		case "hash":
			key.Hash = true
		case "headers":
			key.Headers = parseStringSlice(keyV)
		case "hide":
			key.Hide = cast.ToBool(keyV)
		case "template":
			key.Template = cast.ToString(keyV)
		}
	}

	return key
}

func parseProviderConfiguration(providerConfiguration map[string]interface{}) map[string]interface{} {
	parsed := make(map[string]interface{}, len(providerConfiguration))
	for key, value := range providerConfiguration {
		if nested := asStringMap(value); len(nested) != 0 {
			parsed[key] = parseProviderConfiguration(nested)
			continue
		}
		parsed[key] = value
	}

	return parsed
}

func parseCacheProvider(raw interface{}) configurationtypes.CacheProvider {
	provider := configurationtypes.CacheProvider{}
	providerConfiguration := asStringMap(raw)
	if len(providerConfiguration) == 0 {
		return provider
	}

	for providerKey, providerValue := range providerConfiguration {
		switch providerKey {
		case "url":
			provider.URL = cast.ToString(providerValue)
		case "path":
			provider.Path = cast.ToString(providerValue)
		case "configuration":
			if configuration := asStringMap(providerValue); len(configuration) != 0 {
				provider.Configuration = parseProviderConfiguration(configuration)
			}
		}
	}

	return provider
}

func asStringMap(raw interface{}) map[string]interface{} {
	if raw == nil {
		return map[string]interface{}{}
	}

	return cast.ToStringMap(raw)
}

func parseConfiguration(c map[string]interface{}) Configuration {
	configuration := Configuration{}

	for k, v := range c {
		switch k {
		case "api":
			var a configurationtypes.API
			var prometheusConfiguration, souinConfiguration map[string]interface{}
			apiConfiguration := cast.ToStringMap(v)
			for apiK, apiV := range apiConfiguration {
				switch apiK {
				case "prometheus":
					prometheusConfiguration = make(map[string]interface{})
					if apiV != nil {
						prometheus := cast.ToStringMap(apiV)
						if len(prometheus) != 0 {
							prometheusConfiguration = prometheus
						}
					}
				case "souin":
					souinConfiguration = make(map[string]interface{})
					if apiV != nil {
						souin := cast.ToStringMap(apiV)
						if len(souin) != 0 {
							souinConfiguration = souin
						}
					}
				}
			}
			if prometheusConfiguration != nil {
				a.Prometheus = configurationtypes.APIEndpoint{}
				a.Prometheus.Enable = true
				if prometheusConfiguration["basepath"] != nil {
					a.Prometheus.BasePath = cast.ToString(prometheusConfiguration["basepath"])
				}
			}
			if souinConfiguration != nil {
				a.Souin = configurationtypes.APIEndpoint{}
				a.Souin.Enable = true
				if souinConfiguration["basepath"] != nil {
					a.Souin.BasePath = cast.ToString(souinConfiguration["basepath"])
				}
			}
			configuration.API = a
		case "cache_keys":
			cacheKeys := make(configurationtypes.CacheKeys, 0)
			cacheKeyConfiguration := cast.ToStringMap(v)
			for cacheKeyConfigurationK, cacheKeyConfigurationV := range cacheKeyConfiguration {
				cacheKeyK := configurationtypes.RegValue{
					Regexp: regexp.MustCompile(cacheKeyConfigurationK),
				}
				cacheKeyV := configCacheKey(cast.ToStringMap(cacheKeyConfigurationV))
				cacheKeys = append(cacheKeys, configurationtypes.CacheKey{
					cacheKeyK: cacheKeyV,
				})
			}
			configuration.CacheKeys = cacheKeys
		case "default_cache":
			dc := configurationtypes.DefaultCache{
				Distributed: false,
				Headers:     []string{},
				Olric: configurationtypes.CacheProvider{
					URL:           "",
					Path:          "",
					Configuration: nil,
				},
				Regex:               configurationtypes.Regex{},
				TTL:                 configurationtypes.Duration{},
				DefaultCacheControl: "",
			}
			defaultCache := cast.ToStringMap(v)
			for defaultCacheK, defaultCacheV := range defaultCache {
				switch defaultCacheK {
				case "cache_name":
					dc.CacheName = cast.ToString(defaultCacheV)
				case "cdn":
					cdn := configurationtypes.CDN{
						Dynamic: true,
					}
					cdnConfiguration := cast.ToStringMap(defaultCacheV)
					for cdnK, cdnV := range cdnConfiguration {
						switch cdnK {
						case "api_key":
							cdn.APIKey = cast.ToString(cdnV)
						case "dynamic":
							cdn.Dynamic = cast.ToBool(cdnV)
						case "email":
							cdn.Email = cast.ToString(cdnV)
						case "hostname":
							cdn.Hostname = cast.ToString(cdnV)
						case "network":
							cdn.Network = cast.ToString(cdnV)
						case "provider":
							cdn.Provider = cast.ToString(cdnV)
						case "service_id":
							cdn.ServiceID = cast.ToString(cdnV)
						case "strategy":
							cdn.Strategy = cast.ToString(cdnV)
						case "zone_id":
							cdn.ZoneID = cast.ToString(cdnV)
						}
					}
					dc.CDN = cdn
				case "headers":
					dc.Headers = parseStringSlice(defaultCacheV)
				case "key":
					dc.Key = configCacheKey(cast.ToStringMap(defaultCacheV))
				case "mode":
					dc.Mode = cast.ToString(defaultCacheV)
				case "redis":
					dc.Distributed = true
					dc.Redis = parseCacheProvider(defaultCacheV)
				case "regex":
					exclude := cast.ToString(cast.ToStringMap(defaultCacheV)["exclude"])
					if exclude != "" {
						dc.Regex = configurationtypes.Regex{Exclude: exclude}
					}
				case "timeout":
					timeout := configurationtypes.Timeout{}
					timeoutConfiguration := cast.ToStringMap(defaultCacheV)
					for timeoutK, timeoutV := range timeoutConfiguration {
						switch timeoutK {
						case "backend":
							d := configurationtypes.Duration{}
							ttl, err := time.ParseDuration(cast.ToString(timeoutV))
							if err == nil {
								d.Duration = ttl
							}
							timeout.Backend = d
						case "cache":
							d := configurationtypes.Duration{}
							ttl, err := time.ParseDuration(cast.ToString(timeoutV))
							if err == nil {
								d.Duration = ttl
							}
							timeout.Cache = d
						}
					}
					dc.Timeout = timeout
				case "ttl":
					ttl, err := time.ParseDuration(cast.ToString(defaultCacheV))
					if err == nil {
						dc.TTL = configurationtypes.Duration{Duration: ttl}
					}
				case "allowed_http_verbs":
					dc.AllowedHTTPVerbs = parseStringSlice(defaultCacheV)
				case "allowed_additional_status_codes":
					dc.AllowedAdditionalStatusCodes = parseIntSlice(defaultCacheV)
				case "stale":
					stale, err := time.ParseDuration(defaultCacheV.(string))
					if err == nil {
						dc.Stale = configurationtypes.Duration{Duration: stale}
					}
				case "storers":
					dc.Storers = parseStringSlice(defaultCacheV)
				case "default_cache_control":
					dc.DefaultCacheControl = cast.ToString(defaultCacheV)
				case "max_cachable_body_bytes":
					dc.MaxBodyBytes = parseUint64(defaultCacheV)
				}
			}
			configuration.DefaultCache = &dc
		case "log_level":
			configuration.LogLevel = cast.ToString(v)
		case "urls":
			u := make(map[string]configurationtypes.URL)
			urls := cast.ToStringMap(v)

			for urlK, urlV := range urls {
				currentURL := configurationtypes.URL{
					TTL:     configurationtypes.Duration{},
					Headers: nil,
				}
				currentValue := cast.ToStringMap(urlV)
				currentURL.Headers = parseStringSlice(currentValue["headers"])
				d := cast.ToString(currentValue["ttl"])
				ttl, err := time.ParseDuration(d)
				if err == nil {
					currentURL.TTL = configurationtypes.Duration{Duration: ttl}
				}
				if _, exists := currentValue["default_cache_control"]; exists {
					currentURL.DefaultCacheControl = cast.ToString(currentValue["default_cache_control"])
				}
				u[urlK] = currentURL
			}
			configuration.URLs = u
		case "ykeys":
			ykeys := make(map[string]configurationtypes.SurrogateKeys)
			d, _ := json.Marshal(v)
			_ = json.Unmarshal(d, &ykeys)
			configuration.Ykeys = ykeys
		case "disable_surrogate_key":
			configuration.SurrogateKeyDisabled = cast.ToBool(v)
		}
	}

	return configuration
}

// parseStringSlice returns the string slice corresponding to the given interface.
// The interface can be of type string which contains a comma separated list of values (e.g. foo,bar) or of type []string.
func parseStringSlice(i interface{}) []string {
	if value, ok := i.([]string); ok {
		return value
	}
	if value, ok := i.([]interface{}); ok {
		var arr []string
		for _, v := range value {
			arr = append(arr, v.(string))
		}
		return arr
	}

	if value, ok := i.(string); ok {
		if strings.HasPrefix(value, "║24║") {
			return strings.Split(strings.TrimPrefix(value, "║24║"), "║")
		}
		return strings.Split(value, ",")
	}

	if value, ok := i.([]string); ok {
		return value
	}

	return nil
}

func parseIntSlice(i interface{}) []int {
	if value, ok := i.([]int); ok {
		return value
	}
	if value, ok := i.([]interface{}); ok {
		var arr []int
		for _, v := range value {
			arr = append(arr, parseIntValue(v))
		}
		return arr
	}

	return nil
}

func parseIntValue(i interface{}) int {
	switch value := i.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case uint:
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		return int(value)
	case uint64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	default:
		return cast.ToInt(i)
	}
}

func parseUint64(i interface{}) uint64 {
	switch value := i.(type) {
	case uint64:
		return value
	case uint:
		return uint64(value)
	case uint32:
		return uint64(value)
	case uint16:
		return uint64(value)
	case uint8:
		return uint64(value)
	case int:
		return uint64(value)
	case int64:
		return uint64(value)
	case int32:
		return uint64(value)
	case float64:
		return uint64(value)
	case float32:
		return uint64(value)
	default:
		return cast.ToUint64(i)
	}
}

// New create Souin instance.
func New(_ context.Context, next http.Handler, config *TestConfiguration, name string) (http.Handler, error) {
	c := parseConfiguration(*config)

	return &SouinTraefikMiddleware{
		name:             name,
		next:             next,
		SouinBaseHandler: middleware.NewHTTPCacheHandler(&c),
	}, nil
}

func (s *SouinTraefikMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_ = s.SouinBaseHandler.ServeHTTP(rw, req, func(w http.ResponseWriter, r *http.Request) error {
		s.next.ServeHTTP(w, r)

		return nil
	})
}
