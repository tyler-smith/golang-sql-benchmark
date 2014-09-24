package db_sql_benchmark

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/lann/squirrel"
)

//
// Build, but don't execute, queries
//

//
// dbr
//

func BenchmarkBuilderDbrSimple(b *testing.B) {
	sess := dbrSess()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sess.Select("id").From("tickets").Where("subdomain_id = ? and (state = ? or state = ?)", 1, "open", "spam").ToSql()
	}
}

func BenchmarkBuilderDbrComplex(b *testing.B) {
	sess := dbrSess()

	arg_eq1 := dbr.Eq{"f": 2, "x": "hi"}
	arg_eq2 := map[string]interface{}{"g": 3}
	arg_eq3 := dbr.Eq{"h": []int{1, 2, 3}}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sess.Select("a", "b", "z", "y", "x").
			Distinct().
			From("c").
			Where("d = ? OR e = ?", 1, "wat").
			Where(arg_eq1).
			Where(arg_eq2).
			Where(arg_eq3).
			GroupBy("i").
			GroupBy("ii").
			GroupBy("iii").
			Having("j = k").
			Having("jj = ?", 1).
			Having("jjj = ?", 2).
			OrderBy("l").
			OrderBy("l").
			OrderBy("l").
			Limit(7).
			Offset(8).
			ToSql()
	}
}

//
// Squirrel
//

func BenchmarkBuilderSquirrelSimple(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		squirrel.Select("id").From("tickets").Where("subdomain_id = ? and (state = ? or state = ?)", 1, "open", "spam").ToSql()
	}
}

func BenchmarkBuilderSquirrelComplex(b *testing.B) {
	arg_eq1 := squirrel.Eq{"f": 2, "x": "hi"}
	arg_eq2 := map[string]interface{}{"g": 3}
	arg_eq3 := squirrel.Eq{"h": []int{1, 2, 3}}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		squirrel.Select("a", "b", "z", "y", "x").
			Distinct().
			From("c").
			Where("d = ? OR e = ?", 1, "wat").
			Where(arg_eq1).
			Where(arg_eq2).
			Where(arg_eq3).
			GroupBy("i").
			GroupBy("ii").
			GroupBy("iii").
			Having("j = k").
			Having("jj = ?", 1).
			Having("jjj = ?", 2).
			OrderBy("l").
			OrderBy("l").
			OrderBy("l").
			Limit(7).
			Offset(8).
			ToSql()
	}
}
