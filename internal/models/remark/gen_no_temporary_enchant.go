// Code generated by gen.py; DO NOT EDIT
package remark

type NoTemporaryEnchant struct {
	ItemWowheadAttr string `json:"item_wowhead_attr"`
}

func (NoTemporaryEnchant) getType() Type   { return Type_NoTemporaryEnchant }
func (NoTemporaryEnchant) is(rt Type) bool { return rt == Type_NoTemporaryEnchant }
func (md NoTemporaryEnchant) apply(r *Remark) {
	r.ItemWowheadAttr = md.ItemWowheadAttr
}
