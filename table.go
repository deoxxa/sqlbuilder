package sqlbuilder

import "fmt"

type Table struct {
	name    string
	columns []*BasicColumn
}

func NewTable(name string, columns ...string) *Table {
	t := &Table{name: name}

	for _, c := range columns {
		t.columns = append(t.columns, &BasicColumn{
			table: t,
			name:  c,
		})
	}

	return t
}

func (t *Table) C(name string) *BasicColumn {
	for _, c := range t.columns {
		if c.name == name {
			return c
		}
	}

	panic(fmt.Errorf("Requested column %s does not exist in table %s", name, t.name))

	return nil
}

func (t *Table) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, t, right)
}

func (t *Table) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(t, right)
}

func (t *Table) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(t, right)
}

func (t *Table) AsNamed(s *Serializer) {
	s.N(t.name)
}

func (t *Table) AsTableOrSubquery(s *Serializer) {
	s.N(t.name)
}

type AsTableOrSubqueryAndWithColumns interface {
	AsTableOrSubquery
	WithColumns
}

type TableAlias struct {
	from    AsTableOrSubqueryAndWithColumns
	name    string
	columns map[string]*BasicColumn
}

func AliasTable(t AsTableOrSubqueryAndWithColumns, name string) *TableAlias {
	a := &TableAlias{
		from:    t,
		name:    name,
		columns: make(map[string]*BasicColumn),
	}

	return a
}

func (a *TableAlias) AsNamed(s *Serializer) {
	s.N(a.name)
}

func (a *TableAlias) AsTableOrSubquery(s *Serializer) {
	s.F(a.from.AsTableOrSubquery).D(" ").N(a.name)
}

func (a *TableAlias) C(name string) *BasicColumn {
	if a.columns[name] != nil {
		return a.columns[name]
	}

	if c := a.from.C(name); c != nil {
		a.columns[name] = &BasicColumn{table: a, name: name}
	}

	return a.columns[name]
}

func (a *TableAlias) Join(kind string, right AsTableOrSubquery) *JoinExpr {
	return Join(kind, a, right)
}

func (a *TableAlias) LeftJoin(right AsTableOrSubquery) *JoinExpr {
	return LeftJoin(a, right)
}

func (a *TableAlias) CrossJoin(right AsTableOrSubquery) *JoinExpr {
	return CrossJoin(a, right)
}
