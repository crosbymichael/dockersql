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

const (
	insertContainer = "INSERT INTO containers (id, image, name, status, command) VALUES (?, ?, ?, ? ,?)"
	insertImage     = "INSERT INTO images (id, parent_id, size, virtual_size, tag) VALUES (?, ?, ?, ?, ?)"
)

var (
	logger      = logrus.New()
	globalFlags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "enabled debug output for the logs"},
		cli.StringFlag{Name: "docker", Value: "unix:///var/run/docker.sock", Usage: "url to your docker daemon endpoint"},
	}
	tables = []string{
		"CREATE TABLE containers (id, image, name, status, command)",
		"CREATE TABLE images (id, parent_id, size, virtual_size, tag)",
	}
)

func preload(context *cli.Context) error {
	if context.GlobalBool("debug") {
		logger.Level = logrus.DebugLevel
	}
	return nil
}

func loadDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}

func loadContainers(client *dockerclient.DockerClient, db *sql.DB) error {
	containers, err := client.ListContainers(true)
	if err != nil {
		return err
	}
	for _, c := range containers {
		if _, err := db.Exec(insertContainer, c.Id, c.Image, c.Names[0], c.Status, c.Command); err != nil {
			return err
		}
	}
	return nil
}

func loadImages(client *dockerclient.DockerClient, db *sql.DB) error {
	images, err := client.ListImages()
	if err != nil {
		return err
	}
	for _, i := range images {
		if _, err := db.Exec(insertImage, i.Id, i.ParentId, i.Size, i.VirtualSize, i.RepoTags[0]); err != nil {
			return err
		}
	}
	return nil
}

func mainAction(context *cli.Context) {
	client, err := dockerclient.NewDockerClient(context.GlobalString("docker"), nil)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := loadDatabase()
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

		rows, err := db.Query(s.Text())
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
