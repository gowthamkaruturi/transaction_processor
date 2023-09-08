package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/transaction_processor/models"
)

func Test_processTransactionsFile(t *testing.T) {
	type args struct {
		filename string
		rb       *models.RingBuffer
		workCh   chan models.Transaction
		wg       *sync.WaitGroup
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid File",
			args: args{
				filename: "./test_data/test_validfile.json",
				rb:       models.NewRingBuffer(3),
				workCh:   make(chan models.Transaction, 3),
				wg:       &sync.WaitGroup{},
			},
			wantErr: false,
		},
		{
			name: "Invalid File",
			args: args{
				filename: "./test_data/nonexistent_file.json",
				rb:       models.NewRingBuffer(3),
				workCh:   make(chan models.Transaction, 3),
				wg:       &sync.WaitGroup{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			go func() {
				transactions := make([]models.Transaction, 0)
				for transaction := range tt.args.workCh {
					// Process the transaction here
					transactions = append(transactions, transaction)
					fmt.Println("Processing transaction:", transaction)
				}
				if len(transactions) == cap(tt.args.workCh) {
					close(tt.args.workCh)
				}
				// Send the mock transaction to the channel
				// Close the channel after sending the value(s)
			}()

			err := processTransactionsFile(tt.args.filename, tt.args.rb, tt.args.workCh, tt.args.wg)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTransactionsFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
