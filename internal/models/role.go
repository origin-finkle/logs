package models

type Role string

const (
	Role_Tank     Role = "Tank"
	Role_Heal     Role = "Heal"
	Role_Melee    Role = "Melee"
	Role_Ranged   Role = "Ranged"
	Role_Magic    Role = "Magic"
	Role_Physical Role = "Physical"
)

var (
	isRole = map[string]bool{
		string(Role_Tank):     true,
		string(Role_Heal):     true,
		string(Role_Melee):    true,
		string(Role_Ranged):   true,
		string(Role_Magic):    true,
		string(Role_Physical): true,
	}
)

func stringIsRole(r string) bool {
	return isRole[r]
}
