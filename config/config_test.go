package config

import (
	"os"
	"testing"
)

func TestProcessFile(t *testing.T) {
	// Create a temporary file for testing.
	tempFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some test data to the temporary file.
	testData := `{
        "numberOfWorkers": 5,
        "queueCapacity": 50,
        "filePath": "/path/to/data.txt"
    }`
	err = os.WriteFile(tempFile.Name(), []byte(testData), 0644)
	if err != nil {
		t.Fatalf("error writing test data to file: %v", err)
	}

	// Create a new Config object for testing.
	config := &Config{}

	// Call the ProcessFile function with the temporary file and config object.
	err = processFile(tempFile.Name(), config)
	if err != nil {
		t.Fatalf("error processing file: %v", err)
	}

	// Verify that the config object has been updated correctly.
	if config.Workers != 5 {
		t.Errorf("expected number of workers to be 5, got %d", config.Workers)
	}
	if config.Capacity != 50 {
		t.Errorf("expected queue capacity to be 50, got %d", config.Capacity)
	}
	if config.FilePath != "/path/to/data.txt" {
		t.Errorf("expected file path to be '/path/to/data.txt', got '%s'", config.FilePath)
	}
}

// TestLoadConfig tests the LoadConfig function.
func TestLoadConfig(t *testing.T) {
	// Create a temporary file for testing.
	tmpfile, err := os.CreateTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Write test JSON data to the temporary file.
	jsonData := `{
        "numberOfWorkers": 5,
        "queueCapacity": 50,
        "filePath": "/path/to/data.txt"
    }`
	if _, err := tmpfile.Write([]byte(jsonData)); err != nil {
		t.Fatalf("Failed to write test data to file: %v", err)
	}

	// Call the LoadConfig function with the path to the temporary file.
	config := LoadConfig(tmpfile.Name())

	// Verify that the config object is correctly loaded.
	expectedWorkers := 5
	expectedCapacity := 50
	expectedFilePath := "/path/to/data.txt"
	if int(config.Workers) != expectedWorkers {
		t.Errorf("Unexpected number of workers. Expected: %d, Actual: %d", expectedWorkers, config.Workers)
	}
	if int(config.Capacity) != expectedCapacity {
		t.Errorf("Unexpected queue capacity. Expected: %d, Actual: %d", expectedCapacity, config.Capacity)
	}
	if config.FilePath != expectedFilePath {
		t.Errorf("Unexpected FilePath. Expected: %v, Actual: %v", expectedFilePath, config.FilePath)
	}
}
