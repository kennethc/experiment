package sqlquery

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"testing"
)

type connection struct {
	dsn  string
	conn *sql.DB
}

var (
	dbTypes []string = []string{"mysql", "sqlserver", "postgres"}
	dbs     map[string]connection
)

func BenchmarkMysqlQuery(b *testing.B) {
	db := dbs["mysql"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No MySQL connection.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkMysqlPrepStmtInnerLoop(b *testing.B) {
	db := dbs["mysql"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No MySQL connection.")
	}
	b.ResetTimer()
	stmt, err := db.Prepare("select id from test where id=?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("1").Scan(&id)
	}
}

func BenchmarkMysqlPrepStmtOuterLoop(b *testing.B) {
	db := dbs["mysql"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No MySQL connection.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		stmt, err := db.Prepare("select id from test where id=?")
		defer stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
		db.QueryRow("1").Scan(&id)
	}
}

func BenchmarkPostgresQuery(b *testing.B) {
	db := dbs["postgres"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No PostgreSQL connection.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkPostgresPrepStmtInnerLoop(b *testing.B) {
	db := dbs["postgres"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No PostgreSQL connection.")
	}
	b.ResetTimer()
	stmt, err := db.Prepare("select id from test where id=$1")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("1").Scan(&id)
	}
}

func BenchmarkPostgresPrepStmtOuterLoop(b *testing.B) {
	db := dbs["postgres"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No PostgreSQL connection.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		stmt, err := db.Prepare("select id from test where id=$1")
		defer stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
		db.QueryRow("1").Scan(&id)
	}
}

func BenchmarkSqlserverQuery(b *testing.B) {
	db := dbs["sqlserver"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No SqlServer connection.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkSqlserverPrepStmtInnerLoop(b *testing.B) {
	db := dbs["sqlserver"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No SqlServer connection.")
	}
	b.ResetTimer()
	stmt, err := db.Prepare("select id from test where id=?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("1").Scan(&id)
	}
}

func BenchmarkSqlserverPrepStmtOuterLoop(b *testing.B) {
	db := dbs["sqlserver"].conn
	err := db.Ping()
	if err != nil {
		b.Errorf("No SqlServer connection.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		stmt, err := db.Prepare("select id from test where id=?")
		defer stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
		db.QueryRow("1").Scan(&id)
	}
}

func TestMain(m *testing.M) {
	dbs = make(map[string]connection)
	for _, v := range dbTypes {
		db := connection{}
		db.dsn = readDBCreds(v)
		db.conn = setupDB(v, db.dsn)
		dbs[v] = db
	}
	result := m.Run()
	for _, v := range dbTypes {
		teardownDB(dbs[v].conn)
	}
	os.Exit(result)
}

func readDBCreds(dbType string) string {
	return os.Getenv("DSN_" + strings.ToUpper(dbType))
}

func setupDB(dbType string, dsn string) *sql.DB {
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table test (id int)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("insert into test values (1)")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func teardownDB(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("drop table test")
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}
