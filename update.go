package sqlbuilder

type UpdateColumns map[*BasicColumn]AsExpr

type UpdateStatement struct {
	table *Table
	where AsExpr
	set   UpdateColumns
}

func (s *UpdateStatement) clone() *UpdateStatement {
	return &UpdateStatement{
		table: s.table,
		where: s.where,
		set:   s.set,
	}
}

func Update() *UpdateStatement {
	return &UpdateStatement{set: make(UpdateColumns)}
}

func (s *UpdateStatement) Table(table *Table) *UpdateStatement {
	c := s.clone()
	c.table = table
	return c
}

func (s *UpdateStatement) Set(set UpdateColumns) *UpdateStatement {
	c := s.clone()
	c.set = set
	return c
}

func (s *UpdateStatement) AndSet(column *BasicColumn, expr AsExpr) *UpdateStatement {
	c := s.clone()

	m := make(UpdateColumns)
	for k, e := range c.set {
		m[k] = e
	}

	m[column] = expr

	c.set = m

	return c
}

func (s *UpdateStatement) Where(where AsExpr) *UpdateStatement {
	c := s.clone()
	c.where = where
	return c
}

func (q *UpdateStatement) AsStatement(s *Serializer) {
	s.D("UPDATE ").F(q.table.AsNamed).D(" SET ")

	i := 0
	for k, e := range q.set {
		s.F(k.AsNamedShort).D(" = ").F(e.AsExpr).DC(", ", i < len(q.set)-1)
		i++
	}

	if q.where != nil {
		s.D(" WHERE ").F(q.where.AsExpr)
	}
}
