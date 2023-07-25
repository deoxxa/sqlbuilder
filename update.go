package sqlbuilder

type UpdateColumns map[*BasicColumn]AsExpr

func (m UpdateColumns) asFields() []UpdateField {
	var a []UpdateField
	for k, e := range m {
		a = append(a, UpdateField{k, e})
	}
	return a
}

type UpdateField struct {
	Name  AsNamedShort
	Value AsExpr
}

type UpdateStatement struct {
	with      []AsCommonTableExpression
	target    AsTableOrSubquery
	columns   UpdateColumns
	fields    []UpdateField
	from      AsTableOrSubquery
	where     AsExpr
	returning []AsExpr
}

func (s *UpdateStatement) clone() *UpdateStatement {
	return &UpdateStatement{
		with:      s.with,
		target:    s.target,
		columns:   s.columns,
		fields:    s.fields,
		from:      s.from,
		where:     s.where,
		returning: s.returning,
	}
}

func Update() *UpdateStatement {
	return &UpdateStatement{columns: make(UpdateColumns)}
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

func (s *UpdateStatement) Set(columns UpdateColumns) *UpdateStatement {
	c := s.clone()
	c.columns = columns
	c.fields = nil
	return c
}

func (s *UpdateStatement) AndSet(column *BasicColumn, expr AsExpr) *UpdateStatement {
	c := s.clone()

	m := make(UpdateColumns)
	for k, e := range c.columns {
		m[k] = e
	}

	m[column] = expr

	c.columns = m

	return c
}

func (s *UpdateStatement) Fields(fields []UpdateField) *UpdateStatement {
	c := s.clone()
	c.columns = nil
	c.fields = fields
	return c
}

func (s *UpdateStatement) AndField(field UpdateField) *UpdateStatement {
	c := s.clone()
	c.fields = append(c.fields[:], field)
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

	for i, f := range append(q.columns.asFields(), q.fields...) {
		s.DC(", ", i != 0).F(f.Name.AsNamedShort).D(" = ").F(f.Value.AsExpr)
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
