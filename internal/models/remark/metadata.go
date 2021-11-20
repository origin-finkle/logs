package remark

type Metadata interface {
	getType() Type
	is(Type) bool
	apply(*Remark)
}
