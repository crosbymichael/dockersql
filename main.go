package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samalba/dockerclient"
)

var (
	logger      = logrus.New()
	globalFlags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "enabled debug output for the logs"},
		cli.StringFlag{Name: "docker", Value: "unix:///var/run/docker.sock", Usage: "url to your docker daemon endpoint"},
	}
)

func preload(context *cli.Context) error {
	if context.GlobalBool("debug") {
		logger.Level = logrus.DebugLevel
	}
	return nil
}

func mainAction(context *cli.Context) {
	client, err := dockerclient.NewDockerClient(context.GlobalString("docker"), nil)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		logger.Fatal(err)
	}

	if err := loadContainers(client, db); err != nil {
		db.Close()
		logger.Fatal(err)
	}
	if err := loadImages(client, db); err != nil {
		db.Close()
		logger.Fatal(err)
	}

	s := bufio.NewScanner(os.Stdin)
	for {
		prompt()
		if !s.Scan() {
			break
		}

		query := s.Text()
		if query == "" {
			continue
		}
		rows, err := db.Query(query)
		if err != nil {
			logger.Warn(err)
			continue
		}
		if err := DisplayResults(rows); err != nil {
			db.Close()
			logger.Fatal(err)
		}
	}
	db.Close()
}

func prompt() {
	fmt.Fprintf(os.Stdout, "> ")
}

func main() {
	app := cli.NewApp()
	app.Name = "dockersql"
	app.Author = "@crosbymichael"
	app.Usage = "query your dockers with SQL"
	app.Flags = globalFlags
	app.Before = preload
	app.Action = mainAction

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
