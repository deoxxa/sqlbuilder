package sqlbuilder

func OffsetLimit(offset, limit AsExpr) *OffsetLimitClause {
	return &OffsetLimitClause{offset: offset, limit: limit}
}

type OffsetLimitClause struct {
	offset, limit AsExpr
}

func (c *OffsetLimitClause) AsOffsetLimit(s *Serializer) {
	if c.limit != nil {
		s.D("LIMIT ").F(c.limit.AsExpr)

		if c.offset != nil {
			s.D(" OFFSET ").F(c.offset.AsExpr)
		}
	}
}
