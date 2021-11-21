package models

import (
	"sync"

	"github.com/origin-finkle/logs/internal/models/remark"
)

type FightAnalysis struct {
	mu     sync.Mutex      `json:"-"`
	player *PlayerAnalysis `json:"-"`

	Name     string             `json:"name"`
	Auras    map[int64]*Aura    `json:"auras"`
	Remarks  []*remark.Remark   `json:"remarks"`
	Talents  *Talents           `json:"talents"`
	Casts    map[int64]int64    `json:"casts"`
	Gear     []*Gear            `json:"gear"`
	Analysis *TrueFightAnalysis `json:"analysis"`
}

func (fa *FightAnalysis) AddCast(ability, count int64) {
	fa.Casts[ability] += 1
}

func (fa *FightAnalysis) AddRemark(metadata remark.Metadata) {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	fa.Remarks = append(fa.Remarks, remark.New(metadata, fa.Name))
}

func (fa *FightAnalysis) CouldBenefitFromWindfury(analysis *Analysis) bool {
	if fa.Talents.BenefitsFromWindfuryTotem() {
		shamanPlayers := analysis.FilterPlayers(func(pa *PlayerAnalysis) bool {
			return pa.SubType == string(Class_Shaman) && pa.GetFight(fa.Name) != nil
		})
		return len(shamanPlayers) > 0
	}
	return false
}