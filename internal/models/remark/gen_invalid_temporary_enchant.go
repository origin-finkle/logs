// Code generated by gen.py; DO NOT EDIT
package remark

type InvalidTemporaryEnchant struct {
	ItemWowheadAttr string `json:"item_wowhead_attr"`
}

func (InvalidTemporaryEnchant) getType() Type   { return Type_InvalidTemporaryEnchant }
func (InvalidTemporaryEnchant) is(rt Type) bool { return rt == Type_InvalidTemporaryEnchant }
func (md InvalidTemporaryEnchant) apply(r *Remark) {
	r.ItemWowheadAttr = md.ItemWowheadAttr
}
