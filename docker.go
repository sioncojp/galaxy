package galaxy

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// DockerImageName ...Image name use option "-it" when docker run.
func DockerImageName(image, tag string) string {
	return fmt.Sprintf("%s:%s", image, tag)
}

// ShowAllContainer ...Show running container list
func ShowAllContainer() ([]string, error) {
	var lines []string

	cmd := exec.Command("docker", "ps")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	cmd.Start()

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "galaxy-") {
			if strings.Contains(scanner.Text(), "galaxy-proxy") {
				break
			}
			lines = append(lines, scanner.Text())
		}
	}

	return lines, nil
}

// execCommandString ...Output executing command for ExecDocker function
func execCommandString(cmd, cn string) []string {
	cmdSlice := []string{"exec", "galaxy-" + cn[:7]}
	c := strings.Fields(cmd)
	for _, t := range c {
		cmdSlice = append(cmdSlice, t)
	}
	return cmdSlice
}

// execDocker ...Be able to execute "bin/bash command" to container after container running.
// This logic needed exec.Command().Output().
// Cannot running exec.Command().Run().
func execContainer(cmd, cn string) error {
	_, err := exec.Command("docker",
		execCommandString(cmd, cn)...,
	).Output()
	if err != nil {
		return errors.Wrap(err, "cannot exec container")
	}

	return nil
}

// runContainer ...Run docker
func (config *Config) runContainer(cn string) error {
	err := exec.Command(
		"docker",
		"run",
		"-e",
		"VIRTUAL_HOST="+cn[:7]+"-"+config.Url,
		"-d",
		"--name",
		"galaxy-"+cn[:7],
		"-v",
		CommitDirPath(config.WorkDir, cn)+":/tmp",
		"-it",
		DockerImageName(config.Docker.Image, config.Docker.Tag),
	).Run()
	if err != nil {
		return errors.Wrap(err, "cannot run container")
	}
	if err := execContainer(config.Docker.Exec, cn); err != nil {
		return err
	}

	return nil
}

// createContainer ...Create commit directory->
// Git checkout target commit number->Run container
func (config *Config) createContainer(cn string) error {
	CreateCommitDir(config.WorkDir, cn)
	if err := config.GitCheckoutCommit(cn); err != nil {
		return err
	}
	return config.runContainer(cn)
}

// CreateContainer ...Create container and create a record at commits table
func (config *Config) CreateContainer(cn string) error {
	lock := new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	db, err := config.DBConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	commit := Commits{Number: cn}
	if err := db.Create(&commit).Error; err != nil {
		return err
	}

	if err := config.createContainer(cn); err != nil {
		return err
	}

	return nil
}

// deleteContainer ...Delete container
func deleteContainer(cn string) error {
	err := exec.Command(
		"docker",
		"rm",
		"-f",
		"galaxy-"+cn[:7],
	).Run()
	if err != nil {
		return errors.Wrap(err, "cannot delete container")
	}

	return nil
}

// DeleteContainer ...Detele container and delete a record at commits table
func (config *Config) DeleteContainer(cn string) error {
	db, err := config.DBConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	commit := Commits{Number: cn}
	if err := db.Where(&commit).Find(&Commits{}).Error; err != nil {
		return err
	}
	if err := db.Where(&commit).Delete(&Commits{}).Error; err != nil {
		return err
	}
	deleteContainer(cn)
	DeleteCommitDir(config.WorkDir, cn)

	return nil
}

// DeleteContainerProxy ...Create a proxy container
func (config *Config) CreateContainerProxy() error {
	err := exec.Command(
		"docker",
		"run",
		"-d",
		"-p",
		"80:80",
		"-p",
		"443:443",
		"--name",
		"galaxy-proxy",
		"-v",
		"/var/run/docker.sock:/tmp/docker.sock",
		"-v",
		config.WorkDir+":/tmp",
		"-it",
		DockerImageName(config.Docker.ProxyImage, config.Docker.ProxyTag),
	).Run()
	if err != nil {
		return errors.Wrap(err, "cannot create container proxy")
	}

	return nil
}

// DeleteContainerProxy ...Delete a proxy container
func DeleteContainerProxy() error {
	err := exec.Command(
		"docker",
		"rm",
		"-f",
		"galaxy-proxy",
	).Run()
	if err != nil {
		return errors.Wrap(err, "cannot delete container proxy")
	}

	return nil
}
