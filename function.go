package sqlbuilder

func Func(name string, args ...AsExpr) *FuncExpr {
	return &FuncExpr{name: name, args: args}
}

type FuncExpr struct {
	name string
	args []AsExpr
}

func (e *FuncExpr) AsExpr(s *Serializer) {
	s.D(e.name).D("(")

	for i, arg := range e.args {
		s.F(arg.AsExpr).DC(", ", i != len(e.args)-1)
	}

	s.D(")")
}

func (e *FuncExpr) As(alias string) *ColumnAlias {
	return AliasColumn(e, alias)
}
