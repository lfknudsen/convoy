package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func getExecPath() string {
	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(p)
}

var ConvoyDir = filepath.Join(getExecPath(), ".convoy")
var databaseName = "convoy.db"

var connectionString = fmt.Sprintf("file:%s/%s", ConvoyDir, databaseName)
