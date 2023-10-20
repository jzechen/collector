/**
 * @Time: 2023/10/20 16:19
 * @Author: jzechen
 * @File: rate.go
 * @Software: GoLand collector
 */

package rate

import (
	"github.com/jzechen/collector/pkg/manager/config"
	"golang.org/x/time/rate"
)

var limit *rate.Limiter

func InitRateLimiter(cfg *config.ServerConfig) {
	limit = rate.NewLimiter(rate.Limit(cfg.Limit), cfg.Burst)
}

func GetRateLimiter() *rate.Limiter {
	if limit == nil {
		limit = rate.NewLimiter(10, 100)
	}
	return limit
}
