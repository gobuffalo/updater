package updater

import (
	"fmt"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/pkg/errors"
)

type Check func(*Runner) error

type Runner struct {
	App      meta.App
	Warnings []string
}

func Run(checks ...Check) error {
	r := &Runner{
		App:      meta.New("."),
		Warnings: []string{},
	}

	defer func() {
		if len(r.Warnings) == 0 {
			return
		}

		fmt.Println("\n\n----------------------------")
		fmt.Printf("!!! (%d) Warnings Were Found !!!\n\n", len(r.Warnings))
		for _, w := range r.Warnings {
			fmt.Printf("[WARNING]: %s\n", w)
		}
	}()

	for _, c := range checks {
		if err := c(r); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
