package main

import (
	"flag"
	"log"
	"os"
	"path"
)

type item struct {
	items       []item
	name        string
	isDirectory bool
}

type projectStructure struct {
	items []item
}

func main() {
	projectName := flag.String("project-name", "", "Project name")
	flag.Parse()
	if len(*projectName) == 0 {
		log.Fatal("project-name is required")
	}
	project := projectStructure{
		items: []item{
			{name: "api", isDirectory: true},
			{name: "build", isDirectory: true},
			{name: "cmd", isDirectory: true, items: []item{
				{name: "main.go", isDirectory: false},
			},
			},
			{name: "deploy", isDirectory: true},
			{name: "docs", isDirectory: true},
			{name: "internal", isDirectory: true, items: []item{
				{name: *projectName, isDirectory: true, items: []item{
					{name: "rest", isDirectory: true},
					{name: "service", isDirectory: true},
				},
				},
			},
			},
			{name: "pkg", isDirectory: true},
		},
	}

	visited := make(map[string]struct{})
	curPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	createItem(item{name: *projectName, isDirectory: true}, []string{curPath})
	for _, item := range project.items {
		if _, ok := visited[item.name]; !ok {
			traverse(item, []string{curPath, *projectName}, visited)
		}
	}
}

func traverse(item item, paths []string, visited map[string]struct{}) {
	visited[item.name] = struct{}{}
	createItem(item, paths)
	for _, child := range item.items {
		if _, ok := visited[child.name]; !ok {
			traverse(child, append(paths, item.name), visited)
		}
	}
}

func createItem(item item, paths []string) {
	currentPath := path.Join(append(paths, item.name)...)
	if item.isDirectory {
		if err := os.Mkdir(currentPath, os.FileMode(0755)); err != nil {
			log.Fatal(err)
		}
		return
	}
	if _, err := os.Create(currentPath); err != nil {
		log.Fatal(err)
	}
}
