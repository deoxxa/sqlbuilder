package sqlbuilder

type FilterExpr struct {
	left  AsExpr
	right AsExpr
}

func Filter(left, right AsExpr) *FilterExpr {
	return &FilterExpr{left: left, right: right}
}

func (c *FilterExpr) AsExpr(s *Serializer) {
	s.F(c.left.AsExpr).D(" FILTER (WHERE ").F(c.right.AsExpr).D(")")
}
