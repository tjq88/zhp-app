package model

// GameStartReq 描述开始游戏接口的请求体。
type GameStartReq struct {
	PlatformCode string `json:"platformCode" binding:"required"`
	GameType     int8   `json:"gameType" binding:"required"`
	GameId       string `json:"gameId"`
	Terminal     string `json:"terminal"`
	GameCode     string `json:"gameCode"`
	Username     string `json:"username"`
	UserId       string `json:"userId"`
	Ip           string `json:"ip"`
	TenantCode   string `json:"tenantCode"`
}

// GameStartResp 描述开始游戏接口的响应体。
type GameStartResp struct {
	URL          string `json:"url"`
	PlatformCode string `json:"platformCode"`
	PlatformName string `json:"platformName"`
	GameType     int8   `json:"gameType"`
	GameCode     string `json:"gameCode"`
}

// NewGameStartResp 根据平台和请求组装开始游戏响应。
func NewGameStartResp(platform *GamePlatform, req *GameStartReq) GameStartResp {
	return GameStartResp{
		URL:          "",
		PlatformCode: platform.Code,
		PlatformName: platform.Name,
		GameType:     req.GameType,
		GameCode:     req.GameCode,
	}
}
