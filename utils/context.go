package utils

import (
	"fmt"

	"gorm.io/gorm"
)

type contextKey string

const (
	userIDKey   contextKey = "userID"
	userRoleKey contextKey = "userRole"
)

type ContextKeys struct {
	UserIDKey   string `json:"user_id_key"`
	UserRoleKey string `json:"user_role_key"`
}

func GetContextKeys() ContextKeys {
	return ContextKeys{
		UserIDKey:   string(userIDKey),
		UserRoleKey: string(userRoleKey),
	}
}

func GetCurrentUserID(tx *gorm.DB) (uint, error) {

	if userID, ok := tx.Statement.Context.Value(GetContextKeys().UserIDKey).(uint); ok {
		return userID, nil
	}
	return 0, fmt.Errorf("userID not found in context for transaction: %v", tx)
}
