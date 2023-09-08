package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// Config is a struct that represents configuration data.
type Config struct {
	Once     sync.Once  // Used for lazy initialization, ensures only one goroutine executes the initialization logic.
	Workers  int64      `json:"numberOfWorkers"` // Number of workers specified in the JSON config file.
	Capacity int64      `json:"queueCapacity"`   // Capacity of the queue specified in the JSON config file.
	FilePath string     `json:"filePath"`
	Lock     sync.Mutex // Used for synchronization while accessing and modifying the config.
}

// LoadConfig loads the configuration from a file specified by filePath and returns a pointer to the Config object.
func LoadConfig(filePath string) *Config {
	config := &Config{} // Allocate memory for the Config struct.

	fmt.Println("initial load of config file")

	config.Once.Do(func() {
		processFile(filePath, config)
	})

	return config // Return the loaded config object.
}

/*This below function can be used to reload the config if something changes dynamically,
	all the changes are dynamically configurable like time to reload from the config, etc.
	ticker := time.NewTicker(5 * time.Second)
    quit := make(chan struct{})

    go func() {
        for {
            select {
            case <-ticker.C:
                config.reloadConfig("path/to/file.json")
            case <-quit:
                ticker.Stop()
                return
            }
        }
    }()
	 <-make(chan struct{})




	 func (config *Config) reloadConfig(filePath string) {
	config.Lock.Lock()
	_ = processFile(filePath, config)
	config.Lock.Unlock()
}
*/

func processFile(filePath string, config *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	line, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		return fmt.Errorf("error occurred reading file: %v", err)
	}

	if err := json.Unmarshal(line, config); err != nil {
		return fmt.Errorf("error occurred unmarshaling JSON: %v", err)
	}

	return nil
}
