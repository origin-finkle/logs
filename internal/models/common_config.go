package models

import (
	"context"

	"github.com/origin-finkle/logs/internal/logger"
)

type CommonConfig struct {
	Invalid                      bool                `json:"invalid,omitempty"`
	InvalidReason                string              `json:"invalid_reason,omitempty"`
	RestrictedRoles              []string            `json:"restricted_roles,omitempty"`
	RestrictedSpecializations    []string            `json:"restricted_specializations,omitempty"`
	RestrictedSpecializationsNot []string            `json:"restricted_specializations_not,omitempty"`
	RestrictedFights             []string            `json:"restricted_fights,omitempty"`
	RestrictedAny                []string            `json:"restricted_any,omitempty"`
	RestrictedClass              []string            `json:"restricted_class,omitempty"`
	RestrictedComplex            *ComplexRestriction `json:"restricted_complex"`

	Todo string `json:"__todo,omitempty"`
}

func (cc CommonConfig) IsRestricted(ctx context.Context, fa *FightAnalysis) bool {
	if cc.Invalid {
		return true
	}
	if cc.RestrictedComplex != nil {
		return cc.RestrictedComplex.IsRestricted(ctx, fa)
	}
	if len(cc.RestrictedAny) > 0 {
		valid := false
		for _, restricted := range cc.RestrictedAny {
			switch true {
			case stringIsClass(restricted):
				valid = valid || fa.player.SubType == restricted
			case stringIsRole(restricted):
				valid = valid || fa.Talents.Spec.IsRole(Role(restricted))
			case stringIsSpecialization(restricted):
				valid = valid || fa.Talents.Spec == Specialization(restricted)
			}
		}
		if !valid {
			logger.FromContext(ctx).Debugf("player is not in %v", cc.RestrictedAny)
			return true
		}
	}
	if len(cc.RestrictedFights) > 0 {
		if !in(fa.Name, cc.RestrictedFights) {
			logger.FromContext(ctx).Debugf("fight %s is not in %v", fa.Name, cc.RestrictedFights)
			return true
		}
	}
	if len(cc.RestrictedRoles) > 0 {
		restricted := true
		for _, role := range cc.RestrictedRoles {
			logger.FromContext(ctx).Debugf("checking if player is %s", role)
			if fa.Talents.Spec.IsRole(Role(role)) {
				logger.FromContext(ctx).Debugf("player is %s", role)
				restricted = false
				break
			}
		}
		if restricted {
			logger.FromContext(ctx).Debugf("player spec %s does not have role in %v", fa.Talents.Spec, cc.RestrictedRoles)
			return true
		}
	}
	if len(cc.RestrictedSpecializations) > 0 {
		if !in(string(fa.Talents.Spec), cc.RestrictedSpecializations) {
			logger.FromContext(ctx).Debugf("player spec %s is not in %v", fa.Talents.Spec, cc.RestrictedSpecializations)
			return true
		}
	}
	if len(cc.RestrictedSpecializationsNot) > 0 {
		if in(string(fa.Talents.Spec), cc.RestrictedSpecializationsNot) {
			logger.FromContext(ctx).Debugf("player spec %s is in %v", fa.Talents.Spec, cc.RestrictedSpecializationsNot)
			return true
		}
	}
	if len(cc.RestrictedClass) > 0 {
		if !in(fa.player.SubType, cc.RestrictedClass) {
			logger.FromContext(ctx).Debugf("player class %s is not in %v", fa.player.SubType, cc.RestrictedClass)
			return true
		}
	}
	return false
}

func in(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

type ComplexRestriction struct {
	Operator string       `json:"operator"`
	LHS      CommonConfig `json:"lhs"`
	RHS      CommonConfig `json:"rhs"`
}

func (cr *ComplexRestriction) IsRestricted(ctx context.Context, fa *FightAnalysis) bool {
	switch cr.Operator {
	case "AND":
		lhs := cr.LHS.IsRestricted(ctx, fa)
		rhs := cr.RHS.IsRestricted(ctx, fa)
		if !lhs && !rhs {
			return false
		}
		return true
	case "OR":
		return cr.LHS.IsRestricted(ctx, fa) && cr.RHS.IsRestricted(ctx, fa)
	}
	return false
}
