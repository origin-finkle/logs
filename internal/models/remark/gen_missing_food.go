// Code generated by gen.py; DO NOT EDIT
package remark

type MissingFood struct {
}

func (MissingFood) getType() Type   { return Type_MissingFood }
func (MissingFood) is(rt Type) bool { return rt == Type_MissingFood }
func (md MissingFood) apply(r *Remark) {
}
