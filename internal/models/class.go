package models

type Class string

const (
	Class_Paladin Class = "Paladin"
	Class_Rogue   Class = "Rogue"
	Class_Warrior Class = "Warrior"
	Class_Shaman  Class = "Shaman"
	Class_Druid   Class = "Druid"
	Class_Priest  Class = "Priest"
	Class_Warlock Class = "Warlock"
	Class_Mage    Class = "Mage"
	Class_Hunter  Class = "Hunter"
	Class_Unknown Class = "Unknown"
)

var (
	isClass = map[string]bool{
		string(Class_Paladin): true,
		string(Class_Rogue):   true,
		string(Class_Warrior): true,
		string(Class_Shaman):  true,
		string(Class_Druid):   true,
		string(Class_Priest):  true,
		string(Class_Warlock): true,
		string(Class_Mage):    true,
		string(Class_Hunter):  true,
	}
)

func stringIsClass(c string) bool { return isClass[c] }
