package service

import (
	"errors"

	"gorm.io/gorm"
	"zhp-app/internal/model"
	"zhp-app/pkg/common"
)

type GamePlatformService struct {
	db *gorm.DB
}

var (
	// ErrGamePlatformNotFound 表示未找到对应游戏平台。
	ErrGamePlatformNotFound = errors.New("game platform not found")
	// ErrGamePlatformDisabled 表示游戏平台已被禁用。
	ErrGamePlatformDisabled = errors.New("game platform disabled")
	// ErrGamePlatformMaintaining 表示游戏平台处于维护中。
	ErrGamePlatformMaintaining = errors.New("game platform maintaining")

	ErrGameElectronicNotFound    = errors.New("game electronic not found")
	ErrGameElectronicDisabled    = errors.New("game electronic disabled")
	ErrGameElectronicMaintaining = errors.New("game electronic maintaining")

	ErrGamePlatformKeyNotFound = errors.New("game platform key not found")
)

func NewGamePlatformService() *GamePlatformService {
	return &GamePlatformService{
		db: common.Db,
	}
}

// StartGame 执行开始游戏前的平台查找和状态校验。
// 平台配置和平台密钥优先走 Redis 缓存，电子游优先走进程内缓存。
func (s *GamePlatformService) StartGame(req *model.GameStartReq) (*model.GameStartResp, error) {
	gamePlatform, err := s.loadAvailablePlatform(req)
	if err != nil {
		return nil, err
	}

	gameElectronic, err := s.loadAvailableElectronic(req)
	if err != nil {
		return nil, err
	}

	gamePlatformKey, err := s.loadPlatformKey(req.TenantCode, gamePlatform.KeyId)
	if err != nil {
		return nil, err
	}

	s.attachGameContext(req, gamePlatform, gameElectronic, gamePlatformKey)

	resp := model.NewGameStartResp(gamePlatform, req)
	return &resp, nil
}

// FindPlatformByCodeAndType 按租户、平台编码和游戏类型查找游戏平台。
func (s *GamePlatformService) FindPlatformByCodeAndType(req *model.GameStartReq) (*model.GamePlatform, error) {
	gamePlatform, err := model.FindPlatformCacheByCodeAndType(s.db, req.PlatformCode, req.GameType, req.TenantCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGamePlatformNotFound
		}
		return nil, err
	}

	return gamePlatform, nil
}

func (s *GamePlatformService) loadAvailablePlatform(req *model.GameStartReq) (*model.GamePlatform, error) {
	gamePlatform, err := s.FindPlatformByCodeAndType(req)
	if err != nil {
		return nil, err
	}

	if !gamePlatform.Enable {
		return nil, ErrGamePlatformDisabled
	}
	if gamePlatform.Maintain {
		return nil, ErrGamePlatformMaintaining
	}

	return gamePlatform, nil
}

func (s *GamePlatformService) loadAvailableElectronic(req *model.GameStartReq) (*model.GameElectronic, error) {
	gameElectronic, err := model.FindElectronicById(s.db, req.TenantCode, req.GameId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGameElectronicNotFound
		}
		return nil, err
	}

	if !gameElectronic.Enable {
		return nil, ErrGameElectronicDisabled
	}
	if gameElectronic.Maintain {
		return nil, ErrGameElectronicMaintaining
	}

	return gameElectronic, nil
}

func (s *GamePlatformService) loadPlatformKey(tenantCode string, keyID int64) (*model.GamePlatformKey, error) {
	gamePlatformKey, err := model.FindPlatformKeyCacheById(s.db, tenantCode, keyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGamePlatformKeyNotFound
		}
		return nil, err
	}

	return gamePlatformKey, nil
}

func (s *GamePlatformService) attachGameContext(
	req *model.GameStartReq,
	gamePlatform *model.GamePlatform,
	gameElectronic *model.GameElectronic,
	gamePlatformKey *model.GamePlatformKey,
) {
	req.Platform = gamePlatform
	req.Electronic = gameElectronic
	req.PlatformKey = gamePlatformKey
}
