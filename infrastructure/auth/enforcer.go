package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/rs/zerolog/log"
)

// AuthEnforcer defines the interface for authorization enforcement.
type AuthEnforcer interface {
	Enforce(sub string, obj AuthObject, act AuthAction) (bool, error)
	AddPolicy(sub string, obj AuthObject, act AuthAction) error
	RemovePolicy(sub string, obj AuthObject, act AuthAction) error
}

// authEnforcer implements AuthEnforcer using Casbin.
type authEnforcer struct {
	enforcer *casbin.Enforcer
}

// AuthEnforcerOpts contains configuration options for AuthEnforcer.
type AuthEnforcerOpts struct {
	ModelPath string
}

// DefaultAuthEnforcerOpts provides default configuration for AuthEnforcer.
var DefaultAuthEnforcerOpts = &AuthEnforcerOpts{
	ModelPath: "auth_model.conf",
}

// NewAuthEnforcer creates a new instance of AuthEnforcer.
func NewAuthEnforcer(adapter persist.Adapter, opts ...*AuthEnforcerOpts) *authEnforcer {
	option := DefaultAuthEnforcerOpts
	if len(opts) > 0 {
		option = opts[0]
	}

	enforcer, err := casbin.NewEnforcer(option.ModelPath, adapter)
	if err != nil {
		log.Fatal().Err(err).Msg("💣 failed to create enforcer")
	}

	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatal().Err(err).Msg("💣 failed to load policy")
	}

	return &authEnforcer{enforcer}
}

// Enforce checks if a request is allowed.
func (e *authEnforcer) Enforce(sub string, obj AuthObject, act AuthAction) (bool, error) {
	return e.enforcer.Enforce(sub, obj.String(), act.String())
}

// AddPolicy adds a new policy rule.
func (e *authEnforcer) AddPolicy(sub string, obj AuthObject, act AuthAction) error {
	if ok, err := e.enforcer.AddPolicy(sub, obj.String(), act.String()); !ok || err != nil {
		log.Error().Err(err).Msg("🚨 failed to add policy")
		return fmt.Errorf("failed to add policy")
	}

	if err := e.enforcer.SavePolicy(); err != nil {
		log.Error().Err(err).Msg("🚨 failed to save policy")
		return fmt.Errorf("failed to save policy")
	}

	return nil
}

// RemovePolicy removes a policy rule.
func (e *authEnforcer) RemovePolicy(sub string, obj AuthObject, act AuthAction) error {
	if ok, err := e.enforcer.RemovePolicy(sub, obj.String(), act.String()); !ok || err != nil {
		log.Error().Err(err).Msg("🚨 failed to remove policy")
		return fmt.Errorf("failed to remove policy")
	}

	if err := e.enforcer.SavePolicy(); err != nil {
		log.Error().Err(err).Msg("🚨 failed to save policy")
		return fmt.Errorf("failed to save policy")
	}

	return nil
}
