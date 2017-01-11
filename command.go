package galaxy

import (
	"fmt"
	"os"
	"os/exec"
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
