package idgenx

import (
	"fmt"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/yitter/idgenerator-go/regworkerid"
	"zhp-app/pkg/config"
)

// Init 使用给定 workerId 初始化全局雪花 ID 生成器。
func Init(workerID uint16) {
	options := idgen.NewIdGeneratorOptions(workerID)
	//options.SeqBitLength = 10
	options.WorkerIdBitLength = 10
	idgen.SetIdGenerator(options)
}

// InitFromRedis 先从 Redis 注册 workerId，再用它初始化全局 ID 生成器。
func InitFromRedis(conf config.RedisConf) (uint16, error) {
	workerID := regworkerid.RegisterOne(regworkerid.RegisterConf{
		Address:         conf.Addr,
		Password:        conf.Password,
		DB:              conf.DB,
		MasterName:      conf.MasterName,
		MinWorkerId:     conf.WorkerIDMin,
		MaxWorkerId:     conf.WorkerIDMax,
		LifeTimeSeconds: conf.WorkerIDLifeSeconds,
	})
	if workerID < 0 {
		return 0, fmt.Errorf("register worker id failed: %d", workerID)
	}

	Init(uint16(workerID))
	return uint16(workerID), nil
}

// Unregister 释放当前进程在 Redis 中注册的 workerId。
func Unregister() {
	regworkerid.UnRegister()
}
