package model

import (
	"strings"
	"sync"

	"gorm.io/gorm"
)

type GameElectronic struct {
	TenantBaseEntity
	PlatformId   int32  `json:"platformId"`
	PlatformCode string `json:"platformCode"`
	GameType     int32  `json:"gameType"`
	LangKey      string `json:"langKey"`
	NameEn       string `json:"nameEn"`
	NameCn       string `json:"nameCn"`
	Sort         int32  `json:"sort"`
	ImageH5En    string `json:"imageH5En"`
	ImageH5Cn    string `json:"imageH5Cn"`
	GameCode     string `json:"gameCode"`
	GameCodeExt  string `json:"gameCodeExt"`
	Enable       bool   `json:"enable"`
	Maintain     bool   `json:"maintain"`
	GameLabel    string `json:"gameLabel"`
}

var (
	electronicTenantMap map[string][]*GameElectronic
	electronicCacheMu   sync.RWMutex
)

func (GameElectronic) TableName() string {
	return "zp_game_electronic"
}

func FindElectronicByGameCode(db *gorm.DB, tenantCode, gameCode string, platformCode string, gameType int32) (*GameElectronic, error) {
	gameElectronics, err := findByTenant(db, tenantCode)
	if err != nil {
		return nil, err
	}

	for _, gameElectronic := range gameElectronics {
		if !strings.EqualFold(gameElectronic.PlatformCode, platformCode) || !strings.EqualFold(gameElectronic.GameCode, gameCode) {
			continue
		}
		if gameType == 0 || gameElectronic.GameType == gameType {
			return gameElectronic, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func FindElectronicById(db *gorm.DB, tenantCode string, id int64) (*GameElectronic, error) {
	gameElectronics, err := findByTenant(db, tenantCode)
	if err != nil {
		return nil, err
	}

	for _, gameElectronic := range gameElectronics {
		if gameElectronic.Id == id {
			return gameElectronic, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func FindElectronicAll(db *gorm.DB) (map[string][]*GameElectronic, error) {

	var gameElectronics []GameElectronic
	err := db.Find(&gameElectronics).Error
	if err != nil {
		return nil, err
	}

	tenantMap := make(map[string][]*GameElectronic)

	for i := range gameElectronics {
		gameElectronic := &gameElectronics[i]
		tenantMap[gameElectronic.TenantCode] = append(tenantMap[gameElectronic.TenantCode], gameElectronic)
	}

	electronicCacheMu.Lock()
	electronicTenantMap = tenantMap
	electronicCacheMu.Unlock()

	return electronicTenantMap, nil
}

func findByTenant(db *gorm.DB, tenantCode string) ([]*GameElectronic, error) {
	if _, err := FindElectronicAll(db); err != nil {
		return nil, err
	}
	return electronicTenantMap[tenantCode], nil
}
