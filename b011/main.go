package main

import (
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
}
