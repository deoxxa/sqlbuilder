package sqlbuilder

type BasicColumn struct {
	schema AsNamed
	table  AsNamed
	name   string
}

func (c *BasicColumn) AsNamedShort(s *Serializer) {
	s.N(c.name)
}

func (c *BasicColumn) AsExpr(s *Serializer) {
	if c.table != nil {
		if c.schema != nil {
			s.F(c.schema.AsNamed).D(".")
		}

		s.F(c.table.AsNamed).D(".")
	}

	s.N(c.name)
}

func (c *BasicColumn) As(alias string) *ColumnAlias {
	return AliasColumn(c, alias)
}

type ColumnAlias struct {
	expr  AsExpr
	alias string
}

func AliasColumn(expr AsExpr, alias string) *ColumnAlias {
	return &ColumnAlias{expr: expr, alias: alias}
}

func (e *ColumnAlias) As(alias string) *ColumnAlias {
	return &ColumnAlias{expr: e.expr, alias: alias}
}

func (e *ColumnAlias) AsExpr(s *Serializer) {
	if e.alias != "" {
		s.N(e.alias)
	} else {
		s.F(e.expr.AsExpr)
	}
}

func (e *ColumnAlias) AsResultColumn(s *Serializer) {
	s.F(e.expr.AsExpr)

	if e.alias != "" {
		s.D(" AS ").N(e.alias)
	}
}
