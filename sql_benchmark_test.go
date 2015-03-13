package db_sql_benchmark

import (
	"database/sql"
	"log"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"

	"github.com/coopernurse/gorp"
	"github.com/jmoiron/sqlx"
)

const MysqlDSN = "root@unix(/var/run/mysqld/mysqld.sock)/golang_sql_benchmarks"

// Attempts to compare different situations regarding use of prepared statements and/or interpolation
// Use a fairly simple query, but one that allows us to use placeholders or not as per our individual benchmark
// Each query should be the same regardless of how it is generated or executed
//
// select id, subject, state from tickets where subdomain_id = 1 and (state = 'open' or state = 'spam') limit 1

//
// Benchmark different prepared statment related scenarios
//

// BenchmarkPreparedStatementsNone is the case that you have a SQL string that requires no arguments/interpolation
// The database/sql package is able to issue the query directly to the server
func BenchmarkPreparedStatementsNone(b *testing.B) {
	conn := mysqlConn()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rows, err := conn.Query("select id, subject, state from tickets where subdomain_id = 1 and (state = 'open' or state = 'spam') limit 1")
		if err != nil {
			log.Fatalln(err)
		}

		var id int
		var subject string
		var state string
		for rows.Next() {
			err := rows.Scan(&id, &subject, &state)
			if err != nil {
				log.Fatalln(err)
			}
		}

		rows.Close()
	}

	conn.Close()
}

// BenchmarkPreparedStatementsThrowaway is the (probably more likely) case that you have arguments you need to bind to placeholders in your SQL
// The database/sql package generates a prepared statement implicitly
func BenchmarkPreparedStatementsThrowaway(b *testing.B) {
	conn := mysqlConn()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rows, err := conn.Query("select id, subject, state from tickets where subdomain_id = ? and (state = ? or state = ?) limit 1", 1, "open", "spam")
		if err != nil {
			log.Fatalln(err)
		}

		var id int
		var subject string
		var state string
		for rows.Next() {
			err := rows.Scan(&id, &subject, &state)
			if err != nil {
				log.Fatalln(err)
			}
		}

		rows.Close()
	}

	conn.Close()
}

// BenchmarkPreparedStatementsSingle is the case that the user creates a single prepared statement uses it reuses it
func BenchmarkPreparedStatementsSingle(b *testing.B) {
	conn := mysqlConn()
	b.ResetTimer()

	stmt, err := conn.Prepare("select id, subject, state from tickets where subdomain_id = ? and (state = ? or state = ?) limit 1")
	if err != nil {
		log.Fatalln(err)
	}
	for n := 0; n < b.N; n++ {
		rows, err := stmt.Query(1, "open", "spam")
		if err != nil {
			log.Fatalln(err)
		}

		var id int
		var subject string
		var state string
		for rows.Next() {
			err := rows.Scan(&id, &subject, &state)
			if err != nil {
				log.Fatalln(err)
			}
		}

		rows.Close()
	}

	conn.Close()
}

//
// dbr
//

// Select into integers
func BenchmarkDbrSelectInts1(b *testing.B) {
	benchmarkDbrSelectInts(b, 1)
}

func BenchmarkDbrSelectInts100(b *testing.B) {
	benchmarkDbrSelectInts(b, 100)
}

func BenchmarkDbrSelectInts1000(b *testing.B) {
	benchmarkDbrSelectInts(b, 1000)
}

func BenchmarkDbrSelectInts10000(b *testing.B) {
	benchmarkDbrSelectInts(b, 10000)
}

// Select without query params
func BenchmarkDbrSelectAll1(b *testing.B) {
	benchmarkDbrSelectAll(b, 1)
}

func BenchmarkDbrSelectAll100(b *testing.B) {
	benchmarkDbrSelectAll(b, 100)
}

func BenchmarkDbrSelectAll1000(b *testing.B) {
	benchmarkDbrSelectAll(b, 1000)
}

func BenchmarkDbrSelectAll10000(b *testing.B) {
	benchmarkDbrSelectAll(b, 10000)
}

// Select with query params
func BenchmarkDbrSelectAllWithArgs1(b *testing.B) {
	benchmarkDbrSelectAllWithArgs(b, 1)
}

func BenchmarkDbrSelectAllWithArgs100(b *testing.B) {
	benchmarkDbrSelectAllWithArgs(b, 100)
}

func BenchmarkDbrSelectAllWithArgs1000(b *testing.B) {
	benchmarkDbrSelectAllWithArgs(b, 1000)
}

func BenchmarkDbrSelectAllWithArgs10000(b *testing.B) {
	benchmarkDbrSelectAllWithArgs(b, 10000)
}

//
// sqlx
//

// Select into integers
func BenchmarkSqlxSelectInts1(b *testing.B) {
	benchmarkSqlxSelectInts(b, 1)
}

func BenchmarkSqlxSelectInts100(b *testing.B) {
	benchmarkSqlxSelectInts(b, 100)
}

func BenchmarkSqlxSelectInts1000(b *testing.B) {
	benchmarkSqlxSelectInts(b, 1000)
}

func BenchmarkSqlxSelectInts10000(b *testing.B) {
	benchmarkSqlxSelectInts(b, 10000)
}

// Select without query params
func BenchmarkSqlxSelectAll1(b *testing.B) {
	benchmarkSqlxSelectAll(b, 1)
}

func BenchmarkSqlxSelectAll100(b *testing.B) {
	benchmarkSqlxSelectAll(b, 100)
}

func BenchmarkSqlxSelectAll1000(b *testing.B) {
	benchmarkSqlxSelectAll(b, 1000)
}

func BenchmarkSqlxSelectAll10000(b *testing.B) {
	benchmarkSqlxSelectAll(b, 10000)
}

