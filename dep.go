package updater

import (
	"fmt"
	"os"
	"os/exec"
)

func DepEnsure() error {
	if _, err := os.Stat("Gopkg.toml"); err == nil {
		fmt.Println("~~~ Running dep ensure ~~~")
		cc := exec.Command("dep", "ensure")
		cc.Stdin = os.Stdin
		cc.Stderr = os.Stderr
		cc.Stdout = os.Stdout
		return cc.Run()
	}
	return nil
}
