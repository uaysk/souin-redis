package storage

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/uaysk/souin-redis/configurationtypes"
	"github.com/uaysk/souin-redis/pkg/storage/types"
)

const (
	VarySeparator                   = "{-VARY-}"
	DecodedHeaderSeparator          = ";"
	encodedHeaderSemiColonSeparator = "%3B"
	encodedHeaderColonSeparator     = "%3A"
	StalePrefix                     = "STALE_"
)

type StorerInstanciator func(configurationtypes.AbstractConfigurationInterface) (types.Storer, error)

func NewStorages(configuration configurationtypes.AbstractConfigurationInterface) ([]types.Storer, error) {
	if shouldUseRedis(configuration) {
		if redis, err := newRedisStorage(configuration); err == nil {
			return []types.Storer{redis}, nil
		} else {
			return nil, err
		}
	}

	s, err := Factory(configuration)
	return []types.Storer{s}, err
}

func shouldUseRedis(configuration configurationtypes.AbstractConfigurationInterface) bool {
	for _, storer := range configuration.GetDefaultCache().GetStorers() {
		if strings.EqualFold(storer, "redis") {
			return true
		}
	}

	redis := configuration.GetDefaultCache().GetRedis()
	return redis.URL != "" || redis.Path != "" || redis.Configuration != nil
}

func varyVoter(baseKey string, req *http.Request, currentKey string) bool {
	if currentKey == baseKey {
		return true
	}

	if strings.Contains(currentKey, VarySeparator) && strings.HasPrefix(currentKey, baseKey+VarySeparator) {
		list := currentKey[(strings.LastIndex(currentKey, VarySeparator) + len(VarySeparator)):]
		if len(list) == 0 {
			return false
		}

		for _, item := range strings.Split(list, ";") {
			index := strings.LastIndex(item, ":")
			if len(item) < index+1 {
				return false
			}

			hVal := item[index+1:]
			if strings.Contains(hVal, encodedHeaderSemiColonSeparator) || strings.Contains(hVal, encodedHeaderColonSeparator) {
				hVal, _ = url.QueryUnescape(hVal)
			}
			if req.Header.Get(item[:index]) != hVal {
				return false
			}
		}

		return true
	}

	return false
}
