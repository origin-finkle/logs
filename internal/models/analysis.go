package models

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/origin-finkle/logs/internal/models/remark"
)

type Analysis struct {
	mu   sync.RWMutex
	Data map[int64]*PlayerAnalysis
}

func (a *Analysis) FilterPlayers(predicate func(*PlayerAnalysis) bool) []*PlayerAnalysis {
	a.mu.RLock()
	defer a.mu.RUnlock()

	pas := make([]*PlayerAnalysis, 0)
	for _, pa := range a.Data {
		if predicate(pa) {
			pas = append(pas, pa)
		}
	}
	return pas
}

func (a *Analysis) GetPlayerAnalysis(gameID int64) *PlayerAnalysis {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.Data[gameID]
}

func (a *Analysis) SetPlayerAnalysis(playerID int64, pa *PlayerAnalysis) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Data[playerID] = pa
}

func (a *Analysis) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Data)
}

type PlayerAnalysis struct {
	Actor
	mu      sync.RWMutex              `json:"-"`
	Fights  map[string]*FightAnalysis `json:"fights"`
	Remarks []*remark.Remark          `json:"remarks"`
}

func (pa *PlayerAnalysis) GetFight(name string) *FightAnalysis {
	pa.mu.RLock()
	defer pa.mu.RUnlock()

	return pa.Fights[name]
}

func (pa *PlayerAnalysis) SetFight(fa *FightAnalysis) {
	fa.player = pa

	pa.mu.Lock()
	defer pa.mu.Unlock()

	pa.Fights[fa.Name] = fa
}

func sortRemarks(remarks []*remark.Remark) {
	sort.SliceStable(remarks, func(i, j int) bool {
		return remarks[i].UUID < remarks[j].UUID
	})
}

func (pa *PlayerAnalysis) AggregateRemarks() {
	remarks := map[remark.Type][]*remark.Remark{}
	for _, fa := range pa.Fights {
		sortRemarks(fa.Remarks)
		for _, r := range fa.Remarks {
			if _, ok := remarks[r.Type]; !ok {
				remarks[r.Type] = make([]*remark.Remark, 0, 1)
			}
			rr := *r // need to copy otherwise we'll change everywhere
			remarks[r.Type] = append(remarks[r.Type], &rr)
		}
	}
	removeDup := func(remarkType remark.Type, uniqKey func(*remark.Remark) interface{}) {
		uniq := map[interface{}]*remark.Remark{}
		for _, remark := range remarks[remarkType] {
			uK := uniqKey(remark)
			if _, ok := uniq[uK]; !ok {
				remark.Fight = ""
				remark.ComputeUUID()
				uniq[uK] = remark
			}
		}
		rs := make([]*remark.Remark, 0, len(uniq))
		for _, remark := range uniq {
			rs = append(rs, remark)
		}
		remarks[remarkType] = rs
	}
	removeDup(remark.Type_MissingGems, func(r *remark.Remark) interface{} { return r.ItemWowheadAttr })
	removeDup(remark.Type_NoEnchant, func(r *remark.Remark) interface{} { return r.ItemWowheadAttr })
	pa.Remarks = make([]*remark.Remark, 0)
	for _, rks := range remarks {
		pa.Remarks = append(pa.Remarks, rks...)
	}
	sortRemarks(pa.Remarks)
}

type TrueFightAnalysis struct {
	Items       []FightCast `json:"items"`
	Spells      []FightCast `json:"spells"`
	Unknown     []FightCast `json:"unknown"`
	Consumables []FightCast `json:"consumables"`
}

type FightCast struct {
	SpellID int64 `json:"spell_id"`
	Count   int64 `json:"count"`
}

type GearGem struct {
	ID        int64  `json:"id"`
	ItemLevel int64  `json:"itemLevel"`
	Icon      string `json:"icon"`
}

type Aura struct {
	Source  int64      `json:"source"`
	Ability int64      `json:"ability"`
	Stacks  int64      `json:"stacks"`
	Icon    string     `json:"icon"`
	Name    string     `json:"name,omitempty"`
	Events  []struct{} `json:"events"`
}
