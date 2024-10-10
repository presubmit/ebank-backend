package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gobike/envflag"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	dir    = flag.String("dir", "db/migrations", "Directory in which to save migrations")
	dbHost = flag.String("db_host", "localhost", "database host")
	dbPort = flag.String("db_port", "5432", "database port")
	dbPass = flag.String("db_pass", "ebank123", "database password")
	dbUser = flag.String("db_user", "postgres", "database username")
	dbName = flag.String("db_name", "ebank", "database name")
)

func connectDB() (*migrate.Migrate, error) {
	envflag.Parse()

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", *dbUser, *dbPass, *dbHost, *dbPort, *dbName)
	println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		fmt.Println("not here", err)
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("here", err)
		return nil, err
	}
	return migrate.NewWithDatabaseInstance(
		"file://"+*dir,
		"postgres", driver)
}

func createFile(filename string, body string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(body); err != nil {
		return err
	}

	return f.Close()
}

func main() {
	_ = flag.Set("alsologtostderr", "true")
	flag.Parse()

	println("⏱  Running migrations...")

	cmd := "up"
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	switch cmd {
	case "clean":
		m, err := connectDB()
		if err != nil {
			panic(err)
		}
		if err = m.Drop(); err != nil {
			fmt.Println(err)
		}

	case "up":
		num := 0
		if len(os.Args) > 3 {
			var err error
			num, err = strconv.Atoi(os.Args[3])
			if err != nil || num < 1 {
				panic("Invalid number of migrations: " + err.Error())
			}
		}

		m, err := connectDB()
		if err != nil {
			panic(fmt.Sprintf("Error connecting to db: %v", err))
		}

		if num != 0 {
			if err = m.Steps(num); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = m.Up(); err != nil {
				fmt.Println(err)
			}
		}

	case "down":
		num := 0
		if len(os.Args) > 3 {
			println(os.Args[3])
			var err error
			num, err = strconv.Atoi(os.Args[3])
			if err != nil || num < 1 {
				panic("Invalid number of migrations: " + err.Error())
			}
		}

		m, err := connectDB()
		if err != nil {
			fmt.Println(err)
		}

		if num != 0 {
			if err = m.Steps(-num); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = m.Down(); err != nil {
				panic(err)
			}
		}

	case "force":
		if len(os.Args) < 4 {
			panic("Invalid migration to force")
		}
		num, err := strconv.Atoi(os.Args[3])
		if err != nil {
			panic("Invalid migration to force: " + err.Error())
		}

		m, err := connectDB()
		if err != nil {
			panic(err)
		}
		if err := m.Force(num); err != nil {
			fmt.Println(err)
		}

	case "create":
		if len(os.Args) < 4 {
			panic("Invalid name as the second argument")
		}
		name := os.Args[3]
		version := strconv.FormatInt(time.Now().Unix(), 10)
		dir := filepath.Clean(*dir)

		for _, direction := range []string{"up", "down"} {
			basename := fmt.Sprintf("%s_%s.%s.sql", version, name, direction)
			filename := filepath.Join(dir, basename)

			template := "BEGIN;\n\n"
			if direction == "up" {
				template += "-- TODO: write migration sql"
			} else {
				template += "-- TODO: write rollback sql"
			}
			template += "\n\nCOMMIT;"

			if err := createFile(filename, template); err != nil {
				panic(err)
			}

			absPath, _ := filepath.Abs(filename)
			fmt.Println("Generated: " + absPath)
		}
	}
}
