package cache

import (
	"fmt"

	"github.com/pilillo/igovium/utils"
)

// todo: convert to a map
func NewDMCacheFromConfig(config *utils.DMCacheConfig) (DMCache, error) {
	switch dmType := config.Type; dmType {
	case "olric":
		return NewOlricDMCache(), nil
	case "redis":
		return NewRedisDMCache(), nil
	default:
		return nil, fmt.Errorf("Unknown dm type %s", dmType)
	}
}
