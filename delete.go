package sqlbuilder

type DeleteStatement struct {
	table *Table
	where AsExpr
}

func (s *DeleteStatement) clone() *DeleteStatement {
	return &DeleteStatement{
		table: s.table,
		where: s.where,
	}
}

func Delete() *DeleteStatement {
	return &DeleteStatement{}
}

func (s *DeleteStatement) Table(table *Table) *DeleteStatement {
	c := s.clone()
	c.table = table
	return c
}

func (s *DeleteStatement) Where(where AsExpr) *DeleteStatement {
	c := s.clone()
	c.where = where
	return c
}

func (q *DeleteStatement) AsStatement(s *Serializer) {
	s.D("DELETE FROM ").F(q.table.AsNamed)

	if q.where != nil {
		s.D(" WHERE ").F(q.where.AsExpr)
	}
}
