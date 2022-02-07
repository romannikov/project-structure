package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"
)

type item struct {
	Items       []item `json:"Items"`
	Name        string `json:"Name"`
	IsDirectory bool   `json:"is_directory"`
}

type projectStructure struct {
	Items []item `json:"Items"`
}

func main() {
	projectName := flag.String("name", "", "A project name")
	structureFile := flag.String("structure", "default.json", "A file with a project structure description")
	flag.Parse()
	if len(*projectName) == 0 {
		log.Fatal("Name is required.")
	}
	project := getProjectStructure(*structureFile)

	visited := make(map[string]struct{})
	patterns := make(map[string]string)
	patterns["{project_name}"] = *projectName
	curPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	createItem(item{Name: *projectName, IsDirectory: true}, []string{curPath}, nil)
	for _, item := range project.Items {
		if _, ok := visited[item.Name]; !ok {
			traverse(item, []string{curPath, *projectName}, visited, patterns)
		}
	}
}

func getProjectStructure(filename string) projectStructure {
	byteContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var project projectStructure
	err = json.Unmarshal(byteContent, &project)
	if err != nil {
		log.Fatal(err)
	}
	return project
}

func traverse(item item, paths []string, visited map[string]struct{}, patterns map[string]string) {
	visited[item.Name] = struct{}{}
	createItem(item, paths, patterns)
	for _, child := range item.Items {
		if _, ok := visited[child.Name]; !ok {
			traverse(child, append(paths, patternValueOrPureValue(patterns, item.Name)), visited, patterns)
		}
	}
}

func createItem(item item, paths []string, patterns map[string]string) {
	currentPath := path.Join(append(paths, patternValueOrPureValue(patterns, item.Name))...)
	if item.IsDirectory {
		if err := os.Mkdir(currentPath, os.FileMode(0755)); err != nil {
			log.Fatal(err)
		}
		return
	}
	if _, err := os.Create(currentPath); err != nil {
		log.Fatal(err)
	}
}

func patternValueOrPureValue(patterns map[string]string, name string) string {
	if value, ok := patterns[name]; ok {
		return value
	}
	return name
}
