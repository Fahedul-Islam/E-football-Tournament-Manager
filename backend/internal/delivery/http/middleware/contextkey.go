package middleware

// ContextKey is the type for HTTP request context keys.
// Using a named type prevents key collisions with other packages.
type ContextKey string

const (
	ContextKeyUserID ContextKey = "user_id"
	ContextKeyEmail  ContextKey = "email"
	ContextKeyRole   ContextKey = "roles"
)
