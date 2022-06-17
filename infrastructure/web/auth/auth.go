package auth

import (
	"context"
)

const (
	UserIdKey string = "user_id"
)

func GetUserId(ctx context.Context) string {
	return ctx.Value(UserIdKey).(string)
}

func WrapWithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}
