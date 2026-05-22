package model

import "time"

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
