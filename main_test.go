package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("should report parsing errors", func(t *testing.T) {
		t.Run("should report that unexpected flags", func(t *testing.T) {
			// given
			args := []string{"-x"}

			// when
			config, err := readConfig(args)

			// then
			if config != nil {
				t.Fatal("config should have been nil")
			}
			if err == nil {
				t.Fatal("error should not have been nil")
			}
		})
	})

	t.Run("should skip errors due to the help being called", func(t *testing.T) {
		// given
		args := []string{"-h"}

		// when
		_, err := readConfig(args)

		// then
		if err != nil {
			t.Fatal("error should have been nil")
		}
	})

	t.Run("should use default values when the flags are not specified", func(t *testing.T) {
		// given
		var args []string

		// when
		config, err := readConfig(args)

		// then
		if config == nil { t.Fatal("config should not have been nil") }
		if err != nil { t.Fatal("config should not have been nil") }
		assertEmptyString(config.ProjectName, "project name", t)
		assertEmptyString(config.ProjectLocation, "project location", t)
		assertEmptyString(config.ProjectRepository, "project repository", t)
		if config.WebAppProject != false { t.Fatal("project type should have been false") }
	})
}

func TestCheckParameters(t *testing.T)  {
	t.Run("should return an error message if the parameter is empty", func(t *testing.T) {
		// given
		name := "testParamName"
		value := ""

		// when
		errorMessage := checkParameter(name, value)

		// then
		if len(errorMessage) == 0 { t.Fatal("error message should not have been empty") }
		if errorMessage != fmt.Sprintf("Project %s cannot be empty", name) { t.Fatalf("error message is '%s' instead of being the expected value", errorMessage) }
	})

	t.Run("should return an empty string if the parameter is not empty", func(t *testing.T) {
		// given
		name := "testParamName"
		value := "testValue"

		// when
		errorMessage := checkParameter(name, value)

		// then
		if len(errorMessage) > 0 { t.Fatal("error message should have been empty") }
	})
}

func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		config Config
		numOfErrors int
	}{
		{ config: Config{}, numOfErrors: 3 },
		{ config: Config{ProjectName: "name"}, numOfErrors: 2 },
		{ config: Config{ProjectName: "name", ProjectLocation: "path"}, numOfErrors: 1 },
		{ config: Config{ProjectName: "name", ProjectLocation: "path", ProjectRepository: "repository"}, numOfErrors: 0 },
		{ config: Config{ProjectName: "name", ProjectLocation: "path", ProjectRepository: "http://repository"}, numOfErrors: 0 },
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("should report %d errors", testCase.numOfErrors), func(t *testing.T) {
			// given
			config := testCase.config

			// when
			errors := config.Validate()

			// then
			if len(errors) != testCase.numOfErrors {
				t.Fatalf("Expected %d errors, but got %d instead", testCase.numOfErrors, len(errors))
			}
		})
	}
}

func TestMain(m *testing.M) {
	// given
	output := strings.Builder{}
	log.SetOutput(&output)
	const name = "test-name"
	const path = "./test-path"
	var expectedOutput = fmt.Sprintf("Generating scaffold for project %s in %s", name, path)
	os.Args = []string{"test", "-n", name, "-d", path, "-r", "some-repo"}

	// when
	exitCode := m.Run()

	// then
	if exitCode != 0 { panic("exit code should have been 0") }
	if output.String() != expectedOutput { panic("output is not the expected one") }
}

// using testify would be way better
func assertEmptyString(value string, name string, t *testing.T) {
	if len(value) > 0 { t.Fatalf("%s should have been empty", name) }
}
