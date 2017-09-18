package sqlbuilder

func LeftJoin(left, right AsTableOrSubquery) *JoinExpr  { return Join("LEFT", left, right) }
func CrossJoin(left, right AsTableOrSubquery) *JoinExpr { return Join("CROSS", left, right) }

type JoinExpr struct {
	kind        string
	left, right AsTableOrSubquery
	condition   AsExpr
}

func Join(kind string, left, right AsTableOrSubquery) *JoinExpr {
	return &JoinExpr{kind: kind, left: left, right: right}
}

func (j *JoinExpr) On(condition AsExpr) *JoinExpr {
	return &JoinExpr{kind: j.kind, left: j.left, right: j.right, condition: condition}
}

func (j *JoinExpr) AsTableOrSubquery(s *Serializer) {
	s.F(j.left.AsTableOrSubquery)

	s.DC(" "+j.kind, j.kind != "").D(" JOIN ")

	s.F(j.right.AsTableOrSubquery)

	if j.condition != nil {
		s.D(" ON ").F(j.condition.AsExpr)
	}
}

func (j *JoinExpr) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, j, right)
}

func (j *JoinExpr) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(j, right)
}

func (j *JoinExpr) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(j, right)
}
