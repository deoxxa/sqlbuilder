package sqlbuilder

func Literal(text string) *LiteralExpr {
	return &LiteralExpr{text: text}
}

type LiteralExpr struct {
	text string
}

func (l *LiteralExpr) AsExpr(s *Serializer) {
	s.D(l.text)
}

func (e *LiteralExpr) As(alias string) *ColumnAlias {
	return AliasColumn(e, alias)
}

func TypedLiteral(typ, text string) *TypedLiteralExpr {
	return &TypedLiteralExpr{typ: typ, text: text}
}

type TypedLiteralExpr struct {
	typ  string
	text string
}

func (l *TypedLiteralExpr) AsExpr(s *Serializer) {
	s.D(l.typ).D(" ").D(l.text)
}

func (e *TypedLiteralExpr) As(alias string) *ColumnAlias {
	return AliasColumn(e, alias)
}
