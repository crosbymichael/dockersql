package main

import (
	"database/sql"
	"dockersql/internal/app/dockersql"
	ln "github.com/GeertJohan/go.linenoise"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samalba/dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"strings"
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

func loadDatabase(context *cli.Context) (*sql.DB, error) {
	client, err := dockerclient.NewDockerClient(context.GlobalString("docker"), nil)
	if err != nil {
		logger.Fatal(err)
	}
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err := dockersql.LoadContainers(client, db); err != nil {
		db.Close()
		return nil, err
	}
	if err := dockersql.LoadImages(client, db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

var completions = []string{
	"select",
	"where",
	"from",
	"containers",
	"images",
}

// only complete on the last full word of the completion
func completion(input string) []string {
	var (
		out   []string
		parts = strings.Split(input, " ")
		lower = strings.ToLower(parts[len(parts)-1])
		l     = len(lower)
	)
	for _, c := range completions {
		if len(c) < l {
			continue
		}
		if strings.HasPrefix(c, lower) {
			if len(parts) == 1 {
				out = append(out, c)
			} else {
				parts[len(parts)-1] = c
				out = append(out, strings.Join(parts, " "))
			}
		}
	}
	return out
}

func mainAction(context *cli.Context) {
	db, err := loadDatabase(context)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	ln.SetMultiline(true)
	ln.SetCompletionHandler(completion)
	for {
		query, err := ln.Line("> ")
		if err != nil {
			if err != ln.KillSignalError {
				logger.Error(err)
			}
			return
		}
		if query == "" {
			continue
		}
		if err := ln.AddHistory(query); err != nil {
			logger.Error(err)
		}
		rows, err := db.Query(query)
		if err != nil {
			logger.Warn(err)
			continue
		}
		if err := dockersql.DisplayResults(rows); err != nil {
			db.Close()
			logger.Fatal(err)
		}
	}
}

func prompt() {
	ln.Line("> ")
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
