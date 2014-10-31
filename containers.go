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
		"CREATE TABLE containers (id, image, name, state, command, cpu, memory, cpuset, ip)"); err != nil {
		return err
	}
	for _, c := range containers {
		info, err := client.InspectContainer(c.Id)
		if err != nil {
			return err
		}

		if _, err := db.Exec(
			`INSERT INTO containers 
            (id, image, name, state, command, cpu, memory, cpuset, ip) 
            VALUES (?, ?, ?, ? ,?, ?, ?, ?, ?)`,
			c.Id, c.Image, info.Name, getState(info), c.Command, info.Config.CpuShares, info.Config.Memory,
			info.Config.Cpuset, info.NetworkSettings.IpAddress); err != nil {
			return err
		}
	}
	return nil
}

func getState(info *dockerclient.ContainerInfo) string {
	s := info.State
	switch {
	case s.Paused:
		return "paused"
	case s.Running:
		return "running"
	case s.Restarting:
		return "restarting"
	default:
		return "stopped"
	}
}
