package models

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestTalents(t *testing.T) {
	for _, data := range []struct {
		SubType      string
		Talents      [3]int64
		ExpectedSpec Specialization
	}{
		{
			SubType:      "Warlock",
			Talents:      [3]int64{0, 21, 40},
			ExpectedSpec: Specialization_DestructionWarlock,
		},
	} {
		t.Run(string(data.ExpectedSpec), func(tt *testing.T) {
			talents := NewTalents(&FightAnalysis{
				player: &PlayerAnalysis{
					Actor: Actor{
						SubType: data.SubType,
					},
				},
			}, data.Talents)
			td.Cmp(tt, talents.Spec, data.ExpectedSpec)
		})
	}
}
