package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/transaction_processor/config"
	"github.com/transaction_processor/models"
)

// Transaction represents a single transaction record.

func processTransaction(workerID int, rb *models.RingBuffer, workCh chan models.Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		transaction, ok := <-workCh
		removed := rb.DequeueWithTransaction(transaction)

		if !ok && !removed {
			fmt.Println("nothing to be processed by  worker id " + fmt.Sprintf("%d", workerID) + "\n")
			return // No more transactions to process.
		}

		// Replace this with your actual processing logic.
		fmt.Printf("Processed: %+v by worker ID %d\n", transaction, workerID)
	}
}

// processTransactionsFile reads a file of JSON records, decodes them into Transaction objects,
// and enqueues them into the RingBuffer for processing.
func processTransactionsFile(filename string, rb *models.RingBuffer, workCh chan models.Transaction, wg *sync.WaitGroup) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading line: %w", err)
		}

		var transaction models.Transaction
		if err := json.Unmarshal([]byte(line), &transaction); err != nil {
			return fmt.Errorf("error decoding JSON: %w", err)
		}

		rb.Enqueue(transaction)
		workCh <- transaction // Distribute work to readers.
	}

	close(workCh) // Close the work channel to signal readers to exit when done.

	wg.Wait() // Wait for all processing goroutines to finish.

	return nil
}

func main() {
	config := config.LoadConfig("./config/config.json")
	rb := models.NewRingBuffer(int(config.Capacity))

	// Create a WaitGroup to wait for all processing goroutines to finish.
	var wg sync.WaitGroup

	// Create a channel to distribute work to multiple readers.
	workCh := make(chan models.Transaction, config.Capacity)

	// Start processing goroutines.
	for i := 0; i < int(config.Workers); i++ {
		wg.Add(1)
		go processTransaction(i, rb, workCh, &wg)
	}

	// Open the JSON file, read JSON records, and enqueue them into the RingBuffer.
	err := processTransactionsFile("data.json", rb, workCh, &wg)
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}

}
