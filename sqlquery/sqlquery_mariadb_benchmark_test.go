package sqlquery

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var (
	user   string
	passwd string
	host   string
	dbname string
	dsn    string
	db     *sql.DB
)

func BenchmarkSQLQuery(b *testing.B) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		db.QueryRow("select id from test where id=1").Scan(&id)
	}
}

func BenchmarkSQLPrepStmtInnerLoop(b *testing.B) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
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

func BenchmarkSQLPrepStmtOuterLoop(b *testing.B) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
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
	readDBCreds()
	setupDB()
	result := m.Run()
	teardownDB()
	os.Exit(result)
}

func readDBCreds() {
	user = os.Getenv("DB_USER")
	passwd = os.Getenv("DB_PASSWD")
	host = os.Getenv("DB_HOST")
	dbname = os.Getenv("DB_NAME")
	dsn = user + ":" + passwd + "@" + host + "/" + dbname
}

func setupDB() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("create table if not exists test (id int)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("insert into test values (1)")
	if err != nil {
		log.Fatal(err)
	}
}

func teardownDB() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("drop table test")
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}
