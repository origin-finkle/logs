// Code generated by gen.py; DO NOT EDIT
package remark

type NoTemporaryEnchantButWindfury struct {
	ItemWowheadAttr string `json:"item_wowhead_attr"`
}

func (NoTemporaryEnchantButWindfury) getType() Type { return Type_NoTemporaryEnchantButWindfury }
func (NoTemporaryEnchantButWindfury) is(rt Type) bool {
	return rt == Type_NoTemporaryEnchantButWindfury
}
func (md NoTemporaryEnchantButWindfury) apply(r *Remark) {
	r.ItemWowheadAttr = md.ItemWowheadAttr
}