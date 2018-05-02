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
    GetConn() *sql.DB
}

var (
	drivers []string = []string{"mysql", "sqlserver", "postgres"}
	dbs     map[string]dbsetup
)

func BenchmarkMysqlQuery(b *testing.B) {
	if db, ok := dbs["mysql"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            var id int
            c.QueryRow("select id from test where id=1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkMysqlPrepStmtInnerLoop(b *testing.B) {
	if db, ok := dbs["mysql"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        stmt, err := c.Prepare("select id from test where id=?")
        defer stmt.Close()
        if err != nil {
            log.Fatal(err)
        }
        for i := 0; i < b.N; i++ {
            var id int
            c.QueryRow("1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkMysqlPrepStmtOuterLoop(b *testing.B) {
	if db, ok := dbs["mysql"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            var id int
            stmt, err := c.Prepare("select id from test where id=?")
            defer stmt.Close()
            if err != nil {
                log.Fatal(err)
            }
            c.QueryRow("1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkPostgresQuery(b *testing.B) {
	if db, ok := dbs["postgres"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            var id int
            c.QueryRow("select id from test where id=1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkPostgresPrepStmtInnerLoop(b *testing.B) {
	if db, ok := dbs["postgres"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        stmt, err := c.Prepare("select id from test where id=$1")
        defer stmt.Close()
        if err != nil {
            log.Fatal(err)
        }
        for i := 0; i < b.N; i++ {
            var id int
            c.QueryRow("1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkPostgresPrepStmtOuterLoop(b *testing.B) {
	if db, ok := dbs["postgres"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            var id int
            stmt, err := c.Prepare("select id from test where id=$1")
            defer stmt.Close()
            if err != nil {
                log.Fatal(err)
            }
            c.QueryRow("1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkSqlserverQuery(b *testing.B) {
	if db, ok := dbs["sqlserver"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            var id int
            c.QueryRow("select id from test where id=1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkSqlserverPrepStmtInnerLoop(b *testing.B) {
	if db, ok := dbs["sqlserver"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        stmt, err := c.Prepare("select id from test where id=?")
        defer stmt.Close()
        if err != nil {
            log.Fatal(err)
        }
        for i := 0; i < b.N; i++ {
            var id int
            c.QueryRow("1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }

}

func BenchmarkSqlserverPrepStmtOuterLoop(b *testing.B) {
	if db, ok := dbs["sqlserver"]; ok {
        c := db.GetConn()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            var id int
            stmt, err := c.Prepare("select id from test where id=?")
            defer stmt.Close()
            if err != nil {
                log.Fatal(err)
            }
            c.QueryRow("1").Scan(&id)
        }
    } else {
        b.Errorf("Did not run")
    }
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	dbs = make(map[string]dbsetup)
	for _, d := range drivers {
		var db dbsetup = &connection{}
        if err := db.Setup(d); err == nil {
            defer db.Teardown()
            dbs[d] = db
        }
	}
	return m.Run()
}
