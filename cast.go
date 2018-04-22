package sqlbuilder

type CastExpr struct {
	expr AsExpr
	as   string
}

func Cast(expr AsExpr, as string) *CastExpr {
	return &CastExpr{expr: expr, as: as}
}

func (c *CastExpr) AsExpr(s *Serializer) {
	s.D("CAST (").F(c.expr.AsExpr).D(" AS ").D(c.as).D(")")
}
