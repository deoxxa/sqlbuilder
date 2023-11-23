package sqlbuilder

func OrderAsc(expr AsExpr) *OrderingTerm     { return &OrderingTerm{expr: expr, order: "ASC"} }
func OrderAscNullsFirst(expr AsExpr) *OrderingTerm     { return &OrderingTerm{expr: expr, order: "ASC", nulls: "FIRST"} }
func OrderAscNullsLast(expr AsExpr) *OrderingTerm     { return &OrderingTerm{expr: expr, order: "ASC", nulls: "LAST"} }
func OrderDesc(expr AsExpr) *OrderingTerm    { return &OrderingTerm{expr: expr, order: "DESC"} }
func OrderDescNullsFirst(expr AsExpr) *OrderingTerm    { return &OrderingTerm{expr: expr, order: "DESC", nulls: "FIRST"} }
func OrderDescNullsLast(expr AsExpr) *OrderingTerm    { return &OrderingTerm{expr: expr, order: "DESC", nulls: "LAST"} }
func OrderDefault(expr AsExpr) *OrderingTerm { return &OrderingTerm{expr: expr} }
func OrderDefaultNullsFirst(expr AsExpr) *OrderingTerm { return &OrderingTerm{expr: expr, nulls: "FIRST"} }
func OrderDefaultNullsLast(expr AsExpr) *OrderingTerm { return &OrderingTerm{expr: expr, nulls: "LAST"} }

type OrderingTerm struct {
	expr  AsExpr
	order string
	nulls string
}

func (t *OrderingTerm) AsOrderingTerm(s *Serializer) {
	s.F(t.expr.AsExpr).DC(" ", t.order != "").DC(t.order, t.order != "").DC(" NULLS ", t.nulls != "").DC(t.nulls, t.nulls != "")
}
