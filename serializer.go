package sqlbuilder

import (
	"strings"
)

type Serializer struct {
	d Dialect

	bits []string
	vals []*BoundVariable
}

type serializableString string

func (s serializableString) serialize() string { return string(s) }

type BoundVariable struct {
	index int
	value interface{}
}

func (b *BoundVariable) AsExpr(s *Serializer) {
	s.V(b)
}

// Bind binds a variable to this serializer - this can be used for e.g.
// symbolic representation of user input.
func (s *Serializer) Bind(val interface{}) *BoundVariable {
	b, ok := val.(*BoundVariable)
	if ok {
		return b
	}

	b = &BoundVariable{index: len(s.vals) + 1, value: val}

	s.vals = append(s.vals, b)

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
		b := s.Bind(val)
		s.bits = append(s.bits, dialect(s.d).Bind(b.index))
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
