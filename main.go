package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"scaffold_gen/fileutils"
	"strings"
	"text/template"
)

type Config struct {
	ProjectName       string
	ProjectLocation   string
	ProjectRepository string
	WebAppProject     bool
}

func checkParameter(name string, value string) string {
	if len(value) == 0 {
		return fmt.Sprintf("Project %s cannot be empty", name)
	}

	return ""
}

func validatePath(path string) []string {
	var errors []string
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("The path '%s' doesn't exist", path))
		} else if os.IsPermission(err) {
			errors = append(errors, fmt.Sprintf("The path '%s' doesn't exist", path))
		} else {
			errors = append(errors, err.Error())
		}
	} else if !fileInfo.IsDir() {
		errors = append(errors, fmt.Sprintf("The path '%s' is not a directory", path))
	} else if !fileutils.IsDirectoryWritable(path) {
		errors = append(errors, fmt.Sprintf("The path '%s' is not writeable", path))
	}

	return errors
}

func validateRepositoryUrl(repositoryUrl string) []string {
	var errors []string
	parsedUrl, err := url.Parse(repositoryUrl)
	if err != nil {
		errors = append(errors, err.Error())
	} else if len(parsedUrl.Scheme) > 0 {
		errors = append(errors, fmt.Sprintf("The repository URL must not specify the scheme, but '%s' was found", parsedUrl.Scheme))
	}

	return errors
}

func (c *Config) Validate() []string {
	var errors []string
	var parameters = map[string]string{
		"name": c.ProjectName,
		"path": c.ProjectLocation,
		"repository URL": c.ProjectRepository,
	}
	for name, value := range parameters {
		if err := checkParameter(name, value); len(err) > 0 {
			errors = append(errors, err)
		}
	}
	if len(errors) == 0 {
		errors = append(errors, validateRepositoryUrl(c.ProjectRepository)...)
		errors = append(errors, validatePath(c.ProjectLocation)...)
	}

	return errors
}

func readConfig(args []string) (*Config, error) {
	var flagSet = flag.NewFlagSet("parameters", flag.ContinueOnError)
	var config = Config{}
	flagSet.StringVar(&config.ProjectName, "n", "", "Project name")
	flagSet.StringVar(&config.ProjectLocation, "d", "", "Project location on disk")
	flagSet.StringVar(&config.ProjectRepository, "r", "", "Project remote repository URL")
	flagSet.BoolVar(&config.WebAppProject, "s", false, "Project will have static assets or not")

	if err := flagSet.Parse(args); err != nil && err != flag.ErrHelp {
		return nil, err
	}

	return &config, nil
}

func resolveTemplate(templateFilename string, config *Config) {
	parsedTemplate := template.Must(template.ParseFiles(templateFilename))
	defer os.Remove(templateFilename)

	filename := templateFilename[0: strings.Index(templateFilename, filepath.Ext(templateFilename))]
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("An error occurred while creating the file %s from the template %s:\n%v", filename, templateFilename, err)
	}
	defer file.Close()

	if err := parsedTemplate.Execute(file, config); err != nil {
		log.Fatalf("An error occurred while resolving the template %s:\n%v", templateFilename, err)
	}
}

func main() {
	config, err := readConfig(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	if errors := config.Validate(); len(errors) > 0 {
		log.Fatalln(strings.Join(errors, "\n"))
	}

	log.Printf("Generating scaffold for project %s in %s", config.ProjectName, config.ProjectLocation)

	var projectName string
	if config.WebAppProject {
		projectName = "webapp"
	} else {
		projectName = "api"
	}
	templateFiles := fileutils.CopyTree(filepath.Join(".", "templates", projectName), filepath.Join(config.ProjectLocation, config.ProjectName))
	for _, templateFile := range templateFiles {
		resolveTemplate(templateFile, config)
	}
}
