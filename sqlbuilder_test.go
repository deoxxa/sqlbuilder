package sqlbuilder

import (
	"testing"

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
