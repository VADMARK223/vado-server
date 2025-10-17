package auth

import "context"

type AuthContext struct {
	context.Context
	userID uint
}

func (a *AuthContext) UserID() uint {
	return a.userID
}

func Wrap(ctx context.Context, userID uint) context.Context {
	return &AuthContext{
		Context: ctx,
		userID:  userID,
	}
}

func TryGet(ctx context.Context) (uint, bool) {
	if c, ok := ctx.(*AuthContext); ok {
		return c.userID, true
	}
	return 0, false
}
