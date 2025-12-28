package handler

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"development-of-internet-applications/internal/app/role"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const jwtPrefix = "Bearer "

type BaseHandler struct {
	repo *repository.Repository
}

func NewBaseHandler(repo *repository.Repository) *BaseHandler {
	return &BaseHandler{
		repo: repo,
	}
}

func (h *BaseHandler) WithAuthCheck(assignedRoles ...role.Role) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]

		if h.repo.Redis != nil {
			err := h.repo.Redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
			if err == nil {
				gCtx.AbortWithStatus(http.StatusForbidden)
				return
			}
			if !errors.Is(err, redis.Nil) {
				logrus.Error("Redis error:", err)
				gCtx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.repo.Config.JWT.Token), nil
		})
		if err != nil {
			logrus.Error("JWT parse error:", err)
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !token.Valid {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		myClaims, ok := token.Claims.(*ds.JWTClaims)
		if !ok {
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !slices.Contains(assignedRoles, myClaims.Role) {
			logrus.Warnf("Role %s is not in assigned roles %v", myClaims.Role, assignedRoles)
			gCtx.AbortWithStatus(http.StatusForbidden)
			return
		}

		gCtx.Set("user_id", myClaims.UserID)
		gCtx.Set("user_role", myClaims.Role)

		gCtx.Next()
	}
}

func GetUserIDFromContext(gCtx *gin.Context) (uint64, bool) {
	userID, exists := gCtx.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint64)
	return id, ok
}

func GetUserRoleFromContext(gCtx *gin.Context) (role.Role, bool) {
    if val, exists := gCtx.Get("user_role"); exists {
        logrus.Debugf("Found 'user_role': %v (type: %T)", val, val)
        return convertToRole(val)
    }
    
    if val, exists := gCtx.Get("role"); exists {
        logrus.Debugf("Found 'role': %v (type: %T)", val, val)
        return convertToRole(val)
    }
    
    return 0, false
}

func convertToRole(val interface{}) (role.Role, bool) {
    switch v := val.(type) {
    case role.Role:
        return v, true
    case float64:
        return role.Role(v), true
    case int:
        return role.Role(v), true
    case string:
        if v == "moderator" || v == "1" {
            return role.Moderator, true
        } else if v == "user" || v == "0" {
            return role.User, true
        }
    }
    return 0, false
}