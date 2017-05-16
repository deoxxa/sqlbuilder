package sqlbuilder

import (
	"regexp"
	"strconv"
	"strings"
)

type Dialect interface {
	Bind(i int) string
	QuoteName(name string) string
}

func dialect(d Dialect) Dialect {
	if d != nil {
		return d
	}

	return DialectGeneric{}
}

var dialectGenericNameRegexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z_0-9]*$")

var dialectGenericReservedWords = []string{
	"ABSOLUTE", "EXEC", "OVERLAPS", "ACTION", "EXECUTE", "PAD", "ADA", "EXISTS",
	"PARTIAL", "ADD", "EXTERNAL", "PASCAL", "ALL", "EXTRACT", "POSITION",
	"ALLOCATE", "FALSE", "PRECISION", "ALTER", "FETCH", "PREPARE", "AND", "FIRST",
	"PRESERVE", "ANY", "FLOAT", "PRIMARY", "ARE", "FOR", "PRIOR", "AS", "FOREIGN",
	"PRIVILEGES", "ASC", "FORTRAN", "PROCEDURE", "ASSERTION", "FOUND", "PUBLIC",
	"AT", "FROM", "READ", "AUTHORIZATION", "FULL", "REAL", "AVG", "GET",
	"REFERENCES", "BEGIN", "GLOBAL", "RELATIVE", "BETWEEN", "GO", "RESTRICT",
	"BIT", "GOTO", "REVOKE", "BIT_LENGTH", "GRANT", "RIGHT", "BOTH", "GROUP",
	"ROLLBACK", "BY", "HAVING", "ROWS", "CASCADE", "HOUR", "SCHEMA", "CASCADED",
	"IDENTITY", "SCROLL", "CASE", "IMMEDIATE", "SECOND", "CAST", "IN", "SECTION",
	"CATALOG", "INCLUDE", "SELECT", "CHAR", "INDEX", "SESSION", "CHAR_LENGTH",
	"INDICATOR", "SESSION_USER", "CHARACTER", "INITIALLY", "SET",
	"CHARACTER_LENGTH", "INNER", "SIZE", "CHECK", "INPUT", "SMALLINT", "CLOSE",
	"INSENSITIVE", "SOME", "COALESCE", "INSERT", "SPACE", "COLLATE", "INT", "SQL",
	"COLLATION", "INTEGER", "SQLCA", "COLUMN", "INTERSECT", "SQLCODE", "COMMIT",
	"INTERVAL", "SQLERROR", "CONNECT", "INTO", "SQLSTATE", "CONNECTION", "IS",
	"SQLWARNING", "CONSTRAINT", "ISOLATION", "SUBSTRING", "CONSTRAINTS", "JOIN",
	"SUM", "CONTINUE", "KEY", "SYSTEM_USER", "CONVERT", "LANGUAGE", "TABLE",
	"CORRESPONDING", "LAST", "TEMPORARY", "COUNT", "LEADING", "THEN", "CREATE",
	"LEFT", "TIME", "CROSS", "LEVEL", "TIMESTAMP", "CURRENT", "LIKE",
	"TIMEZONE_HOUR", "CURRENT_DATE", "LOCAL", "TIMEZONE_MINUTE", "CURRENT_TIME",
	"LOWER", "TO", "CURRENT_TIMESTAMP", "MATCH", "TRAILING", "CURRENT_USER",
	"MAX", "TRANSACTION", "CURSOR", "MIN", "TRANSLATE", "DATE", "MINUTE",
	"TRANSLATION", "DAY", "MODULE", "TRIM", "DEALLOCATE", "MONTH", "TRUE", "DEC",
	"NAMES", "UNION", "DECIMAL", "NATIONAL", "UNIQUE", "DECLARE", "NATURAL",
	"UNKNOWN", "DEFAULT", "NCHAR", "UPDATE", "DEFERRABLE", "NEXT", "UPPER",
	"DEFERRED", "NO", "USAGE", "DELETE", "NONE", "USER", "DESC", "NOT", "USING",
	"DESCRIBE", "NULL", "VALUE", "DESCRIPTOR", "NULLIF", "VALUES", "DIAGNOSTICS",
	"NUMERIC", "VARCHAR", "DISCONNECT", "OCTET_LENGTH", "VARYING", "DISTINCT",
	"OF", "VIEW", "DOMAIN", "ON", "WHEN", "DOUBLE", "ONLY", "WHENEVER", "DROP",
	"OPEN", "WHERE", "ELSE", "OPTION", "WITH", "END", "OR", "WORK", "END-EXEC",
	"ORDER", "WRITE", "ESCAPE", "OUTER", "YEAR", "EXCEPT", "OUTPUT", "ZONE",
	"EXCEPTION",
}

type DialectGeneric struct{}

func (DialectGeneric) Bind(i int) string { return "$" + strconv.FormatInt(int64(i), 10) }

func (DialectGeneric) QuoteName(name string) string {
	if !dialectGenericNameRegexp.MatchString(name) {
		return strconv.Quote(name)
	}

	for _, e := range dialectGenericReservedWords {
		if strings.EqualFold(e, name) {
			return strconv.Quote(name)
		}
	}

	return name
}

