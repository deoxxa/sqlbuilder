package sqlbuilder

type CaseExpr struct {
	inputExpr AsExpr
	cases     []*CaseWhenClause
	elseExpr  AsExpr
}

func Case(inputExpr AsExpr) *CaseExpr {
	return &CaseExpr{inputExpr: inputExpr}
}

func (c *CaseExpr) When(whenClause ...*CaseWhenClause) *CaseExpr {
	return &CaseExpr{inputExpr: c.inputExpr, cases: append(c.cases[:], whenClause...), elseExpr: c.elseExpr}
}

func (c *CaseExpr) Else(elseExpr AsExpr) *CaseExpr {
	return &CaseExpr{inputExpr: c.inputExpr, cases: c.cases, elseExpr: elseExpr}
}

func (c *CaseExpr) AsExpr(s *Serializer) {
	s.D("CASE ")

	if c.inputExpr != nil {
		s.F(c.inputExpr.AsExpr).D(" ")
	}

	for _, e := range c.cases {
		s.D("WHEN ").F(e.whenExpr.AsExpr).D(" THEN ").F(e.resultExpr.AsExpr).D(" ")
	}

	if c.elseExpr != nil {
		s.D("ELSE ").F(c.elseExpr.AsExpr).D(" ")
	}

	s.D("END")
}

type CaseWhenClause struct {
	whenExpr, resultExpr AsExpr
}

func CaseWhen(whenExpr, resultExpr AsExpr) *CaseWhenClause {
	return &CaseWhenClause{whenExpr: whenExpr, resultExpr: resultExpr}
}
