package xorm

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	testEngine *Engine
	dbType     string
	connString string

	db         = flag.String("db", "sqlite3", "the tested database")
	showSQL    = flag.Bool("show_sql", true, "show generated SQLs")
	ptrConnStr = flag.String("conn_str", "", "test database connection string")
	mapType    = flag.String("map_type", "snake", "indicate the name mapping")
	cache      = flag.Bool("cache", false, "if enable cache")
)

func createEngine(dbType, connStr string) error {
	if testEngine == nil {
		var err error
		testEngine, err = NewEngine(dbType, connStr)
		if err != nil {
			return err
		}

		testEngine.ShowSQL(*showSQL)
		testEngine.logger.SetLevel(core.LOG_DEBUG)
	}

	tables, err := testEngine.DBMetas()
	if err != nil {
		return err
	}
	var tableNames = make([]interface{}, 0, len(tables))
	for _, table := range tables {
		tableNames = append(tableNames, table.Name)
	}
	return testEngine.DropTables(tableNames...)
}

func prepareEngine() error {
	return createEngine(dbType, connString)
}

func TestMain(m *testing.M) {
	flag.Parse()

	dbType = *db
	if *db == "sqlite3" {
		if ptrConnStr == nil {
			connString = "./test.db?cache=shared&mode=rwc"
		} else {
			connString = *ptrConnStr
		}
	} else {
		if ptrConnStr == nil {
			fmt.Println("you should indicate conn string")
			return
		}
		connString = *ptrConnStr
	}

	dbs := strings.Split(*db, ";")
	conns := strings.Split(connString, ";")

	var res int
	for i := 0; i < len(dbs); i++ {
		dbType = dbs[i]
		connString = conns[i]
		testEngine = nil
		fmt.Println("testing", dbType, connString)

		if err := prepareEngine(); err != nil {
			fmt.Println(err)
			return
		}

		code := m.Run()
		if code > 0 {
			res = code
		}
	}

	os.Exit(res)
}

func TestPing(t *testing.T) {
	if err := testEngine.Ping(); err != nil {
		t.Fatal(err)
	}
}
