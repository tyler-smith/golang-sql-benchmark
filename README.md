golang-db-sql-benchmark
====================

A collection of benchmarks for popular Go database/SQL utilities

# Libraries under test

*  [database/sql](https://golang.org/pkg/database/sql/) + [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
*  [gocraft/dbr](https://github.com/gocraft/dbr)
*  [Gorp](https://github.com/coopernurse/gorp)
*  [Sqlx](https://github.com/jmoiron/sqlx)
*  [Squirrel](https://github.com/lann/squirrel)

# database/sql SQL Execution Benchmarks:

* **BenchmarkPreparedStatementsNone** - Runs simple queries without query arguments, so database/sql doesn't need to create a prepared statement
* **BenchmarkPreparedStatementsThrowaway** - Runs queries with query arguments. database/sql must create and then throwaway a prepared statement each time
* **BenchmarkPreparedStatementsSingle** - Runs queries with query arguments, but creates and reuses the a single prepared statement


# Dbr/Sqlx/Gorp SQL Execution Benchmarks:
Each library under test has the same set of benchmarks, just replace `Dbr` in the examples with `Sqlx` or `Gorp`.
Each one is run with varying number of rows, N.

* **BenchmarkDbrSelectIntsN** - Select rows of integers into []int64's
* **BenchmarkDbrSelectAllN** - Select rows into structs using no query arguments
* **BenchmarkDbrSelectAllWithArgsN** - Select rows into structs using query arguments

# Dbr/Squrrel SQL Building Benchamrks
Test building (but not executing) various SQL statements

* **BenchmarkBuilderDbrSimple** - Simple SQL query with dbr
* **BenchmarkBuilderDbrComplex** - Complex SQL query with dbr
* **BenchmarkBuilderSquirrelSimple** - Simple SQL query with squirrel
* **BenchmarkBuilderSquirrelComplex** - Complex SQL query with squirrel

# Output

`godep go test -bench=. -benchmem 2>/dev/null | column -t` on @tyler-smith's 2.6 GHz i7 Macbook Pro:

```
BenchmarkPreparedStatementsNone       100000       23137     ns/op  612      B/op  20     allocs/op
BenchmarkPreparedStatementsThrowaway  20000        94086     ns/op  795      B/op  25     allocs/op
BenchmarkPreparedStatementsSingle     50000        33426     ns/op  652      B/op  23     allocs/op

BenchmarkDbrSelectInts1               100000       30098     ns/op  1198     B/op  22     allocs/op
BenchmarkDbrSelectInts100             20000        94477     ns/op  11645    B/op  426    allocs/op
BenchmarkDbrSelectInts1000            5000         627958    ns/op  100574   B/op  4041   allocs/op
BenchmarkDbrSelectInts10000           500          6207828   ns/op  1297929  B/op  40177  allocs/op
BenchmarkDbrSelectAll1                50000        35468     ns/op  1876     B/op  60     allocs/op
BenchmarkDbrSelectAll100              10000        157510    ns/op  29256    B/op  863    allocs/op
BenchmarkDbrSelectAll1000             2000         1436856   ns/op  272365   B/op  8078   allocs/op
BenchmarkDbrSelectAll10000            100          12927094  ns/op  2979829  B/op  80171  allocs/op
BenchmarkDbrSelectAllWithArgs1        50000        38959     ns/op  2590     B/op  72     allocs/op
BenchmarkDbrSelectAllWithArgs100      10000        162227    ns/op  29970    B/op  875    allocs/op
BenchmarkDbrSelectAllWithArgs1000     2000         1227154   ns/op  273765   B/op  8093   allocs/op
BenchmarkDbrSelectAllWithArgs10000    100          13342026  ns/op  3000164  B/op  80284  allocs/op

BenchmarkSqlxSelectInts1              20000        86128     ns/op  628      B/op  22     allocs/op
BenchmarkSqlxSelectInts100            10000        203886    ns/op  9480     B/op  426    allocs/op
BenchmarkSqlxSelectInts1000           2000         1289621   ns/op  83799    B/op  4038   allocs/op
BenchmarkSqlxSelectInts10000          100          15156243  ns/op  1134724  B/op  40139  allocs/op
BenchmarkSqlxSelectAll1               50000        32373     ns/op  959      B/op  25     allocs/op
BenchmarkSqlxSelectAll100             10000        177965    ns/op  28366    B/op  828    allocs/op
BenchmarkSqlxSelectAll1000            1000         1634705   ns/op  271352   B/op  8042   allocs/op
BenchmarkSqlxSelectAll10000           100          17483420  ns/op  2989883  B/op  80191  allocs/op
BenchmarkSqlxSelectAllWithArgs1       10000        140460    ns/op  1228     B/op  30     allocs/op
BenchmarkSqlxSelectAllWithArgs100     10000        250621    ns/op  23682    B/op  634    allocs/op
BenchmarkSqlxSelectAllWithArgs1000    2000         1548073   ns/op  229879   B/op  6049   allocs/op
BenchmarkSqlxSelectAllWithArgs10000   100          15441239  ns/op  2587009  B/op  60190  allocs/op

BenchmarkGorpSelectInts1              10000        134277    ns/op  432      B/op  15     allocs/op
BenchmarkGorpSelectAll1               50000        44086     ns/op  1983     B/op  75     allocs/op
BenchmarkGorpSelectAll100             10000        232527    ns/op  34258    B/op  978    allocs/op
BenchmarkGorpSelectAll1000            1000         1563675   ns/op  321846   B/op  9100   allocs/op
BenchmarkGorpSelectAll10000           100          16451807  ns/op  3472508  B/op  90250  allocs/op
BenchmarkGorpSelectAllWithArgs1       10000        131738    ns/op  2306     B/op  80     allocs/op
BenchmarkGorpSelectAllWithArgs100     10000        265323    ns/op  29948    B/op  785    allocs/op
BenchmarkGorpSelectAllWithArgs1000    2000         1423747   ns/op  279838   B/op  7103   allocs/op
BenchmarkGorpSelectAllWithArgs10000   100          14538568  ns/op  3072519  B/op  70263  allocs/op

BenchmarkBuilderDbrSimple             1000000      1699      ns/op  866      B/op  13     allocs/op
BenchmarkBuilderDbrComplex            500000       6193      ns/op  2190     B/op  37     allocs/op

BenchmarkBuilderSquirrelSimple        200000       8981      ns/op  2780     B/op  51     allocs/op
BenchmarkBuilderSquirrelComplex       50000        44721     ns/op  11707    B/op  259    allocs/op
```

# Run yourself

* Use the `godep` tool or manually install all libraries under test
* Create db: `mysql -e "create database golang_sql_benchmarks;"`
* Create schema: `mysql golang_sql_benchmarks < structure.sql`
* Run: `godep go test -bench=. -benchmem`
* You can set the MySQL DSN to use by setting the GOLANG_SQL_BENCHMARKS_DSN env var (defaults to root@unix(/var/run/mysqld/mysqld.sock)/golang_sql_benchmarks)
