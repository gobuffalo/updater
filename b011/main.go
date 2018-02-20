package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gobuffalo/updater"
	"github.com/pkg/errors"
)

var replace = map[string]string{
	"github.com/markbates/pop":      "github.com/gobuffalo/pop",
	"github.com/markbates/validate": "github.com/gobuffalo/validate",
	"github.com/satori/go.uuid":     "github.com/gobuffalo/uuid",
}

var ic = updater.ImportConverter{
	Data: replace,
}

var checks = []updater.Check{
	ic.Process,
	updater.WebpackCheck,
	updater.PackageJSONCheck,
	updater.DepEnsure,
	checkMain,
}

func main() {
	err := updater.Run(checks...)
	if err != nil {
		log.Fatal(err)
	}
}

func checkMain(*updater.Runner) error {
	fmt.Println("~~~ Checking main.go ~~~")
	b, err := ioutil.ReadFile("main.go")
	if err != nil {
		return errors.WithStack(err)
	}
	if bytes.Contains(b, []byte("app.Start")) {
		fmt.Println("[Warning]: app.Start has been removed in v0.11.0. Use app.Serve Instead.")
	}
	return nil
}
