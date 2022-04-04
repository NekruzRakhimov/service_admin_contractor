package middleware

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"service_admin_contractor/application/respond"
	"service_admin_contractor/domain/model"
	"service_admin_contractor/domain/repository"
)

const (
	UserInfoCtxKey = "UserInfo"
)

type authenticationHandler struct {
	r    repository.BpmsUserRepository
	next http.Handler
}

func newAuthenticationHandler(r repository.BpmsUserRepository, next http.Handler) *authenticationHandler {
	return &authenticationHandler{r, next}
}

func (a *authenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var userInfo *model.UserInfo = nil

	if login, _, ok := r.BasicAuth(); ok {
		roles, err := a.r.FindUserRoles(login)
		if err != nil {
			log.Warn(err)
		}
		userRoles := make([]model.RoleCode, len(roles))
		for i, role := range roles {
			userRoles[i] = model.RoleCode(role)
		}
		userInfo = model.NewUserInfo(r.Header.Get("Authorization"), login, userRoles)
	} else {
		respond.WithError(w, r, errors.New("authentication failed"))
		return
	}

	ctx := context.WithValue(r.Context(), UserInfoCtxKey, userInfo)

	a.next.ServeHTTP(w, r.WithContext(ctx))
}

func AuthHandler(r repository.BpmsUserRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return newAuthenticationHandler(r, next)
	}
}

func GetUserInfo(ctx context.Context) *model.UserInfo {
	info, ok := ctx.Value(UserInfoCtxKey).(*model.UserInfo)
	if !ok {
		return nil
	}

	return info
}
