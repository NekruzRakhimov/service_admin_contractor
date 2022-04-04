package model

import "strings"

type RoleCode string

type UserInfo struct {
	basicAuth string
	login     string
	roles     []RoleCode
}

func (u UserInfo) Login() string {
	return u.login
}

func (u UserInfo) Roles() []RoleCode {
	return u.roles
}

func (u UserInfo) BasicAuth() string {
	return u.basicAuth
}

func NewUserInfo(basicAuth string, login string, roles []RoleCode) *UserInfo {
	return &UserInfo{
		basicAuth: basicAuth,
		login:     strings.ToUpper(login),
		roles:     roles,
	}
}
