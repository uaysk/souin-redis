package traefik

import "testing"

func TestParseConfigurationFromMantraeStyleConfig(t *testing.T) {
	cfg := map[string]interface{}{
		"log_level":              "INFO",
		"disable_surrogate_key":  true,
		"default_cache": map[string]interface{}{
			"cache_name": "homelab-public",
			"ttl":        "10m",
			"stale":      "2m",
			"default_cache_control": "public, max-age=600",
			"allowed_http_verbs": []interface{}{"GET", "HEAD"},
			"allowed_additional_status_codes": []interface{}{float64(301), float64(308)},
			"headers": []interface{}{"Accept-Language"},
			"storers": []interface{}{"redis"},
			"regex": map[string]interface{}{
				"exclude": "(?i)(^/(?:api|graphql|rpc)(?:/|$)|^/(?:admin|login|logout|signin|signout|auth|oauth|callback|session|account|profile|settings|user)(?:/|$)|^/(?:wp-admin|wp-login\\.php|wp-json|xmlrpc\\.php)(?:/|$)|^/(?:preview|socket|ws|websocket|stream|events)(?:/|$))",
			},
			"timeout": map[string]interface{}{
				"backend": "15s",
				"cache":   "100ms",
			},
			"redis": map[string]interface{}{
				"configuration": map[string]interface{}{
					"InitAddress": []interface{}{"redis.redis.svc.cluster.local:6379"},
					"ClientName":  "traefik-souin",
					"DialTimeout": "1s",
					"SelectDB":    float64(0),
				},
			},
		},
		"urls": map[string]interface{}{
			"^blog\\.uaysk\\.com(?:/.*)?$": map[string]interface{}{
				"ttl": "15m",
			},
			"^blog\\.uaysk\\.com/.+\\.(?:css|js|mjs|map|png|jpe?g|gif|svg|webp|avif|ico|woff2?|ttf|eot)$": map[string]interface{}{
				"ttl": "24h",
			},
		},
	}

	parsed := parseConfiguration(cfg)
	if parsed.GetDefaultCache() == nil {
		t.Fatal("default cache was nil")
	}
}
