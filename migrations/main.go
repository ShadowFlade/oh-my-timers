package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofor-little/env"

	_ "database/sql"

	"github.com/pressly/goose/v3"
	// _ "modernc.org/sqlite"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()

	curDir, err := os.Getwd()
	fmt.Println("Current working directory:", curDir)

	log.Print(args)

	if len(args) < 2 {
		flags.Usage()
		return
	}

	command := args[1]
	curDir, err = filepath.Abs(curDir + "/.env")
	fmt.Println(curDir, " CURRENT DIRECOTRY")

	if err := env.Load(curDir); err != nil {
		fmt.Println("error")
		panic(err)
	}
	dbstring := env.Get("GOOSE_DBSTRING", "")

	gooseDbDriver := env.Get("GOOSE_DRIVER", "i")

	db, err := goose.OpenDBWithDriver(gooseDbDriver, dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
