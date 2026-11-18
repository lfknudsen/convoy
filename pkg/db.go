package pkg

import (
	"container/list"
	"database/sql"
	"errors"
	"os"
	"sync"
	"time"

	. "convoy/proto"

	"modernc.org/sqlite"
)

var List = list.New()

var lock sync.Mutex

func push(db *sql.DB, url string) error {
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

func pop(db *sql.DB) (*Link, error) {
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

func OpenDB() (*sql.DB, error) {
	err := os.MkdirAll(ConvoyDir, 0644)
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
