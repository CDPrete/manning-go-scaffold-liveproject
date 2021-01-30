package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	ProjectName       string
	ProjectLocation   string
	ProjectRepository string
	ProjectType       bool
}

func checkParameter(name string, value string) string {
	if len(value) == 0 {
		return fmt.Sprintf("Project %s cannot be empty", name)
	}

	return ""
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

	return errors
}

func readConfig(args []string) (*Config, error) {
	var flagSet = flag.NewFlagSet("parameters", flag.ContinueOnError)
	var config = Config{}
	flagSet.StringVar(&config.ProjectName, "n", "", "Project name")
	flagSet.StringVar(&config.ProjectLocation, "d", "", "Project location on disk")
	flagSet.StringVar(&config.ProjectRepository, "r", "", "Project remote repository URL")
	flagSet.BoolVar(&config.ProjectType, "s", false, "Project will have static assets or not")

	if err := flagSet.Parse(args); err != nil && err != flag.ErrHelp {
		return nil, err
	}

	return &config, nil
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
}
