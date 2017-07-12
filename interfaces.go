package sqlbuilder

type AsNamed interface {
	AsNamed(s *Serializer)
}

type AsNamedShort interface {
	AsNamedShort(s *Serializer)
}

type AsExpr interface {
	AsExpr(s *Serializer)
}

type WithColumns interface {
	C(name string) *BasicColumn
}
