// Code generated by gen.py; DO NOT EDIT
package remark

type NoEnchant struct {
	ItemWowheadAttr string `json:"item_wowhead_attr"`
}

func (NoEnchant) getType() Type   { return Type_NoEnchant }
func (NoEnchant) is(rt Type) bool { return rt == Type_NoEnchant }
func (md NoEnchant) apply(r *Remark) {
	r.ItemWowheadAttr = md.ItemWowheadAttr
}
