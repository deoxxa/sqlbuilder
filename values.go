package sqlbuilder

type ValuesExpr struct {
	rows []*ValuesRowExpr
}

func Values(rows ...*ValuesRowExpr) *ValuesExpr {
	return &ValuesExpr{rows: rows}
}

func (v *ValuesExpr) AsExpr(s *Serializer) {
	s.D("VALUES ")

	for i, r := range v.rows {
		s.F(r.AsExpr).DC(", ", i != len(v.rows)-1)
	}
}

type ValuesRowExpr struct {
	values []AsExpr
}

func ValuesRow(values ...AsExpr) *ValuesRowExpr {
	return &ValuesRowExpr{values: values}
}

func (r *ValuesRowExpr) AsExpr(s *Serializer) {
	s.D("(")

	for i, v := range r.values {
		s.F(v.AsExpr).DC(", ", i != len(r.values)-1)
	}

	s.D(")")
}

type ValuesTableExpr struct {
	expr    *ValuesExpr
	name    string
	names   []string
	columns []*BasicColumn
}

func ValuesTable(expr *ValuesExpr, name string, names ...string) *ValuesTableExpr {
	t := &ValuesTableExpr{expr: expr, name: name, names: names}

	columns := make([]*BasicColumn, len(names))
	for i, v := range names {
		columns[i] = &BasicColumn{table: t, name: v}
	}
	t.columns = columns

	return t
}

func (t *ValuesTableExpr) AsNamed(s *Serializer) {
	s.N(t.name)
}

func (t *ValuesTableExpr) AsTableOrSubquery(s *Serializer) {
	s.D("(").F(t.expr.AsExpr).D(") AS ").N(t.name).D(" ")

	s.D("(")
	for i, name := range t.names {
		s.N(name).DC(", ", i != len(t.names)-1)
	}
	s.D(")")
}

func (t *ValuesTableExpr) C(name string) *BasicColumn {
	for _, c := range t.columns {
		if c.name == name {
			return c
		}
	}

	return nil
}

func (t *ValuesTableExpr) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, t, right)
}

func (t *ValuesTableExpr) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(t, right)
}

func (t *ValuesTableExpr) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(t, right)
}
