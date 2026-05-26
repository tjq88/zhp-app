package model

import (
	"strings"
	"sync"
	"sync/atomic"

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
	electronicCacheHolderValue atomic.Pointer[electronicCacheHolder]
	electronicCacheResetMu     sync.Mutex
)

type electronicCacheHolder struct {
	once      sync.Once
	tenantMap map[string][]*GameElectronic
	err       error
}

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
	holder := getElectronicCacheHolder()
	holder.once.Do(func() {
		var gameElectronics []GameElectronic
		holder.err = db.Find(&gameElectronics).Error
		if holder.err != nil {
			return
		}

		tenantMap := make(map[string][]*GameElectronic)
		for i := range gameElectronics {
			gameElectronic := &gameElectronics[i]
			tenantKey := normalizeTenantCode(gameElectronic.TenantCode)
			tenantMap[tenantKey] = append(tenantMap[tenantKey], gameElectronic)
		}
		holder.tenantMap = tenantMap
	})

	return holder.tenantMap, holder.err
}

func findByTenant(db *gorm.DB, tenantCode string) ([]*GameElectronic, error) {
	if _, err := FindElectronicAll(db); err != nil {
		return nil, err
	}

	holder := getElectronicCacheHolder()
	return holder.tenantMap[normalizeTenantCode(tenantCode)], nil
}

// ClearElectronicCache 清空电子游进程内缓存，供后台刷新后重新加载使用。
func ClearElectronicCache() {
	electronicCacheResetMu.Lock()
	defer electronicCacheResetMu.Unlock()

	electronicCacheHolderValue.Store(&electronicCacheHolder{})
}

// ReloadElectronicCache 清空电子游进程内缓存并立即重新从数据库加载。
func ReloadElectronicCache(db *gorm.DB) (map[string][]*GameElectronic, error) {
	ClearElectronicCache()
	return FindElectronicAll(db)
}

func normalizeTenantCode(tenantCode string) string {
	return strings.ToLower(tenantCode)
}

func getElectronicCacheHolder() *electronicCacheHolder {
	holder := electronicCacheHolderValue.Load()
	if holder != nil {
		return holder
	}

	newHolder := &electronicCacheHolder{}
	if electronicCacheHolderValue.CompareAndSwap(nil, newHolder) {
		return newHolder
	}

	return electronicCacheHolderValue.Load()
}
