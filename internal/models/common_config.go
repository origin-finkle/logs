package models

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/origin-finkle/logs/internal/logger"
)

type CommonConfig struct {
	Invalid                      bool                `json:"invalid,omitempty"`
	InvalidReason                string              `json:"invalid_reason,omitempty"`
	RestrictedRoles              []string            `json:"restricted_roles,omitempty" stringer:"PlayerRole"`
	RestrictedSpecializations    []string            `json:"restricted_specializations,omitempty" stringer:"PlayerSpec"`
	RestrictedSpecializationsNot []string            `json:"restricted_specializations_not,omitempty" stringer:"PlayerSpec,not"`
	RestrictedFights             []string            `json:"restricted_fights,omitempty" stringer:"Fight"`
	RestrictedAny                []string            `json:"restricted_any,omitempty" stringer:"Player"`
	RestrictedClass              []string            `json:"restricted_class,omitempty" stringer:"PlayerClass"`
	RestrictedComplex            *ComplexRestriction `json:"restricted_complex,omitempty"`

	Todo     string `json:"__todo,omitempty"`
	TextRule string `json:"text_rule,omitempty"`
}

func (cc CommonConfig) String() string {
	if cc.Invalid {
		return "INVALID"
	}
	if cc.RestrictedComplex != nil {
		return cc.RestrictedComplex.String()
	}
	parts := make([]string, 0)
	t := reflect.TypeOf(cc)
	v := reflect.ValueOf(cc)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if str := field.Tag.Get("stringer"); str != "" {
			fieldVal := v.Field(i)
			if fieldVal.Len() == 0 {
				continue
			}
			ss := strings.Split(str, ",")
			var equality bool
			if len(ss) == 2 && ss[1] == "not" {
				equality = false
			} else {
				equality = true
			}
			var operator string
			values := strings.Join(fieldVal.Interface().([]string), ", ")
			if fieldVal.Len() == 1 {
				if equality {
					operator = "="
				} else {
					operator = "!="
				}
			} else {
				operator = "IN"
				if !equality {
					operator = "NOT " + operator
				}
				values = "(" + values + ")"
			}
			parts = append(parts, fmt.Sprintf("%s %s %s", ss[0], operator, values))
		}
	}
	return strings.Join(parts, " AND ")
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
		if !in(fa.FightName(), cc.RestrictedFights) {
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

func (cr *ComplexRestriction) String() string {
	return fmt.Sprintf("(%s) %s (%s)", cr.LHS.String(), cr.Operator, cr.RHS.String())
}
