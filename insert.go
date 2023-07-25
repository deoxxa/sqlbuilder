package sqlbuilder

type InsertColumns map[*BasicColumn]AsExpr

type InsertStatement struct {
	table     *Table
	columns   InsertColumns
	returning []AsExpr
}

func (s *InsertStatement) clone() *InsertStatement {
	return &InsertStatement{
		table:     s.table,
		columns:   s.columns,
		returning: s.returning,
	}
}

func Insert() *InsertStatement {
	return &InsertStatement{columns: make(InsertColumns)}
}

func (s *InsertStatement) Table(table *Table) *InsertStatement {
	c := s.clone()
	c.table = table
	return c
}

func (s *InsertStatement) Columns(columns InsertColumns) *InsertStatement {
	c := s.clone()
	c.columns = columns
	return c
}

func (s *InsertStatement) Returning(returning ...AsExpr) *InsertStatement {
	c := s.clone()
	c.returning = returning
	return c
}

func (q *InsertStatement) AsStatement(s *Serializer) {
	s.D("INSERT INTO ").F(q.table.AsNamed).D(" ")

	var keys []AsNamedShort
	var vals []AsExpr

	for k, v := range q.columns {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	s.D("(")
	for i, k := range keys {
		s.F(k.AsNamedShort).DC(", ", i < len(keys)-1)
	}
	s.D(") VALUES (")
	for i, k := range vals {
		s.F(k.AsExpr).DC(", ", i < len(keys)-1)
	}
	s.D(")")

	if len(q.returning) > 0 {
		s.D(" RETURNING ")

		for i, c := range q.returning {
			if a, ok := c.(AsResultColumn); ok {
				s.F(a.AsResultColumn)
			} else {
				s.F(c.AsExpr)
			}

			s.DC(", ", i < len(q.returning)-1)
		}
	}
}
