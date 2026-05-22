package model

import "time"

// MemberView 是返回给调用方的安全响应对象。
// 像 Password 这种仅供持久化使用的敏感字段不能出现在这里。
type MemberView struct {
	Id            int64     `json:"id"`
	Username      string    `json:"username"`
	TenantCode    string    `json:"tenantCode"`
	CreateUser    string    `json:"createUser"`
	UpdateUser    string    `json:"updateUser"`
	CreateTime    time.Time `json:"createTime"`
	UpdateTime    time.Time `json:"updateTime"`
	LastLoginTime time.Time `json:"lastLoginTime"`
}

// NewMemberView 把持久化模型转换成接口响应模型。
func NewMemberView(member *MemberInfo) MemberView {
	return MemberView{
		Id:            member.Id,
		Username:      member.Username,
		TenantCode:    member.TenantCode,
		CreateUser:    member.CreateUser,
		UpdateUser:    member.UpdateUser,
		CreateTime:    member.CreateTime,
		UpdateTime:    member.UpdateTime,
		LastLoginTime: member.LastLoginTime,
	}
}
