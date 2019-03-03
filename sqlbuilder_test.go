package sqlbuilder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	tbl := NewTable("users", "id", "name", "email")

	q := Update().Table(tbl).Set(UpdateColumns{
		tbl.C("name"): Bind("jim"),
	}).Where(Eq(tbl.C("id"), Bind(5)))

	qs, qv, err := NewSerializer(DialectPostgres{}).F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal("UPDATE \"users\" SET \"name\" = $1 WHERE (\"users\".\"id\" = $2)", qs)
	a.Equal([]interface{}{"jim", 5}, qv)
}

func TestSelect(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectMSSQL{})

	tbl := AliasTable(NewTable("tblproducts",
		"PartNo",
		"Type",
		"Product",
		"Grade",
		"Coating",
		"Finish",
		"Thickness",
		"Width",
		"Length",
		"Dim1",
		"Dim2",
		"ClassFBR",
		"ClassFME",
		"ClassFSY",
		"SLOB",
	), "p")

	region := s.Bind("REGION_1")
	customer := s.Bind("CUSTOMER_1")

	q := Select().From(tbl).Columns(
		tbl.C("PartNo"),
		tbl.C("Type"),
		tbl.C("Product"),
		tbl.C("Grade"),
		tbl.C("Coating"),
		tbl.C("Finish"),
		tbl.C("Thickness"),
		tbl.C("Width"),
		tbl.C("Length"),
		tbl.C("Dim1"),
		tbl.C("Dim2"),
		tbl.C("ClassFBR"),
		tbl.C("ClassFME"),
		tbl.C("ClassFSY"),
		AliasColumn(s.Bind("D"), "ClassFHO"),
		tbl.C("SLOB"),
		AliasColumn(Func("dbo.productAvailablePlusALTOAmountOnHand", tbl.C("PartNo"), region), "OnHandAmount"),
		AliasColumn(Func("dbo.productAvailableWeightOnHand", tbl.C("PartNo"), region), "OnHandWeight"),
		AliasColumn(Func("dbo.productReservedAmount", tbl.C("PartNo"), region), "ReservedAmount"),
		AliasColumn(Func("dbo.productALTOAmount", tbl.C("PartNo"), region), "ALTO"),
		AliasColumn(Func("dbo.getPartAvgCost", tbl.C("PartNo")), "AverageCost"),
		AliasColumn(Func("dbo.productOnOrderAmount", tbl.C("PartNo"), region), "OnOrderAmount"),
		AliasColumn(Func("dbo.productOnOrderWeight", tbl.C("PartNo"), region), "OnOrderWeight"),
		AliasColumn(Func("dbo.getListPriceGivenAvgCost", customer, tbl.C("PartNo"), Func("dbo.getPartAvgCost", tbl.C("PartNo"))), "MinimumPrice"),
		AliasColumn(Func("dbo.customerLastPrice", tbl.C("PartNo"), customer), "CustomerLastPrice"),
		AliasColumn(Func("CONVERT", Func("VARCHAR", Literal("23")), Func("dbo.customerLastSoldDate", tbl.C("PartNo"), customer), Literal("126")), "CustomerLastSoldDate"),
		AliasColumn(InfixOperator("+", Literal("1"), Literal("2"), Literal("3")), "Five"),
	).Where(
		In(tbl.C("PartNo"), s.BindAllAsExpr(1000, 1001, 1002)...),
	).OrderBy(
		OrderAsc(tbl.C("Type")),
		OrderAsc(tbl.C("Product")),
		OrderAsc(tbl.C("Grade")),
		OrderAsc(tbl.C("Coating")),
		OrderAsc(tbl.C("Finish")),
		OrderAsc(tbl.C("Thickness")),
		OrderAsc(tbl.C("Width")),
		OrderAsc(tbl.C("Dim1")),
		OrderAsc(tbl.C("Dim2")),
		OrderAsc(tbl.C("Length")),
	).OffsetLimit(MSSQLOffsetLimit(s.Bind(30), s.Bind(30)))

	qs, qv, err := s.F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal("SELECT p.PartNo, p.Type, p.Product, p.Grade, p.Coating, p.Finish, p.Thickness, p.Width, p.Length, p.Dim1, p.Dim2, p.ClassFBR, p.ClassFME, p.ClassFSY, @p3 AS ClassFHO, p.SLOB, dbo.productAvailablePlusALTOAmountOnHand(p.PartNo, @p1) AS OnHandAmount, dbo.productAvailableWeightOnHand(p.PartNo, @p1) AS OnHandWeight, dbo.productReservedAmount(p.PartNo, @p1) AS ReservedAmount, dbo.productALTOAmount(p.PartNo, @p1) AS ALTO, dbo.getPartAvgCost(p.PartNo) AS AverageCost, dbo.productOnOrderAmount(p.PartNo, @p1) AS OnOrderAmount, dbo.productOnOrderWeight(p.PartNo, @p1) AS OnOrderWeight, dbo.getListPriceGivenAvgCost(@p2, p.PartNo, dbo.getPartAvgCost(p.PartNo)) AS MinimumPrice, dbo.customerLastPrice(p.PartNo, @p2) AS CustomerLastPrice, CONVERT(VARCHAR(23), dbo.customerLastSoldDate(p.PartNo, @p2), 126) AS CustomerLastSoldDate, (1 + 2 + 3) AS Five FROM tblproducts p WHERE p.PartNo IN (@p4, @p5, @p6) ORDER BY p.Type ASC, p.Product ASC, p.Grade ASC, p.Coating ASC, p.Finish ASC, p.Thickness ASC, p.Width ASC, p.Dim1 ASC, p.Dim2 ASC, p.Length ASC OFFSET @p7 ROWS FETCH NEXT @p8 ROWS ONLY", qs)
	a.Equal([]interface{}{"REGION_1", "CUSTOMER_1", "D", 1000, 1001, 1002, 30, 30}, qv)
}

