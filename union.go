package sqlbuilder

type UnionExpr struct {
	all  bool
	stmt *SelectStatement
}

func Union(all bool, stmt *SelectStatement) *UnionExpr {
	return &UnionExpr{all: all, stmt: stmt}
}

func (e *UnionExpr) AsUnion(s *Serializer) {
	s.D("UNION ").DC("ALL ", e.all).F(e.stmt.AsStatement)
}