// Select with query params
func BenchmarkSqlxSelectAllWithArgs1(b *testing.B) {
	benchmarkSqlxSelectAllWithArgs(b, 1)
}

func BenchmarkSqlxSelectAllWithArgs100(b *testing.B) {
	benchmarkSqlxSelectAllWithArgs(b, 100)
}

func BenchmarkSqlxSelectAllWithArgs1000(b *testing.B) {
	benchmarkSqlxSelectAllWithArgs(b, 1000)
}

func BenchmarkSqlxSelectAllWithArgs10000(b *testing.B) {
	benchmarkSqlxSelectAllWithArgs(b, 10000)
}

//
// Gorp
//

// Select into integer; Gorp only seems to support this for single ints, not slices
func BenchmarkGorpSelectInts1(b *testing.B) {
	dbmap := gorpConn()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, err := dbmap.SelectInt("select id from tickets where subdomain_id = 1 limit ?", 1)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// Select without query params
func BenchmarkGorpSelectAll1(b *testing.B) {
	benchmarkGorpSelectAll(b, 1)
}

func BenchmarkGorpSelectAll100(b *testing.B) {
	benchmarkGorpSelectAll(b, 100)
}

func BenchmarkGorpSelectAll1000(b *testing.B) {
	benchmarkGorpSelectAll(b, 1000)
}

func BenchmarkGorpSelectAll10000(b *testing.B) {
	benchmarkGorpSelectAll(b, 10000)
}

// Select with query params
func BenchmarkGorpSelectAllWithArgs1(b *testing.B) {
	benchmarkGorpSelectAllWithArgs(b, 1)
}

func BenchmarkGorpSelectAllWithArgs100(b *testing.B) {
	benchmarkGorpSelectAllWithArgs(b, 100)
}

func BenchmarkGorpSelectAllWithArgs1000(b *testing.B) {
	benchmarkGorpSelectAllWithArgs(b, 1000)
}

func BenchmarkGorpSelectAllWithArgs10000(b *testing.B) {
	benchmarkGorpSelectAllWithArgs(b, 10000)
}

//
// Benchmark helpers
//

func benchmarkDbrSelectAll(b *testing.B, limit int) {
	sess := dbrSess()
	query := "select id, subject, state from tickets where subdomain_id = 1 and (state = 'open' or state = 'spam') limit " + strconv.Itoa(limit)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tickets := []*struct {
			Id      int64
			Subject string
			State   string
		}{}
		// _, err := sess.Select("id", "subject", "state").From("tickets").Where("subdomain_id = 1 and (state = 'open' or state = 'spam')").Limit(limit).LoadStructs(&tickets)
		_, err := sess.SelectBySql(query).LoadStructs(&tickets)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkDbrSelectAllWithArgs(b *testing.B, limit int) {
	sess := dbrSess()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tickets := []*struct {
			Id      int64
			Subject string
			State   string
		}{}
		// _, err := sess.Select("id", "subject", "state").From("tickets").Where("subdomain_id = ? and (state = ? or state = ?)", 1, "open", "spam").Limit(limit).LoadStructs(&tickets)
		_, err := sess.SelectBySql("select id, subject, state from tickets where subdomain_id = ? and (state = ? or state = ?) limit ?", 1, "open", "spam", limit).LoadStructs(&tickets)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkDbrSelectInts(b *testing.B, limit int) {
	sess := dbrSess()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var counts []int
		_, err := sess.SelectBySql("select id from tickets where subdomain_id = 1 limit ?", limit).LoadValues(&counts)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkSqlxSelectAll(b *testing.B, limit int) {
	db := sqlxConn()
	query := "select id, subject, state from tickets where subdomain_id = 1 and (state = 'open' or state = 'spam') limit " + strconv.Itoa(limit)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tickets := []*struct {
			Id      int64
			Subject string
			State   string
		}{}

		err := db.Select(&tickets, query)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkSqlxSelectAllWithArgs(b *testing.B, limit int) {
	db := sqlxConn()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tickets := []*struct {
			Id      int64
			Subject string
			State   string
		}{}

		err := db.Select(&tickets, "select id, subject, state from tickets where subdomain_id = ? and (state = ? or state = ?) limit ?", 1, "open", "spam", limit)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkSqlxSelectInts(b *testing.B, limit int) {
	db := sqlxConn()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var counts []int
		err := db.Select(&counts, "select id from tickets where subdomain_id = 1 limit ?", limit)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkGorpSelectAll(b *testing.B, limit int) {
	dbmap := gorpConn()
	query := "select id, subject, state from tickets where subdomain_id = 1 and (state = 'open' or state = 'spam') limit " + strconv.Itoa(limit)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tickets := []*struct {
			Id      int64
			Subject string
			State   string
		}{}

		_, err := dbmap.Select(&tickets, query)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func benchmarkGorpSelectAllWithArgs(b *testing.B, limit int) {
	dbmap := gorpConn()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tickets := []*struct {
			Id      int64
			Subject string
			State   string
		}{}

		_, err := dbmap.Select(&tickets, "select id, subject, state from tickets where subdomain_id = ? and (state = ? or state = ?) limit ?", 1, "open", "spam", limit)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

//
// Connection helpers
//

func mysqlConn() *sql.DB {
	db, err := sql.Open("mysql", MysqlDSN)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func dbrSess() *dbr.Session {
	return dbr.NewConnection(mysqlConn(), nil).NewSession(nil)
}

func sqlxConn() *sqlx.DB {
	db, err := sqlx.Connect("mysql", MysqlDSN)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func gorpConn() *gorp.DbMap {
	db, err := sql.Open("mysql", MysqlDSN)
	if err != nil {
		log.Fatalln(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDb", "UTF8"}}
	return dbmap
}