func TestJoin(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectSQLite{})

	t1 := NewTable("table1", "id", "title1")
	t2 := NewTable("table2", "id", "table1_id", "title2")
	t3 := NewTable("table3", "id", "table2_id", "title3")

	q := Select().Columns(
		t1.C("title1"),
		t2.C("title2"),
		t3.C("title3"),
	).From(
		t1.
			LeftJoin(t2).On(Eq(t2.C("table1_id"), t1.C("id"))).
			LeftJoin(t3).On(Eq(t3.C("table2_id"), t2.C("id"))),
	).Where(Eq(t3.C("title3"), Bind("asdf")))

	qs, qv, err := s.F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal(`SELECT "table1"."title1", "table2"."title2", "table3"."title3" FROM "table1" LEFT JOIN "table2" ON ("table2"."table1_id" = "table1"."id") LEFT JOIN "table3" ON ("table3"."table2_id" = "table2"."id") WHERE ("table3"."title3" = $1)`, qs)
	a.Equal([]interface{}{"asdf"}, qv)
}

func TestMultipleQueries(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectSQLite{})

	t1 := NewTable("table1", "id", "title")

	q := Update().Table(t1)

	for _, e := range []struct {
		id    int
		title string
	}{
		{1, "title one"},
		{2, "title two"},
		{3, "title three"},
	} {
		s = s.
			F(q.
				Set(UpdateColumns{t1.C("title"): Bind(e.title)}).Where(Eq(t1.C("id"), Bind(e.id))).AsStatement).
			D(";")
	}

	qs, qv, err := s.ToSQL()

	a.NoError(err)
	a.Equal(`UPDATE "table1" SET "title" = $1 WHERE ("table1"."id" = $2);UPDATE "table1" SET "title" = $3 WHERE ("table1"."id" = $4);UPDATE "table1" SET "title" = $5 WHERE ("table1"."id" = $6);`, qs)
	a.Equal([]interface{}{"title one", 1, "title two", 2, "title three", 3}, qv)
}

func TestExpressionTable(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectPostgres{})

	t1 := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)

	tbl := ExpressionTable(Func("generate_series", Bind(t1), Bind(t2), Literal("'1 month'")), "dt", "d")

	q := Select().From(tbl).Columns(tbl.C("d"))

	qs, qv, err := s.F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal(`SELECT "dt"."d" FROM generate_series($1, $2, '1 month') AS "dt" ("d")`, qs)
	a.Equal([]interface{}{t1, t2}, qv)
}

func TestCast(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectPostgres{})

	q := Select().Columns(Cast(Bind(1), "numeric"))

	qs, qv, err := s.F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal(`SELECT CAST ($1 AS numeric)`, qs)
	a.Equal([]interface{}{1}, qv)
}

func TestFilter(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectPostgres{})

	tbl := NewTable("table1", "d", "n")

	q := Select().Columns(Filter(Func("sum", tbl.C("n")), Gt(tbl.C("d"), Bind(5))))

	qs, qv, err := s.F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal(`SELECT sum("table1"."n") FILTER (WHERE ("table1"."d" > $1))`, qs)
	a.Equal([]interface{}{5}, qv)
}

func TestFuncTable(t *testing.T) {
	a := assert.New(t)

	s := NewSerializer(DialectPostgres{})

	f := FuncTable(Func("generate_series", Bind(1), Bind(10)), "num")

	q := Select().Columns(f.C("num")).From(f)

	qs, qv, err := s.F(q.AsStatement).ToSQL()

	a.NoError(err)
	a.Equal(`SELECT "num" FROM generate_series($1, $2) "num"`, qs)
	a.Equal([]interface{}{1, 10}, qv)
}
