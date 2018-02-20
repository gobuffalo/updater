package updater

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"html/template"

	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

func WebpackCheck(r *Runner) error {
	fmt.Println("~~~ Checking webpack.config.js ~~~")

	if !r.App.WithWebpack {
		return nil
	}

	path := filepath.Join(envy.GoPath(), "src", "github.com", "gobuffalo", "buffalo", "generators", "assets", "webpack", "templates", "webpack.config.js.tmpl")
	if _, err := os.Stat(path); err != nil {
		return errors.Errorf("could not find webpack.config.js template at %s", path)
	}

	g := newapp.Generator{
		App:       r.App,
		Bootstrap: 3,
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return errors.WithStack(err)
	}

	bb := &bytes.Buffer{}
	err = tmpl.Execute(bb, map[string]interface{}{
		"opts": g,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	b, err := ioutil.ReadFile("webpack.config.js")
	if err != nil {
		return errors.WithStack(err)
	}

	if string(b) == bb.String() {
		return nil
	}

	fmt.Println("Your webpack.config.js file is different from the latest Buffalo template.\nWould you like to replace yours with the latest template? [y/n]")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = strings.ToLower(strings.TrimSpace(text))

	if text == "y" || text == "yes" {
		f, err := os.Create("webpack.config.js")
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = f.Write(bb.Bytes())
		if err != nil {
			return errors.WithStack(err)
		}
		return f.Close()
	}

	return nil
}
