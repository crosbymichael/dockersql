package main

import (
	"database/sql"
	"time"

	"github.com/samalba/dockerclient"
)

func loadImages(client *dockerclient.DockerClient, db *sql.DB) error {
	images, err := client.ListImages()
	if err != nil {
		return err
	}
	if _, err := db.Exec(
		"CREATE TABLE images (id, parent_id, size, virtual_size, tag, created)"); err != nil {
		return err
	}
	for _, i := range images {
		created := time.Unix(i.Created, 0)
		if _, err := db.Exec(
			"INSERT INTO images (id, parent_id, size, virtual_size, tag, created) VALUES (?, ?, ?, ?, ?, ?)",
			i.Id, i.ParentId, i.Size, i.VirtualSize, i.RepoTags[0], created); err != nil {
			return err
		}
	}
	return nil
}
