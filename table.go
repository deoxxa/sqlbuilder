package sqlbuilder

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

	return nil
}

func (t *Table) AsNamed(s *Serializer) {
	s.N(t.name)
}

func (t *Table) AsTableOrSubquery(s *Serializer) {
	s.N(t.name)
}

type TableAlias struct {
	from    *Table
	name    string
	columns []*BasicColumn
}

func AliasTable(t *Table, name string) *TableAlias {
	a := &TableAlias{
		from:    t,
		name:    name,
		columns: make([]*BasicColumn, len(t.columns)),
	}

	for i, c := range t.columns {
		a.columns[i] = &BasicColumn{table: a, name: c.name}
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
	c := a.from.C(name)

	return &BasicColumn{
		table: a,
		name:  c.name,
	}
}
