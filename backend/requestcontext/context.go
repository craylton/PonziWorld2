package requestcontext

import "context"

// ctxKey is a private type for context keys in this package
type ctxKey string

const usernameKey ctxKey = "username"

// WithUsername returns a new Context that carries the provided username value.
func WithUsername(ctx context.Context, username string) context.Context {
    return context.WithValue(ctx, usernameKey, username)
}

// UsernameFromContext retrieves the username stored in the Context.
func UsernameFromContext(ctx context.Context) (string, bool) {
    username, ok := ctx.Value(usernameKey).(string)
    return username, ok
}
