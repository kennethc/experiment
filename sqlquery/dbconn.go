package sqlquery

import (
	"database/sql"
	"os"
	"strings"
)

type connection struct {
	driver string
	dsn    string
	conn   *sql.DB
}

func (c *connection) readCred() {
	c.dsn = os.Getenv("DSN_" + strings.ToUpper(c.driver))
}

func (c *connection) open() error {
	var err error
	if c.conn, err = sql.Open(c.driver, c.dsn); err != nil {
		return err
	}
	return c.conn.Ping()
}

func (c *connection) Setup(driver string) error {
	c.driver = driver
	c.readCred()
	var err error
	if err = c.open(); err != nil {
		return err
	}
	_, err = c.conn.Exec("create table test (id int)")
	if err != nil {
		return err
	}
	_, err = c.conn.Exec("insert into test values (1)")
	return err
}

func (c *connection) Teardown() error {
	_, err := c.conn.Exec("drop table test")
	if err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *connection) Conn() *sql.DB {
	return c.conn
}
