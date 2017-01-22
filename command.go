package galaxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// RunScript ...Run script file. be able to set "commit number path" in script argument.
func RunScript(script, cnPath string) error {
	if len(script) == 0 {
		return nil
	}

	if err := exec.Command("/bin/bash", script, cnPath).Run(); err != nil {
		return err
	}
	return nil
}

// CommitDirPath ...Output"workdir/commit number".
func CommitDirPath(workdir, cn string) string {
	return fmt.Sprintf("%s/%s", workdir, cn[:7])
}

// repoWorkDirString ...Output repository_work_dir.
func (config *Config) repoWorkDirString() string {
	workdir := config.WorkDir
	reponame := config.Github.Name

	return fmt.Sprintf("%s/%s", workdir, reponame)
}

// ChangeRepositoryDir ...Change repository dir.
func (config *Config) ChangeRepositoryDir() (string, error) {
	prev, err := filepath.Abs(".")
	if err != nil {
		return "", errors.Wrap(err, "cannot change repository dir")
	}
	os.Chdir(config.repoWorkDirString())

	return prev, nil
}

// RevertDir ...Revert before changed repository dir.
func RevertDir(prev string) {
	os.Chdir(prev)
}

// CreateCommitDir ...Create "workdir/commit number" directory for mounting docker conainer.
func CreateCommitDir(workdir, cn string) {
	path := CommitDirPath(workdir, cn)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
}

// DeleteCommitDir ...Create "workdir/commit number" directory for mounting docker conainer.
func DeleteCommitDir(workdir, cn string) {
	path := CommitDirPath(workdir, cn[:7])
	os.RemoveAll(path)
}
