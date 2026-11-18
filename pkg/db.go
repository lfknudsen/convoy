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
	statement, err := tx.Prepare("INSERT INTO queue(version, added, message) VALUES(?,?,?)")
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

type QueueElement struct {
	version int
	added   int64
	url     string
}

func pop(db *sql.DB) (*Link, error) {
	lock.Lock()
	defer lock.Unlock()
	front := List.Front()
	if front == nil {
		return nil, errors.New("list is empty")
	}

	// Retrieve the matching element in the database, so we can persist its
	// column values in subsequent tables.
	rows, err := db.Query("SELECT version, added, message FROM queue WHERE message=?", front.Value.(string))
	if err != nil {
		return nil, err
	}
	var element QueueElement
	if !rows.Next() {
		_ = rows.Close()
		return nil, errors.New("database queue is empty; linked list and database are inconsistent")
	}
	err = rows.Scan(&element.version, &element.added, &element.url)
	if err != nil {
		_ = rows.Close()
		return nil, err
	}
	_ = rows.Close()

	// Move element from `queue` to `in_progress`
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	statement, err := tx.Prepare("DELETE FROM queue WHERE message=?")
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

	statement, err = tx.Prepare(`
		INSERT INTO in_progress (version, added, started, message)
		VALUES(?,?,?,?)`)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_, err = statement.Exec(element.version, element.added, time.Now().Unix(), element.url)
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

func LoadDB(db *sql.DB) (*list.List, error) {
	lock.Lock()
	defer lock.Unlock()
	rows, err := db.Query("SELECT message FROM queue ORDER BY added")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	l := list.New()
	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		_ = l.PushBack(url)
	}
	err = rows.Err()
	return l, err
}
