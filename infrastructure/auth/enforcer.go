package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/rs/zerolog/log"
)

type AuthEnforcer interface {
	Enforce(sub string, obj AuthObject, act AuthAction) (bool, error)
	AddPolicy(sub string, obj AuthObject, act AuthAction) error
	RemovePolicy(sub string, obj AuthObject, act AuthAction) error
}

type authEnforcer struct {
	enforcer *casbin.Enforcer
}

type AuthEnforcerOpts struct {
	ModelPath string
}

var DefaultAuthEnforcerOpts = &AuthEnforcerOpts{
	ModelPath: "auth_model.conf",
}

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

func (e *authEnforcer) Enforce(sub string, obj AuthObject, act AuthAction) (bool, error) {
	return e.enforcer.Enforce(sub, obj.String(), act.String())
}

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
