package models

type CommonConfig struct {
	Invalid                   bool     `json:"invalid"`
	InvalidReason             string   `json:"invalid_reason"`
	RestrictedRoles           []string `json:"restricted_roles"`
	RestrictedSpecializations []string `json:"restricted_specializations"`
	RestrictedFights          []string `json:"restricted_fights"`
}
