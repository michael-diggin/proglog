package auth

import (
	"fmt"

	"github.com/casbin/casbin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authorizer is the type to store and enforce ACL rules
type Authorizer struct {
	enforcer *casbin.Enforcer
}

// New returns a new Authorizer instance
func New(model, policy string) *Authorizer {
	enforcer := casbin.NewEnforcer(model, policy)
	return &Authorizer{enforcer}
}

// Authorize determines if the given subject can act on an object
func (a *Authorizer) Authorize(subject, object, action string) error {
	if !a.enforcer.Enforce(subject, object, action) {
		msg := fmt.Sprintf("%s not permitted to %s to %s", subject, object, action)
		st := status.New(codes.PermissionDenied, msg)
		return st.Err()
	}
	return nil
}
