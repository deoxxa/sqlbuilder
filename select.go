package sqlbuilder

type AsCommonTableExpression interface {
	AsCommonTableExpression(s *Serializer)
}

type AsDistinct interface {
	AsDistinct(s *Serializer)
}

type AsTableOrSubquery interface {
	AsTableOrSubquery(s *Serializer)
}

type AsResultColumn interface {
	AsResultColumn(s *Serializer)
}

type AsOrderingTerm interface {
	AsOrderingTerm(s *Serializer)
}

type AsOffsetLimit interface {
	AsOffsetLimit(s *Serializer)
}

type SelectStatement struct {
	with        []AsCommonTableExpression
	distinct    AsDistinct
	from        AsTableOrSubquery
	columns     []AsExpr
	where       AsExpr
	orderBy     []AsOrderingTerm
	groupBy     []AsExpr
	offsetLimit AsOffsetLimit
}

func (s *SelectStatement) clone() *SelectStatement {
	return &SelectStatement{
		with:        s.with,
		distinct:    s.distinct,
		from:        s.from,
		columns:     s.columns,
		where:       s.where,
		orderBy:     s.orderBy,
		groupBy:     s.groupBy,
		offsetLimit: s.offsetLimit,
	}
}

func Select() *SelectStatement {
	return &SelectStatement{}
}

func (s *SelectStatement) With(with ...AsCommonTableExpression) *SelectStatement {
	c := s.clone()
	c.with = with
	return c
}

func (s *SelectStatement) Distinct(distinct AsDistinct) *SelectStatement {
	c := s.clone()
	c.distinct = distinct
	return c
}

func (s *SelectStatement) From(from AsTableOrSubquery) *SelectStatement {
	c := s.clone()
	c.from = from
	return c
}

func (s *SelectStatement) Columns(columns ...AsExpr) *SelectStatement {
	c := s.clone()
	c.columns = columns
	return c
}

func (s *SelectStatement) Where(where AsExpr) *SelectStatement {
	c := s.clone()
	c.where = where
	return c
}

func (s *SelectStatement) OrderBy(orderBy ...AsOrderingTerm) *SelectStatement {
	c := s.clone()
	c.orderBy = orderBy
	return c
}

func (s *SelectStatement) GroupBy(groupBy ...AsExpr) *SelectStatement {
	c := s.clone()
	c.groupBy = groupBy
	return c
}

func (s *SelectStatement) OffsetLimit(offsetLimit AsOffsetLimit) *SelectStatement {
	c := s.clone()
	c.offsetLimit = offsetLimit
	return c
}

func (q *SelectStatement) As(alias string) AsExpr {
	return AliasColumn(q, alias)
}

func (q *SelectStatement) AsStatement(s *Serializer) {
	if len(q.with) > 0 {
		s.D("WITH ")

		for i, w := range q.with {
			s.F(w.AsCommonTableExpression).DC(",", i != len(q.with)-1).D(" ")
		}
	}

	s.D("SELECT")

	if q.distinct != nil {
		s.D(" ").F(q.distinct.AsDistinct)
	}

	for i, c := range q.columns {
		s.D(" ")

		if a, ok := c.(AsResultColumn); ok {
			s.F(a.AsResultColumn)
		} else {
			s.F(c.AsExpr)
		}

		s.DC(",", i < len(q.columns)-1)
	}

	if q.from != nil {
		s.D(" FROM ").F(q.from.AsTableOrSubquery)
	}

	if q.where != nil {
		s.D(" WHERE ").F(q.where.AsExpr)
	}

	if len(q.groupBy) != 0 {
		s.D(" GROUP BY")

		for i, e := range q.groupBy {
			s.D(" ").F(e.AsExpr).DC(",", i < len(q.groupBy)-1)
		}
	}

	if len(q.orderBy) > 0 {
		s.D(" ORDER BY")

		for i, e := range q.orderBy {
			s.D(" ").F(e.AsOrderingTerm).DC(",", i < len(q.orderBy)-1)
		}
	}

	if q.offsetLimit != nil {
		s.D(" ").F(q.offsetLimit.AsOffsetLimit)
	}
}

func (q *SelectStatement) AsExpr(s *Serializer) {
	s.D("(").F(q.AsStatement).D(")")
}

func (q *SelectStatement) AsTableOrSubquery(s *Serializer) {
	s.D("(").F(q.AsStatement).D(")")
}

func (q *SelectStatement) C(name string) *BasicColumn {
	return &BasicColumn{name: name}
}
