package main

import (
	"database/sql"

	"github.com/samalba/dockerclient"
)

func loadImages(client *dockerclient.DockerClient, db *sql.DB) error {
	images, err := client.ListImages()
	if err != nil {
		return err
	}
	if _, err := db.Exec(
		"CREATE TABLE images (id, parent_id, size, virtual_size, tag)"); err != nil {
		return err
	}
	for _, i := range images {
		if _, err := db.Exec(
			"INSERT INTO images (id, parent_id, size, virtual_size, tag) VALUES (?, ?, ?, ?, ?)",
			i.Id, i.ParentId, i.Size, i.VirtualSize, i.RepoTags[0]); err != nil {
			return err
		}
	}
	return nil
}
