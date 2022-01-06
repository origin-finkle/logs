package models

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
)

func TestTalents(t *testing.T) {
	for _, data := range []struct {
		SubType                   string
		Talents                   [3]int64
		ExpectedSpec              Specialization
		BenefitsFromWindfuryTotem bool
		ExpectedRoles             []Role
	}{
		{
			SubType:                   "Warlock",
			Talents:                   [3]int64{0, 21, 40},
			ExpectedSpec:              Specialization_DestructionWarlock,
			BenefitsFromWindfuryTotem: false,
		},
		{
			SubType:                   "Warrior",
			Talents:                   [3]int64{0, 21, 40},
			ExpectedSpec:              Specialization_ProtectionWarrior,
			BenefitsFromWindfuryTotem: true,
		},
		{
			SubType:                   "Warrior",
			Talents:                   [3]int64{40, 21, 0},
			ExpectedSpec:              Specialization_ArmsWarrior,
			BenefitsFromWindfuryTotem: true,
		},
		{
			SubType:                   "Warrior",
			Talents:                   [3]int64{21, 40, 0},
			ExpectedSpec:              Specialization_FuryWarrior,
			BenefitsFromWindfuryTotem: true,
		},
		{
			SubType:                   "Rogue",
			Talents:                   [3]int64{0, 21, 40},
			ExpectedSpec:              Specialization_SubtletyRogue,
			BenefitsFromWindfuryTotem: true,
			ExpectedRoles:             []Role{Role_Physical, Role_Melee},
		},
		{
			SubType:                   "Rogue",
			Talents:                   [3]int64{40, 21, 0},
			ExpectedSpec:              Specialization_AssassinationRogue,
			BenefitsFromWindfuryTotem: true,
			ExpectedRoles:             []Role{Role_Physical, Role_Melee},
		},
		{
			SubType:                   "Rogue",
			Talents:                   [3]int64{21, 40, 0},
			ExpectedSpec:              Specialization_CombatRogue,
			BenefitsFromWindfuryTotem: true,
			ExpectedRoles:             []Role{Role_Physical, Role_Melee},
		},
		{
			SubType:                   "Paladin",
			Talents:                   [3]int64{0, 21, 40},
			ExpectedSpec:              Specialization_RetributionPaladin,
			BenefitsFromWindfuryTotem: true,
		},
		{
			SubType:      "Druid",
			Talents:      [3]int64{36, 0, 25},
			ExpectedSpec: Specialization_RestorationDruid,
		},
		{
			SubType:      "Druid",
			Talents:      [3]int64{34, 0, 27},
			ExpectedSpec: Specialization_RestorationDruid,
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
			td.Cmp(tt, talents.BenefitsFromWindfuryTotem(), data.BenefitsFromWindfuryTotem)
			for _, role := range data.ExpectedRoles {
				td.CmpTrue(tt, talents.Spec.IsRole(role))
			}
		})
	}
}
