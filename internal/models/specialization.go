package models

type Specialization string

const (
	Specialization_RetributionPaladin Specialization = "RetributionPaladin"
	Specialization_HolyPaladin        Specialization = "HolyPaladin"
	Specialization_ProtectionPaladin  Specialization = "ProtectionPaladin"
	Specialization_CombatRogue        Specialization = "CombatRogue"
	Specialization_AssassinationRogue Specialization = "AssassinationRogue"
	Specialization_SubtletyRogue      Specialization = "SubtletyRogue"
	Specialization_ArmsWarrior        Specialization = "ArmsWarrior"
	Specialization_FuryWarrior        Specialization = "FuryWarrior"
	Specialization_ProtectionWarrior  Specialization = "ProtectionWarrior"
	Specialization_EnhancementShaman  Specialization = "EnhancementShaman"
	Specialization_ElementalShaman    Specialization = "ElementalShaman"
	Specialization_RestorationShaman  Specialization = "RestorationShaman"
	Specialization_BalanceDruid       Specialization = "BalanceDruid"
	Specialization_CatDruid           Specialization = "CatDruid"
	Specialization_BearDruid          Specialization = "BearDruid"
	Specialization_FeralDruid         Specialization = "FeralDruid"
	Specialization_RestorationDruid   Specialization = "RestorationDruid"
	Specialization_DisciplinePriest   Specialization = "DisciplinePriest"
	Specialization_HolyPriest         Specialization = "HolyPriest"
	Specialization_ShadowPriest       Specialization = "ShadowPriest"
	Specialization_AfflictionWarlock  Specialization = "AfflictionWarlock"
	Specialization_DemonologyWarlock  Specialization = "DemonologyWarlock"
	Specialization_DestructionWarlock Specialization = "DestructionWarlock"
	Specialization_ArcaneMage         Specialization = "ArcaneMage"
	Specialization_FireMage           Specialization = "FireMage"
	Specialization_FrostMage          Specialization = "FrostMage"
	Specialization_SurvivalHunter     Specialization = "SurvivalHunter"
	Specialization_MarksmanshipHunter Specialization = "MarksmanshipHunter"
	Specialization_BeastMasteryHunter Specialization = "BeastMasteryHunter"
	Specialization_Unknown            Specialization = "Unknown"
)

var (
	isSpecialization = map[string]bool{
		string(Specialization_RetributionPaladin): true,
		string(Specialization_HolyPaladin):        true,
		string(Specialization_ProtectionPaladin):  true,
		string(Specialization_CombatRogue):        true,
		string(Specialization_AssassinationRogue): true,
		string(Specialization_SubtletyRogue):      true,
		string(Specialization_ArmsWarrior):        true,
		string(Specialization_FuryWarrior):        true,
		string(Specialization_ProtectionWarrior):  true,
		string(Specialization_EnhancementShaman):  true,
		string(Specialization_ElementalShaman):    true,
		string(Specialization_RestorationShaman):  true,
		string(Specialization_BalanceDruid):       true,
		string(Specialization_CatDruid):           true,
		string(Specialization_BearDruid):          true,
		string(Specialization_FeralDruid):         true,
		string(Specialization_RestorationDruid):   true,
		string(Specialization_DisciplinePriest):   true,
		string(Specialization_HolyPriest):         true,
		string(Specialization_ShadowPriest):       true,
		string(Specialization_AfflictionWarlock):  true,
		string(Specialization_DemonologyWarlock):  true,
		string(Specialization_DestructionWarlock): true,
		string(Specialization_ArcaneMage):         true,
		string(Specialization_FireMage):           true,
		string(Specialization_FrostMage):          true,
		string(Specialization_SurvivalHunter):     true,
		string(Specialization_MarksmanshipHunter): true,
		string(Specialization_BeastMasteryHunter): true,
		string(Specialization_Unknown):            true,
	}
)

func (s Specialization) IsRole(role Role) bool {
	return specs[s].Role[role]
}

func stringIsSpecialization(s string) bool {
	return isSpecialization[s]
}
