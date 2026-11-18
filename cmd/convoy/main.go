package main

import (
	"database/sql"
	"fmt"
	"log"

	. "convoy/pkg"

	flag "github.com/spf13/pflag"
)

const version = "0.1.0"

var Port = flag.IntP("port", "p", 8080, "Port to listen on.")

func printHelp() {
	fmt.Printf("CONVOY message queue        v%s\n\nUsage:\n\tconvoy [OPTIONS]\n\nOptions:\n", version)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printHelp
	flag.Parse()
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	List, err = LoadDB(db)
	if err != nil {
		log.Fatal(err)
	}
	wait := StartServer(*Port, db)
	stopped := false
	for !stopped {
		select {
		case <-wait:
			stopped = true
		}
	}
}

var waitingWorkers []WorkerProcess

type WorkerProcess struct {
	id string
}
