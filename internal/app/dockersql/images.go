package dockersql

import (
	"database/sql"
	"github.com/samalba/dockerclient"
	"time"
)

func LoadImages(client *dockerclient.DockerClient, db *sql.DB) error {
	images, err := client.ListImages(true)
	if err != nil {
		return err
	}
	if _, err := db.Exec(
		"CREATE TABLE images (id, parent_id, size, virtual_size, tag, created)"); err != nil {
		return err
	}
	for _, i := range images {
		created := time.Unix(i.Created, 0)
		var tag = ""
		if i.RepoTags != nil {
			tag = i.RepoTags[0]
		}
		if _, err := db.Exec(
			"INSERT INTO images (id, parent_id, size, virtual_size, tag, created) VALUES (?, ?, ?, ?, ?, ?)",
			i.Id, i.ParentId, i.Size, i.VirtualSize, tag, created); err != nil {
			return err
		}
	}
	return nil
}
