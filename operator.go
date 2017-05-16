package sqlbuilder

func In(expr AsExpr, v ...AsExpr) *InOperatorExpr    { return InOperator("IN", expr, v...) }
func NotIn(expr AsExpr, v ...AsExpr) *InOperatorExpr { return InOperator("NOT IN", expr, v...) }

type InOperatorExpr struct {
	operator string
	expr     AsExpr
	elements []AsExpr
}

func InOperator(operator string, expr AsExpr, v ...AsExpr) *InOperatorExpr {
	return &InOperatorExpr{operator: operator, expr: expr, elements: v}
}

func (o *InOperatorExpr) AsExpr(s *Serializer) {
	s.F(o.expr.AsExpr).D(" ").D(o.operator).D(" ")

	s.D("(")

	for i, e := range o.elements {
		s.F(e.AsExpr).DC(", ", i < len(o.elements)-1)
	}

	s.D(")")
}

func Like(left, right AsExpr) *LikeOperatorExpr {
	return &LikeOperatorExpr{left: left, right: right}
}

type LikeOperatorExpr struct {
	left, right AsExpr
}

func (o *LikeOperatorExpr) AsExpr(s *Serializer) {
	s.F(o.left.AsExpr).D(" LIKE ").F(o.right.AsExpr)
}

func And(left, right AsExpr) *BinaryOperatorExpr { return BinaryOperator("AND", left, right) }
func Or(left, right AsExpr) *BinaryOperatorExpr  { return BinaryOperator("OR", left, right) }

type BooleanOperatorExpr struct {
	operator string
	elements []AsExpr
}

func BooleanOperator(operator string, elements ...AsExpr) *BooleanOperatorExpr {
	return &BooleanOperatorExpr{operator: operator, elements: elements}
}

func (b *BooleanOperatorExpr) AsExpr(s *Serializer) {
	if len(b.elements) == 1 {
		s.F(b.elements[0].AsExpr)
		return
	}

	s.D("(")

	for i, e := range b.elements {
		s.F(e.AsExpr).DC(" "+b.operator+" ", i < len(b.elements)-1)
	}

	s.D(")")
}

func Eq(left, right AsExpr) *BinaryOperatorExpr  { return BinaryOperator("=", left, right) }
func Ne(left, right AsExpr) *BinaryOperatorExpr  { return BinaryOperator("!=", left, right) }
func Gt(left, right AsExpr) *BinaryOperatorExpr  { return BinaryOperator(">", left, right) }
func Lt(left, right AsExpr) *BinaryOperatorExpr  { return BinaryOperator("<", left, right) }
func Gte(left, right AsExpr) *BinaryOperatorExpr { return BinaryOperator(">=", left, right) }
func Lte(left, right AsExpr) *BinaryOperatorExpr { return BinaryOperator("<=", left, right) }

type BinaryOperatorExpr struct {
	operator    string
	left, right AsExpr
}

func BinaryOperator(operator string, left, right AsExpr) *BinaryOperatorExpr {
	return &BinaryOperatorExpr{operator: operator, left: left, right: right}
}

func (o *BinaryOperatorExpr) AsExpr(s *Serializer) {
	s.D("(").F(o.left.AsExpr).D(" " + o.operator + " ").F(o.right.AsExpr).D(")")
}

func Not(expr AsExpr) *PrefixOperatorExpr { return PrefixOperator("NOT", expr) }

type PrefixOperatorExpr struct {
	operator string
	expr     AsExpr
}

func PrefixOperator(operator string, expr AsExpr) *PrefixOperatorExpr {
	return &PrefixOperatorExpr{operator: operator, expr: expr}
}

func (o *PrefixOperatorExpr) AsExpr(s *Serializer) {
	s.D(o.operator).D(" ").F(o.expr.AsExpr)
}

func IsNull(expr AsExpr) *PostfixOperatorExpr    { return PostfixOperator(expr, "IS NULL") }
func IsNotNull(expr AsExpr) *PostfixOperatorExpr { return PostfixOperator(expr, "IS NOT NULL") }

type PostfixOperatorExpr struct {
	expr     AsExpr
	operator string
}

func PostfixOperator(expr AsExpr, operator string) *PostfixOperatorExpr {
	return &PostfixOperatorExpr{expr: expr, operator: operator}
}

func (o *PostfixOperatorExpr) AsExpr(s *Serializer) {
	s.F(o.expr.AsExpr).D(" ").D(o.operator)
}
