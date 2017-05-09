package sqlbuilder

import (
	"fmt"
)

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
	from        AsTableOrSubquery
	columns     []AsExpr
	where       AsExpr
	orderBy     []AsOrderingTerm
	groupBy     []AsExpr
	offsetLimit AsOffsetLimit
}

func (s *SelectStatement) clone() *SelectStatement {
	return &SelectStatement{
		from:        s.from,
		columns:     s.columns[:],
		where:       s.where,
		orderBy:     s.orderBy[:],
		groupBy:     s.groupBy[:],
		offsetLimit: s.offsetLimit,
	}
}

func Select() *SelectStatement {
	return &SelectStatement{}
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

func (q *SelectStatement) Serialize(s *Serializer) {
	s.D("SELECT")

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

func OrderAsc(expr AsExpr) *OrderingTerm  { return &OrderingTerm{expr: expr, order: "ASC"} }
func OrderDesc(expr AsExpr) *OrderingTerm { return &OrderingTerm{expr: expr, order: "DESC"} }

type OrderingTerm struct {
	expr  AsExpr
	order string
}

func (t *OrderingTerm) AsOrderingTerm(s *Serializer) {
	s.F(t.expr.AsExpr).D(" ").D(t.order)
}

func OffsetLimit(offset, limit uint) *OffsetLimitClause {
	return &OffsetLimitClause{offset: offset, limit: limit}
}

type OffsetLimitClause struct {
	offset, limit uint
}

func (c *OffsetLimitClause) AsOffsetLimit(s *Serializer) {
	if c.limit != 0 {
		s.D(fmt.Sprintf("LIMIT %d", c.limit))

		if c.offset != 0 {
			s.D(fmt.Sprintf(" OFFSET %d", c.offset))
		}
	}
}
