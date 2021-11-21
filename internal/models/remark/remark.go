package remark

import "strings"

func New(metadata Metadata, fightName string) *Remark {
	r := &Remark{
		Type:  metadata.getType(),
		Fight: fightName,
	}
	metadata.apply(r)
	r.ComputeUUID()
	return r
}

func (r *Remark) ComputeUUID() {
	uuid := []string{string(r.Type)}
	if r.WowheadAttr != "" {
		uuid = append(uuid, r.WowheadAttr)
	}
	if r.ItemWowheadAttr != "" {
		uuid = append(uuid, r.ItemWowheadAttr)
	}
	if r.SpellWowheadAttr != "" {
		uuid = append(uuid, r.SpellWowheadAttr)
	}
	if r.Slot != "" {
		uuid = append(uuid, r.Slot)
	}
	if r.Fight != "" {
		uuid = append(uuid, r.Fight)
	}
	r.UUID = strings.Join(uuid, ":")
}
