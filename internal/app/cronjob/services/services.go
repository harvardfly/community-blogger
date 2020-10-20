package services

import (
	"community-blogger/internal/pkg/transports/cron"
	"github.com/google/wire"
)

// CreateInitServersFn CronJob服务入口
func CreateInitServersFn(cronJob *DefaultCronJobService) cron.InitServers {
	return map[string]func(){
		"sync.es.total.per.5s": func() {
			go cronJob.RedisToES()
		},
	}
}

// ProviderSet CronJob Service wire 注入
var ProviderSet = wire.NewSet(NewDefaultCronJobService, CreateInitServersFn)
