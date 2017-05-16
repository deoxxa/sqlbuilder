package sqlbuilder

type CommonTableExpressionExpr struct {
	name      string
	recursive bool
	as        *SelectStatement
}

func CommonTableExpression(name string) *CommonTableExpressionExpr {
	return &CommonTableExpressionExpr{name: name}
}

func (e *CommonTableExpressionExpr) Recursive(recursive bool) *CommonTableExpressionExpr {
	return &CommonTableExpressionExpr{name: e.name, as: e.as, recursive: recursive}
}

func (e *CommonTableExpressionExpr) As(as *SelectStatement) *CommonTableExpressionExpr {
	return &CommonTableExpressionExpr{name: e.name, recursive: e.recursive, as: as}
}

func (e *CommonTableExpressionExpr) AsCommonTableExpression(s *Serializer) {
	if e.recursive {
		s.D("RECURSIVE ")
	}

	s.N(e.name).D(" AS (").F(e.as.AsStatement).D(")")
}

func (e *CommonTableExpressionExpr) AsNamed(s *Serializer) {
	s.N(e.name)
}

func (e *CommonTableExpressionExpr) C(name string) *BasicColumn {
	return &BasicColumn{table: e, name: name}
}

func (e *CommonTableExpressionExpr) AsTableOrSubquery(s *Serializer) {
	s.N(e.name)
}
