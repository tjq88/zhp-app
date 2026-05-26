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
func (s *GamePlatformService) StartGame(req *model.GameStartReq) (*model.GameStartResp, error) {
	//gamePlatform
	gamePlatform, err := s.FindPlatformByCodeAndType(req)
	if err != nil {
		return nil, err
	}
	if gamePlatform == nil {
		return nil, ErrGamePlatformNotFound
	}
	if !gamePlatform.Enable {
		return nil, ErrGamePlatformDisabled
	}
	if gamePlatform.Maintain {
		return nil, ErrGamePlatformMaintaining
	}

	//electronic
	gameElectronic, err := model.FindElectronicById(s.db, req.TenantCode, req.GameId)
	if err != nil {
		return nil, err
	}
	if gameElectronic == nil {
		return nil, ErrGameElectronicNotFound
	}
	if gameElectronic.Enable {
		return nil, ErrGameElectronicDisabled
	}
	if gameElectronic.Maintain {
		return nil, ErrGameElectronicMaintaining
	}
	//
	gamePlatformKey, err := model.FindPlatformKeyById(s.db, req.TenantCode, gamePlatform.KeyId)
	if err != nil {
		return nil, err
	}
	if gamePlatformKey == nil {
		return nil, ErrGamePlatformKeyNotFound
	}
	req.Electronic = gameElectronic
	req.Platform = gamePlatform
	req.PlatformKey = gamePlatformKey

	resp := model.NewGameStartResp(gamePlatform, req)
	return &resp, nil
}

// FindPlatformByCodeAndType 按租户、平台编码和游戏类型查找游戏平台。
func (s *GamePlatformService) FindPlatformByCodeAndType(req *model.GameStartReq) (*model.GamePlatform, error) {
	gamePlatforms, err := s.FindAll(req.TenantCode)
	if err != nil {
		return nil, err
	}
	for _, gamePlatform := range gamePlatforms {
		if gamePlatform.Code == req.PlatformCode && gamePlatform.GameType == req.GameType {
			return gamePlatform, nil
		}
	}

	gamePlatform, err := model.FindPlatformByCodeAndType(s.db, req.PlatformCode, req.GameType, req.TenantCode)
	if err != nil {
		return nil, err
	}
	return gamePlatform, nil
}

// FindAll 查询当前租户下的全部游戏平台配置。
func (s *GamePlatformService) FindAll(tenantCode string) ([]*model.GamePlatform, error) {
	return model.FindPlatformAll(s.db, tenantCode)
}
