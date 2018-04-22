package sqlbuilder

type ExpressionTableExpr struct {
	expr    AsExpr
	name    string
	names   []string
	columns []*BasicColumn
}

func ExpressionTable(expr AsExpr, name string, names ...string) *ExpressionTableExpr {
	t := &ExpressionTableExpr{expr: expr, name: name, names: names}

	columns := make([]*BasicColumn, len(names))
	for i, v := range names {
		columns[i] = &BasicColumn{table: t, name: v}
	}
	t.columns = columns

	return t
}

func (t *ExpressionTableExpr) AsNamed(s *Serializer) {
	s.N(t.name)
}

func (t *ExpressionTableExpr) AsTableOrSubquery(s *Serializer) {
	s.F(t.expr.AsExpr).D(" AS ").N(t.name).D(" ")

	s.D("(")
	for i, name := range t.names {
		s.N(name).DC(", ", i != len(t.names)-1)
	}
	s.D(")")
}

func (t *ExpressionTableExpr) C(name string) *BasicColumn {
	for _, c := range t.columns {
		if c.name == name {
			return c
		}
	}

	return nil
}

func (t *ExpressionTableExpr) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, t, right)
}

func (t *ExpressionTableExpr) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(t, right)
}

func (t *ExpressionTableExpr) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(t, right)
}
