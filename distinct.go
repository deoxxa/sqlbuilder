package sqlbuilder

type DistinctExpr struct {
	on []AsExpr
}

func Distinct(on ...AsExpr) *DistinctExpr {
	return &DistinctExpr{on: on}
}

func (e *DistinctExpr) AsDistinct(s *Serializer) {
	s.D("DISTINCT")

	if len(e.on) > 0 {
		s.D(" ON (")

		for i, v := range e.on {
			s.F(v.AsExpr).DC(", ", i == len(e.on)-1)
		}

		s.D(")")
	}
}

type MSSQLDistinctExpr struct{}

func MSSQLDistinct() *MSSQLDistinctExpr { return &MSSQLDistinctExpr{} }

func (e *MSSQLDistinctExpr) AsDistinct(s *Serializer) {
	s.D("DISTINCT")
}
