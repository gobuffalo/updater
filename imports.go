package updater

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ImportConverter struct {
	Data map[string]string
}

func (c ImportConverter) Process() error {
	fmt.Println("~~~ Rewriting Imports ~~~")
	err := filepath.Walk(".", func(p string, info os.FileInfo, err error) error {
		for _, n := range []string{"vendor", "node_modules", ".git"} {
			if strings.HasPrefix(p, n+string(filepath.Separator)) {
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(p)
		if ext == ".go" {
			if err := c.rewriteFile(p); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := os.Stat("Gopkg.toml"); err == nil {
		b, err := ioutil.ReadFile("Gopkg.toml")
		if err != nil {
			return errors.WithStack(err)
		}
		warn := []string{}
		for k := range c.Data {
			if bytes.Contains(b, []byte(k)) {
				warn = append(warn, k)
			}
		}
		if len(warn) > 0 {
			fmt.Println("[WARNING] Your Gopkg.toml contains the following imports that need to be changed MANUALLY:")
			for _, n := range warn {
				fmt.Printf("\t%s -> %s\n", n, c.Data[n])
			}
		}
	}
	return nil
}

// TAKEN FROM https://gist.github.com/jackspirou/61ce33574e9f411b8b4a
// rewriteFile rewrites import statments in the named file
// according to the rules supplied by the map of strings.
func (c ImportConverter) rewriteFile(name string) error {

	// create an empty fileset.
	fset := token.NewFileSet()

	// parse the .go file.
	// we are parsing the entire file with comments, so we don't lose anything
	// if we need to write it back out.
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		e := err.Error()
		msg := "expected 'package', found 'EOF'"
		if e[len(e)-len(msg):] == msg {
			return nil
		}
		return err
	}

	// iterate through the import paths. if a change occurs update bool.
	change := false
	for _, i := range f.Imports {

		// unquote the import path value.
		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			return err
		}

		// match import path with the given replacement map
		if path, ok := c.match(path); ok {
			i.Path.Value = strconv.Quote(path)
			change = true
		}
	}

	for _, cg := range f.Comments {
		for _, cl := range cg.List {
			if strings.HasPrefix(cl.Text, "// import \"") {

				// trim off extra comment stuff
				ctext := cl.Text
				ctext = strings.TrimPrefix(ctext, "// import")
				ctext = strings.TrimSpace(ctext)

				// unquote the comment import path value
				ctext, err := strconv.Unquote(ctext)
				if err != nil {
					return err
				}

				// match the comment import path with the given replacement map
				if ctext, ok := c.match(ctext); ok {
					cl.Text = "// import " + strconv.Quote(ctext)
					change = true
				}
			}
		}
	}

	// if no change occured, then we don't need to write to disk, just return.
	if !change {
		return nil
	}

	// since the imports changed, resort them.
	ast.SortImports(fset, f)

	// create a temporary file, this easily avoids conflicts.
	temp := name + ".temp"
	w, err := os.Create(temp)
	if err != nil {
		return err
	}

	// write changes to .temp file, and include proper formatting.
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(w, fset, f)
	if err != nil {
		return err
	}

	// close the writer
	err = w.Close()
	if err != nil {
		return err
	}

	// rename the .temp to .go
	return os.Rename(temp, name)
}

// match takes an import path and replacement map.
func (c ImportConverter) match(importpath string) (string, bool) {
	for key, value := range c.Data {
		if len(importpath) >= len(key) {
			if importpath[:len(key)] == key {
				result := path.Join(value, importpath[len(key):])
				return result, true
			}
		}
	}
	return importpath, false
}
