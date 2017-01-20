package galaxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// GitClone ...Git clone target repository to repository_work_dir.
func (config *Config) GitClone() error {
	err := exec.Command(
		"git",
		"clone",
		config.Github.Repository,
		config.WorkDir+"/"+config.Github.Name,
	).Run()
	if err != nil {
		return errors.Wrap(err, "Failed git clone")
	}
	return nil
}

// gitDirString ...Output "--git-dir=repository_work_dir".
func (config *Config) gitDirString() string {
	workdir := config.WorkDir
	gitname := config.Github.Name

	return fmt.Sprintf("--git-dir=%s/%s/.git", workdir, gitname)
}

// repoWorkDirString ...Output repository_work_dir.
func (config *Config) repoWorkDirString() string {
	workdir := config.WorkDir
	reponame := config.Github.Name

	return fmt.Sprintf("%s/%s", workdir, reponame)
}

// GitCheckoutMaster ...Git checkout master.
func (config *Config) gitCheckoutMaster(cn string) error {
	err := exec.Command(
		"git",
		"checkout",
		"-f",
		"master",
	).Run()
	if err != nil {
		return errors.Wrap(err, "git checkout master")
	}
	return config.gitPull(cn)
}

// GitPull ...Git pull origin master.
func (config *Config) gitPull(cn string) error {
	err := exec.Command(
		"git",
		"pull",
		"origin",
		"master",
	).Run()
	if err != nil {
		return errors.Wrap(err, "git pull")
	}
	return config.gitCheckoutCommitNumber(cn)
}

// gitCheckoutCommit ...Git checkout commit-number.
func (config *Config) gitCheckoutCommitNumber(cn string) error {
	err := exec.Command(
		"git",
		"checkout",
		cn,
	).Run()
	if err != nil {
		return errors.Wrap(err, "git checkout commit number")
	}
	return RunScript(config.Script, config.WorkDir+"/"+cn[:7])
}

func (config *Config) gitCheckoutCommit(cn string) error {
	return config.gitCheckoutMaster(cn)
}

// GitCheckoutCommit ...Do as below to target repository.
// Chdir to repository directory -> git checkout master
// -> git pull origin master -> git checkout commit_number
// -> run script file
func (config *Config) GitCheckoutCommit(cn string) error {
	prev, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	defer os.Chdir(prev)
	os.Chdir(config.repoWorkDirString())

	if err := config.gitCheckoutCommit(cn); err != nil {
		return err
	}

	return nil
}

// checkCommitNumber ...Check commit number. set minimum size is 7-digit.
func checkCommitNumber(cn string) error {
	if len(cn) >= 7 || cn == "master" {
		return nil
	}
	return fmt.Errorf("commit number is very short")
}

// GitCommitNumerTo40digit ...Output commit number to 40-digit from 7-digit.
func (config *Config) GitCommitNumerTo40digit(cn string) (string, error) {
	if err := checkCommitNumber(cn); err != nil {
		return "", err
	}

	out, err := exec.Command(
		"git",
		config.gitDirString(),
		"show",
		cn,
		"-s",
		"--format=%H",
	).Output()

	if err != nil {
		return "", errors.Wrap(err, "unknown revision or path not in the working tree")
	}
	return strings.TrimRight(string(out), "\n"), nil
}
