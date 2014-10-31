package main

import (
	"database/sql"

	"github.com/samalba/dockerclient"
)

func loadContainers(client *dockerclient.DockerClient, db *sql.DB) error {
	containers, err := client.ListContainers(true)
	if err != nil {
		return err
	}
	if _, err := db.Exec(
		"CREATE TABLE containers (id, image, name, status, command, cpu, memory, cpuset, ip)"); err != nil {
		return err
	}
	for _, c := range containers {
		info, err := client.InspectContainer(c.Id)
		if err != nil {
			return err
		}

		if _, err := db.Exec(
			"INSERT INTO containers (id, image, name, status, command, cpu, memory, cpuset, ip) VALUES (?, ?, ?, ? ,?, ?, ?, ?, ?)",
			c.Id, c.Image, c.Names[0], c.Status, c.Command, info.Config.CpuShares, info.Config.Memory, info.Config.Cpuset, info.NetworkSettings.IpAddress); err != nil {
			return err
		}
	}
	return nil
}
