package models

import "sync"

type Transaction struct {
	Datetime  string `json:"datetime"`
	Value     string `json:"value"`
	Partition string `json:"partition"`
}

type RingBuffer struct {
	Data  []Transaction
	Size  int
	Head  int
	Tail  int
	Count int
	Lock  sync.Mutex
}

// NewRingBuffer creates a new RingBuffer with the specified capacity.
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		Data: make([]Transaction, capacity),
		Size: capacity,
	}
}

// Enqueue adds a transaction to the RingBuffer.
func (rb *RingBuffer) Enqueue(transaction Transaction) {
	rb.Lock.Lock()
	defer rb.Lock.Unlock()

	if rb.Count == rb.Size {
		// Buffer is full, overwrite the oldest data
		rb.Head = (rb.Head + 1) % rb.Size
	}

	rb.Data[rb.Tail] = transaction
	rb.Tail = (rb.Tail + 1) % rb.Size
	rb.Count++
}

// Dequeue removes and returns the oldest transaction from the RingBuffer.
func (rb *RingBuffer) Dequeue() Transaction {
	rb.Lock.Lock()
	defer rb.Lock.Unlock()

	if rb.Count == 0 {
		return Transaction{}
	}

	transaction := rb.Data[rb.Head]
	rb.Head = (rb.Head + 1) % rb.Size
	rb.Count--
	return transaction
}

func (rb *RingBuffer) DequeueWithTransaction(transaction Transaction) bool {
	rb.Lock.Lock()
	defer rb.Lock.Unlock()

	if rb.Count == 0 {
		return false
	}

	// Find the index of the specified transaction in the buffer
	index := -1
	for i := 0; i < rb.Count; i++ {
		currentIndex := (rb.Head + i) % rb.Size
		if rb.Data[currentIndex] == transaction {
			index = currentIndex
			break
		}
	}

	if index == -1 {
		return false
	}

	// Shift the subsequent transactions to fill the gap
	for i := index; i < rb.Count-1; i++ {
		current := (rb.Head + i) % rb.Size
		next := (current + 1) % rb.Size
		rb.Data[current] = rb.Data[next]
	}

	// Decrement the tail index and count
	rb.Tail = (rb.Tail - 1 + rb.Size) % rb.Size
	rb.Count--

	return true
}
