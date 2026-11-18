package main

import (
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	. "convoy/pkg"
	. "convoy/proto"

	flag "github.com/spf13/pflag"
	"modernc.org/sqlite"
)

const version = "0.1.0"
const entryVersion = 1

var Port = flag.IntP("port", "p", 8080, "Port to listen on.")

var List = list.New()

func getExecPath() string {
	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(p)
}

var convoyDir = filepath.Join(getExecPath(), ".convoy")
var databaseName = "convoy.db"

var connectionString = fmt.Sprintf("file:%s/%s", convoyDir, databaseName)

var db *sql.DB
var lock sync.Mutex

func printHelp() {
	fmt.Printf("CONVOY message queue        v%s\n\nUsage:\n\tconvoy [OPTIONS]\n\nOptions:\n", version)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printHelp
	flag.Parse()
	var err error
	db, err = openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	wait := StartServer(*Port)
	stopped := false
	for !stopped {
		select {
		case <-wait:
			stopped = true
		}
	}
}

func push(url string) error {
	lock.Lock()
	defer lock.Unlock()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	statement, err := tx.Prepare("INSERT INTO queue(version, added, url) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(1, time.Now().Unix(), url)
	if err != nil {
		_ = statement.Close()
		return err
	}
	err = statement.Close()
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	_ = List.PushBack(url)
	return nil
}

var waitingWorkers []WorkerProcess

type WorkerProcess struct {
	id string
}

func pop() (*Link, error) {
	lock.Lock()
	defer lock.Unlock()
	front := List.Front()
	if front == nil {
		return nil, errors.New("list is empty")
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	statement, err := tx.Prepare("DELETE FROM queue WHERE url=?")
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_, err = statement.Exec(front.Value.(*Link).Url)
	if err != nil {
		_ = tx.Rollback()
		_ = statement.Close()
		return nil, err
	}
	_ = statement.Close()

	statement, err = tx.Prepare("INSERT INTO in_progress VALUES(?,?,?,?)")
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_, err = statement.Exec(1, time.Now().Unix(), front.Value.(*Link).Url)
	if err != nil {
		_ = tx.Rollback()
		_ = statement.Close()
		return nil, err
	}
	_ = statement.Close()

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	_ = List.Remove(front)
	val := front.Value.(*Link)
	return val, nil
}

func closeStmt(stmt *sql.Stmt) {
	if stmt != nil {
		_ = stmt.Close()
	}
}

func openDB() (*sql.DB, error) {
	err := os.MkdirAll(convoyDir, 0644)
	if err != nil {
		return nil, err
	}
	conn, err := sqlite.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(conn)
	return db, db.Ping()
}
