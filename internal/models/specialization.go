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

func (s Specialization) IsRole(role Role) bool {
	return specs[s].Role[role]
}
