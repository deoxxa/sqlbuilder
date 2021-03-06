package sqlbuilder

import (
	"strings"
)

type Serializer struct {
	d Dialect

	bits []string
	vals []*BoundVariable
	vpos map[*BoundVariable]int
}

func NewSerializer(d Dialect) *Serializer {
	return &Serializer{
		d:    d,
		vpos: make(map[*BoundVariable]int),
	}
}

type BoundVariable struct {
	value interface{}
}

func Bind(value interface{}) *BoundVariable {
	return &BoundVariable{value: value}
}

// BindAllAsExpr binds multiple variables at once, returning them as []AsExpr.
// This is most useful for In and NotIn conditions.
func BindAllAsExpr(vals ...interface{}) []AsExpr {
	l := make([]AsExpr, len(vals))

	for i, val := range vals {
		l[i] = Bind(val)
	}

	return l
}

// BindAllStringsAsExpr binds multiple variables at once, returning them as
// []AsExpr. This is most useful for In and NotIn conditions.
func BindAllStringsAsExpr(vals ...string) []AsExpr {
	l := make([]AsExpr, len(vals))

	for i, val := range vals {
		l[i] = Bind(val)
	}

	return l
}

// BindAllIntsAsExpr binds multiple variables at once, returning them as
// []AsExpr. This is most useful for In and NotIn conditions.
func BindAllIntsAsExpr(vals ...int) []AsExpr {
	l := make([]AsExpr, len(vals))

	for i, val := range vals {
		l[i] = Bind(val)
	}

	return l
}

func (b *BoundVariable) AsExpr(s *Serializer) {
	s.V(b)
}

func (b *BoundVariable) As(alias string) *ColumnAlias {
	return AliasColumn(b, alias)
}

// Bind binds a variable to this serializer - this can be used for e.g.
// symbolic representation of user input.
func (s *Serializer) Bind(val interface{}) *BoundVariable {
	b, ok := val.(*BoundVariable)
	if !ok {
		b = Bind(val)
	}

	if _, ok := s.vpos[b]; !ok {
		s.vals = append(s.vals, b)
		s.vpos[b] = len(s.vals)
	}

	return b
}

// BindAllAsExpr binds multiple variables to this serializer at once,
// returning them as []AsExpr. This is most useful for In and NotIn
// conditions.
func (s *Serializer) BindAllAsExpr(vals ...interface{}) []AsExpr {
	l := make([]AsExpr, len(vals))

	for i, val := range vals {
		l[i] = s.Bind(val)
	}

	return l
}

// BindAllStringsAsExpr binds multiple variables to this serializer at once,
// returning them as []AsExpr. This is most useful for In and NotIn
// conditions.
func (s *Serializer) BindAllStringsAsExpr(vals ...string) []AsExpr {
	l := make([]AsExpr, len(vals))

	for i, val := range vals {
		l[i] = s.Bind(val)
	}

	return l
}

// BindAllIntsAsExpr binds multiple variables to this serializer at once,
// returning them as []AsExpr. This is most useful for In and NotIn
// conditions.
func (s *Serializer) BindAllIntsAsExpr(vals ...int) []AsExpr {
	l := make([]AsExpr, len(vals))

	for i, val := range vals {
		l[i] = s.Bind(val)
	}

	return l
}

// SetDialect sets the dialect for this serializer
func (s *Serializer) SetDialect(d Dialect) {
	s.d = d
}

// ToSQL serializes the whole query, returning the query itself and any
// variables requred to execute it
func (s *Serializer) ToSQL() (string, []interface{}, error) {
	var vars []interface{}
	for _, v := range s.vals {
		vars = append(vars, v.value)
	}

	return strings.Join(s.bits, ""), vars, nil
}

// N adds a "name" value to the query, which should be quoted as per the
// dialect
func (s *Serializer) N(name string) *Serializer { return s.NC(name, true) }

// NC adds a "name" value to the query, which should be quoted as per the
// dialect, only if "w" is true
func (s *Serializer) NC(name string, w bool) *Serializer {
	if w {
		s.bits = append(s.bits, dialect(s.d).QuoteName(name))
	}
	return s
}

// D adds string data to the query
func (s *Serializer) D(data string) *Serializer { return s.DC(data, true) }

// DC adds string data to the query, only if "w" is true
func (s *Serializer) DC(data string, w bool) *Serializer {
	if w {
		s.bits = append(s.bits, data)
	}
	return s
}

// V adds a value to the query, where the value should be replaced with a
// placeholder
func (s *Serializer) V(val interface{}) *Serializer { return s.VC(val, true) }

// VC adds a value to the query, where the value should be replaced with a
// placeholder, only if "w" is true
func (s *Serializer) VC(val interface{}, w bool) *Serializer {
	if w {
		s.bits = append(s.bits, dialect(s.d).Bind(s.vpos[s.Bind(val)]))
	}
	return s
}

// F runs a function taking the Serializer as its only argument
func (s *Serializer) F(f func(s *Serializer)) *Serializer { return s.FC(f, true) }

// FC runs a function taking the Serializer as its only argument, only if "w"
// is true
func (s *Serializer) FC(f func(s *Serializer), w bool) *Serializer {
	if w {
		f(s)
	}
	return s
}
