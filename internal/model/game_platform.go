package model

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"zhp-app/pkg/common"
)

// GamePlatform 是游戏平台配置对应的持久化模型。
type GamePlatform struct {
	TenantBaseEntity
	Code         string `json:"code"`
	Name         string `json:"name"`
	GameType     int8   `json:"game_type"`
	Maintain     bool   `json:"maintain"`
	LangKey      string `json:"lang_key"`
	Enable       bool   `json:"enable"`
	Hall         int    `json:"hall"`
	Sort         int    `json:"sort"`
	Image        string `json:"image"`
	GameNum      int32  `json:"game_num"`
	ShowControls bool   `json:"show_controls"`
	ShowTerminal bool   `json:"show_terminal"`
	Detail       int    `json:"detail"`
	KeyId        int64  `json:"key_id"`
	Transfer     int    `json:"transfer"`
}

func (GamePlatform) TableName() string {
	return "zp_game_platform"
}

// FindPlatformByCodeAndType 根据平台编码、游戏类型和租户查询平台配置。
func FindPlatformByCodeAndType(db *gorm.DB, code string, gameType int8, tenantCode string) (*GamePlatform, error) {
	var gamePlatform GamePlatform
	err := db.Where("code = ?", code).
		Where("game_type = ?", gameType).
		Where("tenant_code = ?", tenantCode).
		First(&gamePlatform).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &gamePlatform, nil
}

// FindPlatformAll 查询某个租户下的全部游戏平台配置。
// 这里优先读取 Redis 缓存，缓存未命中时回源 MySQL。
func FindPlatformAll(db *gorm.DB, tenantCode string) ([]*GamePlatform, error) {
	var cachedPlatforms []*GamePlatform
	redisKey := "Game:GamePlatform:" + tenantCode + ":all"
	cached, _ := common.RedisClient.Get(context.Background(), redisKey).Result()
	if cached != "" {
		if unmarshalErr := json.Unmarshal([]byte(cached), &cachedPlatforms); unmarshalErr == nil && len(cachedPlatforms) > 0 {
			return cachedPlatforms, nil
		}
	}

	var gamePlatforms []GamePlatform
	err := db.Where("tenant_code = ?", tenantCode).Find(&gamePlatforms).Error
	if err != nil {
		return nil, err
	}

	platforms := make([]*GamePlatform, 0, len(gamePlatforms))
	for i := range gamePlatforms {
		platform := gamePlatforms[i]
		platforms = append(platforms, &platform)
	}

	if payload, marshalErr := json.Marshal(platforms); marshalErr == nil {
		common.RedisClient.Set(context.Background(), redisKey, payload, 30*24*time.Hour)
	}

	return platforms, nil
}
