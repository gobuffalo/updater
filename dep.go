package updater

import (
	"fmt"
	"os"
	"os/exec"
)

func DepEnsure(r *Runner) error {
	if !r.App.WithDep {
		return nil
	}
	fmt.Println("~~~ Running dep ensure ~~~")
	cc := exec.Command("dep", "ensure", "-v")
	cc.Stdin = os.Stdin
	cc.Stderr = os.Stderr
	cc.Stdout = os.Stdout
	return cc.Run()
}
