package sqlbuilder

func Func(name string, args ...AsExpr) *FuncExpr {
	return &FuncExpr{name: name, args: args}
}

type FuncExpr struct {
	name string
	args []AsExpr
}

func (e *FuncExpr) AsExpr(s *Serializer) {
	s.D(e.name).D("(")

	for i, arg := range e.args {
		s.F(arg.AsExpr).DC(", ", i != len(e.args)-1)
	}

	s.D(")")
}

func (e *FuncExpr) As(alias string) *ColumnAlias {
	return AliasColumn(e, alias)
}

// FuncTableExpr is deprecated in favour of FuncTableWithoutNameExpr
type FuncTableExpr struct {
	expr   *FuncExpr
	name   string
	column *BasicColumn
}

// FuncTable is deprecated in favour of FuncTableWithoutName
func FuncTable(expr *FuncExpr, name string) *FuncTableExpr {
	t := &FuncTableExpr{
		expr: expr,
		name: name,
	}

	t.column = &BasicColumn{name: name}

	return t
}

func (t *FuncTableExpr) AsNamed(s *Serializer) {
	s.N(t.name)
}

func (t *FuncTableExpr) AsTableOrSubquery(s *Serializer) {
	s.F(t.expr.AsExpr).D(" ").N(t.name)
}

func (t *FuncTableExpr) C(name string) *BasicColumn {
	if t.name != name {
		return nil
	}

	return t.column
}

func (t *FuncTableExpr) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, t, right)
}

func (t *FuncTableExpr) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(t, right)
}

func (t *FuncTableExpr) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(t, right)
}

type FuncTableWithoutNameExpr struct {
	expr    *FuncExpr
	names   []string
	columns []*BasicColumn
}

func FuncTableWithoutName(expr *FuncExpr, names ...string) *FuncTableWithoutNameExpr {
	t := &FuncTableWithoutNameExpr{expr: expr, names: names}

	columns := make([]*BasicColumn, len(names))
	for i, v := range names {
		columns[i] = &BasicColumn{name: v}
	}
	t.columns = columns

	return t
}

func (t *FuncTableWithoutNameExpr) AsTableOrSubquery(s *Serializer) {
	s.F(t.expr.AsExpr)
}

func (t *FuncTableWithoutNameExpr) C(name string) *BasicColumn {
	for i, n := range t.names {
		if n == name {
			return t.columns[i]
		}
	}

	return nil
}

func (t *FuncTableWithoutNameExpr) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, t, right)
}

func (t *FuncTableWithoutNameExpr) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(t, right)
}

func (t *FuncTableWithoutNameExpr) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(t, right)
}
