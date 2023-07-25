package sqlbuilder

type UpdateColumns map[*BasicColumn]AsExpr

type UpdateStatement struct {
	with      []AsCommonTableExpression
	target    AsTableOrSubquery
	set       UpdateColumns
	from      AsTableOrSubquery
	where     AsExpr
	returning []AsExpr
}

func (s *UpdateStatement) clone() *UpdateStatement {
	return &UpdateStatement{
		with:      s.with,
		target:    s.target,
		set:       s.set,
		from:      s.from,
		where:     s.where,
		returning: s.returning,
	}
}

func Update() *UpdateStatement {
	return &UpdateStatement{set: make(UpdateColumns)}
}

func (s *UpdateStatement) With(with ...AsCommonTableExpression) *UpdateStatement {
	c := s.clone()
	c.with = with
	return c
}

func (s *UpdateStatement) AndWith(with ...AsCommonTableExpression) *UpdateStatement {
	c := s.clone()
	c.with = append(c.with, with...)
	return c
}

func (s *UpdateStatement) Table(table *Table) *UpdateStatement {
	return s.Target(table)
}

func (s *UpdateStatement) Target(target AsTableOrSubquery) *UpdateStatement {
	c := s.clone()
	c.target = target
	return c
}

func (s *UpdateStatement) GetTarget() AsTableOrSubquery {
	return s.target
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

func (s *UpdateStatement) From(from AsTableOrSubquery) *UpdateStatement {
	c := s.clone()
	c.from = from
	return c
}

func (s *UpdateStatement) Where(where AsExpr) *UpdateStatement {
	c := s.clone()
	c.where = where
	return c
}

func (s *UpdateStatement) Returning(returning ...AsExpr) *UpdateStatement {
	c := s.clone()
	c.returning = returning
	return c
}

func (q *UpdateStatement) AsStatement(s *Serializer) {
	if len(q.with) > 0 {
		s.D("WITH ")

		for _, w := range q.with {
			if w.IsRecursive() {
				s.D("RECURSIVE ")
				break
			}
		}

		for i, w := range q.with {
			s.F(w.AsCommonTableExpression).DC(",", i != len(q.with)-1).D(" ")
		}
	}

	s.D("UPDATE ").F(q.target.AsTableOrSubquery).D(" SET ")

	i := 0
	for k, e := range q.set {
		s.F(k.AsNamedShort).D(" = ").F(e.AsExpr).DC(", ", i < len(q.set)-1)
		i++
	}

	if q.from != nil {
		s.D(" FROM ").F(q.from.AsTableOrSubquery)
	}

	if q.where != nil {
		s.D(" WHERE ").F(q.where.AsExpr)
	}

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
