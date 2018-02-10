package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gobuffalo/updater"
)

var replace = map[string]string{
	"github.com/markbates/pop":      "github.com/gobuffalo/pop",
	"github.com/markbates/validate": "github.com/gobuffalo/validate",
	"github.com/satori/go.uuid":     "github.com/gobuffalo/uuid",
}

func main() {
	ic := updater.ImportConverter{
		Data: replace,
	}
	if err := ic.Process(); err != nil {
		log.Fatal(err)
	}
	if err := updater.DepEnsure(); err != nil {
		log.Fatal(err)
	}
	checkMain()
}

func checkMain() {
	fmt.Println("~~~ Checking main.go ~~~")
	b, err := ioutil.ReadFile("main.go")
	if err != nil {
		log.Fatal(err)
	}
	if bytes.Contains(b, []byte("app.Start")) {
		fmt.Println("[Warning]: app.Start has been removed in v0.11.0. Use app.Serve Instead.")
	}
}
