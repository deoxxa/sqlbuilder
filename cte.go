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
	s.N(e.name).D(" AS (").F(e.as.AsStatement).D(")")
}

func (e *CommonTableExpressionExpr) IsRecursive() bool {
	return e.recursive
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

func (e *CommonTableExpressionExpr) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, e, right)
}

func (e *CommonTableExpressionExpr) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(e, right)
}

func (e *CommonTableExpressionExpr) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(e, right)
}
