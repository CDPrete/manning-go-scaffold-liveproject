package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		{ config: Config{ProjectName: "name", ProjectLocation: "."}, numOfErrors: 1 },
		{ config: Config{ProjectName: "name", ProjectLocation: ".", ProjectRepository: "http://localhost"}, numOfErrors: 1 },
		{ config: Config{ProjectName: "name", ProjectLocation: ".", ProjectRepository: "localhost"}, numOfErrors: 0 },
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

func TestMainRun(t *testing.T) {
	// given
	output := strings.Builder{}
	log.SetOutput(&output)
	const name = "test-name"
	var path = t.TempDir()
	var expectedOutput = fmt.Sprintf("Generating scaffold for project %s in %s", name, path)
	os.Args = []string{"test", "-n", name, "-d", path, "-r", "test-repo"}

	// when
	main()

	// then
	if !strings.Contains(output.String(), expectedOutput) {
		t.Fatalf("output is not the expected one\n found: %s, expected %s", output.String(), expectedOutput)
	}
}

func TestMain_Generation(t *testing.T) {
	const name = "test-name"
	const repository = "main"
	for _, projectType := range []bool{true, false} {
		t.Run(fmt.Sprintf("Generating %s project", getProjectTypeName(projectType)), func(t *testing.T) {
			// given
			var path = t.TempDir()
			var rootGeneratedPath = filepath.Join(path, name)
			os.Args = []string{"test", "-n", name, "-d", path, "-r", repository}
			if projectType {
				os.Args = append(os.Args, "-s")
			}
			exampleFiles := make(map[string]*string)
			var rootExamplesPath = filepath.Join(".", "examples", getProjectTypeName(projectType))
			walkFunc := func(path string, info os.FileInfo, err error) error {
				pathWithoutRoot := strings.TrimPrefix(path, rootExamplesPath)
				if info.IsDir() {
					exampleFiles[pathWithoutRoot] = nil
				} else {
					var data []byte
					if data, err = ioutil.ReadFile(path); err == nil {
						dataStr := string(data)
						exampleFiles[pathWithoutRoot] = &dataStr
					}
				}
				return err
			}
			if err := filepath.Walk(rootExamplesPath, walkFunc); err != nil {
				t.Fatal(err)
			}

			// when
			main()

			// then
			assertWalkFunc := func(path string, info os.FileInfo, err error) error {
				pathWithoutRoot := strings.TrimPrefix(path, rootGeneratedPath)
				if value, exists := exampleFiles[pathWithoutRoot]; !exists {
					err = fmt.Errorf("%s not found in the expected file tree", path)
				} else if value != nil {
					var data []byte
					if data, err = ioutil.ReadFile(path); err == nil && !strings.EqualFold(string(data), *value) {
						err = fmt.Errorf("\ncontent expected:\n%s\ncontent found:\n%s", *value, string(data))
					}
				}
				return err
			}
			if err := filepath.Walk(rootGeneratedPath, assertWalkFunc); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// using testify would be way better
func assertEmptyString(value string, name string, t *testing.T) {
	if len(value) > 0 { t.Fatalf("%s should have been empty", name) }
}
