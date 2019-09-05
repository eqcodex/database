package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

var (
	syntax = `usage: preview database value
supported databases: item, mob, zone

for example, 'preview mob a bixie'`
	database = ""
)

// Mob represents a mob
type Mob struct {
	ID         string
	Name       string
	DropGroups []DropGroup
}

// DropGroup represents a drop database
type DropGroup struct {
	Name   string
	Chance int
	Items  []Item
}

// Item represents an item
type Item struct {
	ID         string
	Chance     int
	Name       string
	Allakhazam int
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args
	if len(args) < 3 {
		return fmt.Errorf(syntax)
	}
	database = args[1]
	value := strings.Join(args[2:], " ")
	fmt.Println("searching", database, "for", value)

	err := filepath.Walk(database, findPattern)
	if err != nil {
		return errors.Wrapf(err, "failed to find path %s", database)
	}

	return nil
}

func findPattern(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	if filepath.Ext(info.Name()) != ".yaml" {
		return nil
	}
	switch database {
	case "mob":
		m := &Mob{}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read file %s", path)
		}

		err = yaml.Unmarshal(data, m)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal file %s", path)
		}
		fmt.Println(m)
	default:
		return fmt.Errorf("unknown database: %s", database)
	}

	fmt.Printf("visited file or dir: %q\n", path)
	return nil
}
