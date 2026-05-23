package model

type GamePlatformKey struct {
	TenantBaseEntity
	Code    string `json:"code"`
	Name    string `json:"name"`
	KeyInfo string `json:"key_info"`
	Remark  string `json:"remark"`
}
