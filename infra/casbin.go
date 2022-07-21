package infra

import "github.com/casbin/casbin/v2"

type CasbinConfig struct {
	ModelPath  string
	PolicyPath string
}

func NewCasbin(cfg CasbinConfig) (*casbin.Enforcer, error) {
	ce, err := casbin.NewEnforcer(cfg.ModelPath, cfg.PolicyPath)
	if err != nil {
		return nil, err
	}
	return ce, nil
}
