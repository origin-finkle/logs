package models

type Talents struct {
	fight *FightAnalysis `json:"-"`

	Points [3]int64       `json:"points"`
	Spec   Specialization `json:"spec"`
}

func NewTalents(fight *FightAnalysis, points [3]int64) *Talents {
	t := &Talents{
		fight:  fight,
		Points: points,
	}
	t.guessSpec()
	return t
}

func (t *Talents) guessSpec() {
	for spec, data := range specs {
		if t.fight.player.SubType == string(data.Class) {
			match := false
			switch data.HasMorePointsIn {
			case 0:
				match = t.Points[0] >= t.Points[1] && t.Points[0] >= t.Points[2]
			case 1:
				match = t.Points[1] >= t.Points[0] && t.Points[1] >= t.Points[2]
			case 2:
				match = t.Points[2] >= t.Points[0] && t.Points[2] >= t.Points[1]
			}
			if match {
				t.Spec = spec
				return
			}
		}
	}
}

func (t *Talents) BenefitsFromWindfuryTotem() bool {
	return specs[t.Spec].BenefitsFromWindfuryTotem
}

var (
	specs = map[Specialization]struct {
		HasMorePointsIn           int
		Class                     Class
		Role                      map[Role]bool
		BenefitsFromWindfuryTotem bool
	}{
		Specialization_HolyPaladin: {
			HasMorePointsIn: 0,
			Class:           Class_Paladin,
			Role:            map[Role]bool{Role_Heal: true},
		},
		Specialization_ProtectionPaladin: {
			HasMorePointsIn: 1,
			Class:           Class_Paladin,
			Role:            map[Role]bool{Role_Tank: true, Role_Magic: true},
		},
		Specialization_RetributionPaladin: {
			HasMorePointsIn:           2,
			Class:                     Class_Paladin,
			Role:                      map[Role]bool{Role_Melee: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_AssassinationRogue: {
			HasMorePointsIn:           0,
			Class:                     Class_Rogue,
			Role:                      map[Role]bool{Role_Melee: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_CombatRogue: {
			HasMorePointsIn:           1,
			Class:                     Class_Rogue,
			Role:                      map[Role]bool{Role_Melee: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_SubtletyRogue: {
			HasMorePointsIn:           2,
			Class:                     Class_Rogue,
			Role:                      map[Role]bool{Role_Melee: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_ArmsWarrior: {
			HasMorePointsIn:           0,
			Class:                     Class_Warrior,
			Role:                      map[Role]bool{Role_Melee: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_FuryWarrior: {
			HasMorePointsIn:           1,
			Class:                     Class_Warrior,
			Role:                      map[Role]bool{Role_Melee: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_ProtectionWarrior: {
			HasMorePointsIn:           2,
			Class:                     Class_Warrior,
			Role:                      map[Role]bool{Role_Tank: true, Role_Physical: true},
			BenefitsFromWindfuryTotem: true,
		},
		Specialization_ElementalShaman: {
			HasMorePointsIn: 0,
			Class:           Class_Shaman,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_EnhancementShaman: {
			HasMorePointsIn: 1,
			Class:           Class_Shaman,
			Role:            map[Role]bool{Role_Melee: true, Role_Physical: true},
		},
		Specialization_RestorationShaman: {
			HasMorePointsIn: 2,
			Class:           Class_Shaman,
			Role:            map[Role]bool{Role_Heal: true},
		},
		Specialization_BalanceDruid: {
			HasMorePointsIn: 0,
			Class:           Class_Druid,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_FeralDruid: {
			HasMorePointsIn: 1,
			Class:           Class_Druid,
			Role:            map[Role]bool{Role_Melee: true, Role_Tank: true, Role_Physical: true},
		},
		Specialization_RestorationDruid: {
			HasMorePointsIn: 2,
			Class:           Class_Druid,
			Role:            map[Role]bool{Role_Heal: true},
		},
		Specialization_DisciplinePriest: {
			HasMorePointsIn: 0,
			Class:           Class_Priest,
			Role:            map[Role]bool{Role_Heal: true},
		},
		Specialization_HolyPriest: {
			HasMorePointsIn: 1,
			Class:           Class_Priest,
			Role:            map[Role]bool{Role_Heal: true},
		},
		Specialization_ShadowPriest: {
			HasMorePointsIn: 2,
			Class:           Class_Priest,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_AfflictionWarlock: {
			HasMorePointsIn: 0,
			Class:           Class_Warlock,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_DemonologyWarlock: {
			HasMorePointsIn: 1,
			Class:           Class_Warlock,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_DestructionWarlock: {
			HasMorePointsIn: 2,
			Class:           Class_Warlock,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_ArcaneMage: {
			HasMorePointsIn: 0,
			Class:           Class_Mage,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_FireMage: {
			HasMorePointsIn: 1,
			Class:           Class_Mage,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_FrostMage: {
			HasMorePointsIn: 2,
			Class:           Class_Mage,
			Role:            map[Role]bool{Role_Ranged: true, Role_Magic: true},
		},
		Specialization_SurvivalHunter: {
			HasMorePointsIn: 0,
			Class:           Class_Hunter,
			Role:            map[Role]bool{Role_Ranged: true, Role_Physical: true},
		},
		Specialization_MarksmanshipHunter: {
			HasMorePointsIn: 1,
			Class:           Class_Hunter,
			Role:            map[Role]bool{Role_Ranged: true, Role_Physical: true},
		},
		Specialization_BeastMasteryHunter: {
			HasMorePointsIn: 2,
			Class:           Class_Hunter,
			Role:            map[Role]bool{Role_Ranged: true, Role_Physical: true},
		},
		Specialization_Unknown: {
			HasMorePointsIn: 0,
			Class:           Class_Unknown,
			Role:            map[Role]bool{},
		},
	}
)
