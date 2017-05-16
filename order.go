package sqlbuilder

func OrderAsc(expr AsExpr) *OrderingTerm     { return &OrderingTerm{expr: expr, order: "ASC"} }
func OrderDesc(expr AsExpr) *OrderingTerm    { return &OrderingTerm{expr: expr, order: "DESC"} }
func OrderDefault(expr AsExpr) *OrderingTerm { return &OrderingTerm{expr: expr} }

type OrderingTerm struct {
	expr  AsExpr
	order string
}

func (t *OrderingTerm) AsOrderingTerm(s *Serializer) {
	s.F(t.expr.AsExpr).DC(" ", t.order != "").DC(t.order, t.order != "")
}
