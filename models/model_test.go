package models

import "testing"

func TestRingBufferEnqueue(t *testing.T) {
	rb := NewRingBuffer(3)

	transactions := []Transaction{
		{"2023-06-27 22:22:19.62710.192501", "1", "p5"},
		{"2023-06-27 22:22:19.62710.193409", "2", "p4"},
		{"2023-06-27 22:22:19.62710.193441", "3", "p4"},
	}

	for _, transaction := range transactions {
		rb.Enqueue(transaction)
	}

	expectedTransactions := []Transaction{
		{"2023-06-27 22:22:19.62710.192501", "1", "p5"},
		{"2023-06-27 22:22:19.62710.193409", "2", "p4"},
		{"2023-06-27 22:22:19.62710.193441", "3", "p4"},
	}

	for i, expectedTransaction := range expectedTransactions {
		transaction := rb.Data[(rb.Head+i)%rb.Size]
		if transaction != expectedTransaction {
			t.Errorf("Expected: %+v, Got: %+v", expectedTransaction, transaction)
		}
	}

	rb.Enqueue(Transaction{"2023-06-27 22:22:19.62710.193470", "4", "p1"})

	expectedTransactions = []Transaction{
		{"2023-06-27 22:22:19.62710.193409", "2", "p4"},
		{"2023-06-27 22:22:19.62710.193441", "3", "p4"},
		{"2023-06-27 22:22:19.62710.193470", "4", "p1"},
	}

	for i, expectedTransaction := range expectedTransactions {
		transaction := rb.Data[(rb.Head+i)%rb.Size]
		if transaction != expectedTransaction {
			t.Errorf("Expected: %+v, Got: %+v", expectedTransaction, transaction)
		}
	}
}

func TestRingBufferDequeue(t *testing.T) {
	rb := NewRingBuffer(3)

	transactions := []Transaction{
		{"2023-06-27 22:22:19.62710.192501", "1", "p5"},
		{"2023-06-27 22:22:19.62710.193409", "2", "p4"},
		{"2023-06-27 22:22:19.62710.193441", "3", "p4"},
	}
	for _, transaction := range transactions {
		rb.Enqueue(transaction)
	}

	for _, expectedTransaction := range transactions {
		transaction := rb.Dequeue()
		if transaction != expectedTransaction {
			t.Errorf("Expected: %+v, Got: %+v", expectedTransaction, transaction)
		}
	}

	emptyTransaction := rb.Dequeue()
	if emptyTransaction != (Transaction{}) {
		t.Errorf("Expected an empty buffer, Got: %+v", emptyTransaction)
	}
}

func TestDequeueWithTransaction(t *testing.T) {
	rb := NewRingBuffer(3)

	// Enqueue some transactions
	t1 := Transaction{Datetime: "2022-01-01", Value: "100", Partition: "A"}
	t2 := Transaction{Datetime: "2022-01-02", Value: "200", Partition: "B"}
	t3 := Transaction{Datetime: "2022-01-03", Value: "300", Partition: "C"}
	rb.Enqueue(t1)
	rb.Enqueue(t2)
	rb.Enqueue(t3)

	// Dequeue a transaction that exists in the buffer
	found := rb.DequeueWithTransaction(t2)
	if !found {
		t.Errorf("Expected DequeueWithTransaction to return true, but got false")
	}
	if rb.Count != 2 {
		t.Errorf("Expected Count to be 2, but got %d", rb.Count)
	}
	if rb.Data[0] != t1 {
		t.Errorf("Expected Data[0] to be %v, but got %v", t1, rb.Data[0])
	}
	if rb.Data[1] != t3 {
		t.Errorf("Expected Data[1] to be %v, but got %v", t3, rb.Data[1])
	}

	// Dequeue a transaction that doesn't exist in the buffer
	found = rb.DequeueWithTransaction(Transaction{Datetime: "2022-01-04", Value: "400", Partition: "D"})
	if found {
		t.Errorf("Expected DequeueWithTransaction to return false, but got true")
	}
	if rb.Count != 2 {
		t.Errorf("Expected Count to still be 2, but got %d", rb.Count)
	}
	if rb.Data[0] != t1 {
		t.Errorf("Expected Data[0] to remain %v, but got %v", t1, rb.Data[0])
	}
	if rb.Data[1] != t3 {
		t.Errorf("Expected Data[1] to remain %v, but got %v", t3, rb.Data[1])
	}
}
