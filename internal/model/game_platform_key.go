package model

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"strings"
	"time"
	"zhp-app/pkg/common"
)

type GamePlatformKey struct {
	TenantBaseEntity
	Code    string `json:"code"`
	Name    string `json:"name"`
	KeyInfo string `json:"key_info"`
	Remark  string `json:"remark"`
}

func (GamePlatformKey) TableName() string {
	return "zp_game_platform_key"
}

func FindPlatformKeyCacheById(db *gorm.DB, tenantCode string, id int64) (*GamePlatformKey, error) {
	gamePlatformKeys, err := FindPlatformKeyCacheAll(db, tenantCode)
	if err != nil {
		return nil, err
	}
	for _, gamePlatformKey := range gamePlatformKeys {
		if gamePlatformKey.Id == id {
			return gamePlatformKey, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func FindPlatformKeyCacheByCode(db *gorm.DB, tenantCode string, code string) (*GamePlatformKey, error) {
	gamePlatformKeys, err := FindPlatformKeyCacheAll(db, tenantCode)
	if err != nil {
		return nil, err
	}
	for _, gamePlatformKey := range gamePlatformKeys {
		if strings.EqualFold(gamePlatformKey.Code, code) {
			return gamePlatformKey, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func FindPlatformKeyCacheAll(db *gorm.DB, tenantCode string) ([]*GamePlatformKey, error) {
	var cachePlatformKeys []*GamePlatformKey
	redisKey := "Game:GamePlatformKey:" + tenantCode + ":all"
	cached, _ := common.RedisClient.Get(context.Background(), redisKey).Result()
	if cached != "" {
		if err := json.Unmarshal([]byte(cached), &cachePlatformKeys); err == nil && len(cachePlatformKeys) > 0 {
			return cachePlatformKeys, nil
		}
	}

	var gamePlatformKeys []GamePlatformKey
	err := db.Where("tenant_code = ?", tenantCode).Find(&gamePlatformKeys).Error
	if err != nil {
		return nil, err
	}

	platformKeys := make([]*GamePlatformKey, 0, len(gamePlatformKeys))
	for i := range gamePlatformKeys {
		platformKey := gamePlatformKeys[i]
		platformKeys = append(platformKeys, &platformKey)
	}
	payload, marshalErr := json.Marshal(platformKeys)
	if marshalErr == nil {
		common.RedisClient.Set(context.Background(), redisKey, string(payload), 30*24*time.Hour)
	}
	return platformKeys, nil
}
