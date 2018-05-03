package sqlquery

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

type dbsetup interface {
	Setup(driver string) error
	Teardown() error
	Conn() *sql.DB
}

var (
	drivers []string = []string{"mysql", "sqlserver", "postgres"}
	dbs     map[string]dbsetup
)

func BenchmarkMysqlQuery(b *testing.B) {
	db, ok := dbs["mysql"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		c.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkMysqlPrepStmtInnerLoop(b *testing.B) {
	db, ok := dbs["mysql"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	stmt, err := c.Prepare("select id from test where id=?")
	defer stmt.Close()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var id int
		c.QueryRow("1").Scan(&id)
	}
}

func BenchmarkMysqlPrepStmtOuterLoop(b *testing.B) {
	db, ok := dbs["mysql"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			var id int
			stmt, err := c.Prepare("select id from test where id=?")
			defer stmt.Close()
			if err != nil {
				b.Fatal(err)
			}
			c.QueryRow("1").Scan(&id)
		}()
	}
}

func BenchmarkPostgresQuery(b *testing.B) {
	db, ok := dbs["postgres"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		c.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkPostgresPrepStmtInnerLoop(b *testing.B) {
	db, ok := dbs["postgres"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	stmt, err := c.Prepare("select id from test where id=$1")
	defer stmt.Close()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var id int
		c.QueryRow("1").Scan(&id)
	}
}

func BenchmarkPostgresPrepStmtOuterLoop(b *testing.B) {
	db, ok := dbs["postgres"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			var id int
			stmt, err := c.Prepare("select id from test where id=$1")
			defer stmt.Close()
			if err != nil {
				b.Fatal(err)
			}
			c.QueryRow("1").Scan(&id)
		}()
	}
}

func BenchmarkSqlserverQuery(b *testing.B) {
	db, ok := dbs["sqlserver"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		c.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkSqlserverPrepStmtInnerLoop(b *testing.B) {
	db, ok := dbs["sqlserver"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	stmt, err := c.Prepare("select id from test where id=?")
	defer stmt.Close()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		var id int
		c.QueryRow("1").Scan(&id)
	}
}

func BenchmarkSqlserverPrepStmtOuterLoop(b *testing.B) {
	db, ok := dbs["sqlserver"]
	if !ok {
		b.Fatalf("Did not run")
	}
	c := db.Conn()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			var id int
			stmt, err := c.Prepare("select id from test where id=?")
			defer stmt.Close()
			if err != nil {
				b.Fatal(err)
			}
			c.QueryRow("1").Scan(&id)
		}()
	}
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	dbs = make(map[string]dbsetup)
	for _, d := range drivers {
		var db dbsetup = &connection{}
		err := db.Setup(d)
		if err != nil {
			log.Println(err)
			continue
		}
		defer db.Teardown()
		dbs[d] = db
	}
	return m.Run()
}