var dialectMSSQLReservedWords = []string{
	"ADD", "EXTERNAL", "PROCEDURE", "ALL", "FETCH", "PUBLIC", "ALTER", "FILE",
	"RAISERROR", "AND", "FILLFACTOR", "READ", "ANY", "FOR", "READTEXT", "AS",
	"FOREIGN", "RECONFIGURE", "ASC", "FREETEXT", "REFERENCES", "AUTHORIZATION",
	"FREETEXTTABLE", "REPLICATION", "BACKUP", "FROM", "RESTORE", "BEGIN",
	"FULL", "RESTRICT", "BETWEEN", "FUNCTION", "RETURN", "BREAK", "GOTO",
	"REVERT", "BROWSE", "GRANT", "REVOKE", "BULK", "GROUP", "RIGHT", "BY",
	"HAVING", "ROLLBACK", "CASCADE", "HOLDLOCK", "ROWCOUNT", "CASE", "IDENTITY",
	"ROWGUIDCOL", "CHECK", "IDENTITY_INSERT", "RULE", "CHECKPOINT",
	"IDENTITYCOL", "SAVE", "CLOSE", "IF", "SCHEMA", "CLUSTERED", "IN",
	"SECURITYAUDIT", "COALESCE", "INDEX", "SELECT", "COLLATE", "INNER",
	"SEMANTICKEYPHRASETABLE", "COLUMN", "INSERT",
	"SEMANTICSIMILARITYDETAILSTABLE", "COMMIT", "INTERSECT",
	"SEMANTICSIMILARITYTABLE", "COMPUTE", "INTO", "SESSION_USER", "CONSTRAINT",
	"IS", "SET", "CONTAINS", "JOIN", "SETUSER", "CONTAINSTABLE", "KEY",
	"SHUTDOWN", "CONTINUE", "KILL", "SOME", "CONVERT", "LEFT", "STATISTICS",
	"CREATE", "LIKE", "SYSTEM_USER", "CROSS", "LINENO", "TABLE", "CURRENT",
	"LOAD", "TABLESAMPLE", "CURRENT_DATE", "MERGE", "TEXTSIZE", "CURRENT_TIME",
	"NATIONAL", "THEN", "CURRENT_TIMESTAMP", "NOCHECK", "TO", "CURRENT_USER",
	"NONCLUSTERED", "TOP", "CURSOR", "NOT", "TRAN", "DATABASE", "NULL",
	"TRANSACTION", "DBCC", "NULLIF", "TRIGGER", "DEALLOCATE", "OF", "TRUNCATE",
	"DECLARE", "OFF", "TRY_CONVERT", "DEFAULT", "OFFSETS", "TSEQUAL", "DELETE",
	"ON", "UNION", "DENY", "OPEN", "UNIQUE", "DESC", "OPENDATASOURCE",
	"UNPIVOT", "DISK", "OPENQUERY", "UPDATE", "DISTINCT", "OPENROWSET",
	"UPDATETEXT", "DISTRIBUTED", "OPENXML", "USE", "DOUBLE", "OPTION", "USER",
	"DROP", "OR", "VALUES", "DUMP", "ORDER", "VARYING", "ELSE", "OUTER", "VIEW",
	"END", "OVER", "WAITFOR", "ERRLVL", "PERCENT", "WHEN", "ESCAPE", "PIVOT",
	"WHERE", "EXCEPT", "PLAN", "WHILE", "EXEC", "PRECISION", "WITH", "EXECUTE",
	"PRIMARY", "WITHIN GROUP", "EXISTS", "PRINT", "WRITETEXT", "EXIT", "PROC",
}

type DialectMSSQL struct{ QuoteEverything bool }

func (DialectMSSQL) Bind(i int) string { return "@p" + strconv.FormatInt(int64(i), 10) }

func (d DialectMSSQL) QuoteName(name string) string {
	if d.QuoteEverything {
		return "[" + name + "]"
	}

	if !dialectGenericNameRegexp.MatchString(name) {
		return "[" + name + "]"
	}

	for _, e := range dialectGenericReservedWords {
		if strings.EqualFold(e, name) {
			return "[" + name + "]"
		}
	}

	for _, e := range dialectMSSQLReservedWords {
		if strings.EqualFold(e, name) {
			return "[" + name + "]"
		}
	}

	return name
}

func MSSQLOffsetLimit(offset, limit AsExpr) *MSSQLOffsetLimitClause {
	return &MSSQLOffsetLimitClause{offset: offset, limit: limit}
}

type MSSQLOffsetLimitClause struct {
	offset, limit AsExpr
}

func (c *MSSQLOffsetLimitClause) AsOffsetLimit(s *Serializer) {
	if c.offset != nil {
		s.D("OFFSET ").F(c.offset.AsExpr).D(" ROWS")
	}

	if c.limit != nil {
		s.DC(" ", c.offset != nil).D("FETCH NEXT ").F(c.limit.AsExpr).D(" ROWS ONLY")
	}
}

type DialectPostgres struct{}

func (DialectPostgres) Bind(i int) string { return "$" + strconv.FormatInt(int64(i), 10) }

func (DialectPostgres) QuoteName(name string) string {
	return strconv.Quote(name)
}

type DialectSQLite struct{}

func (DialectSQLite) Bind(i int) string { return "$" + strconv.FormatInt(int64(i), 10) }

func (DialectSQLite) QuoteName(name string) string {
	return strconv.Quote(name)
}
