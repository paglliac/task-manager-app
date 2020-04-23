package auth

import (
	"context"
	"net/http"
)

var credentialsKey = "credentials"

// FromRequest returns the Credentials value stored in ctx, if any.
func FromRequest(r *http.Request) (*Credentials, bool) {
	ctx := r.Context()
	u, ok := ctx.Value(credentialsKey).(*Credentials)
	return u, ok
}

func NewContextWithCredentials(ctx context.Context, c *Credentials) context.Context {
	return context.WithValue(ctx, credentialsKey, c)
}
