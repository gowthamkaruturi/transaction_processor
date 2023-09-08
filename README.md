**Introduction**

The transaction_processor designed in the assignment would deal with all the questions asked.

1. The design consists of three packages
    config package ---> the code in this package is loaded during the start of the main function.
   
        * config contains three parameters in it
        * Workers --> This is the config parameter used for number of parallel processors/workers in go 
        * Capacity --> this is the configurable circular queue size
        FilePath --> file path , where the json is stored. This is a     static location in the root director. This can be anything like s3 or a remote location , necessary tooling and code hasto be completed for this to be done.

   models package -->  the models package contains the needed code for all circular queue  

       * RingerBuffer is designed to hold transaction and processed in the order of circular queue's enqueue and dequeue 
       * Enqueue --> put the transaction in to the circular queue.
       * Dequeue --> remove the transaction from the circular queue.
       * DequeuWithTrasaction --> remove specific transaction from the circular queue which takes o(n)
   
     main package --> the processing logic for transactions are done here.
        
        * transactions are designed as buffered channels so multiple readers(go routines) can process the messages that come in to channel.
        * processTransaction is treated as a goroutine/worker which would do processing of putting it to standard output by reading from the worker channel and remove it from the circular queue.
        * processTransactions is the method which would do the needed task of reading the data from the file and put it to the circular queue. The processing order of the transaction is the order the transaction came in file.

  Multiple workers/go routines process the data from the queue, so the concurreny is handled using routines, config is designed to configure the number of workers and capacity dynamically from a file. Hence both the problem1 and problem2  are solved.    

Test cases :
  
  1. Implemented the test cases for 

         config , model and main pkgs
            following is coverage  output 
            ok      github.com/transaction_processor        (cached)        coverage: 44.7% of statements
            ok      github.com/transaction_processor/config 0.147s  coverage: 80.0% of statements
            ok      github.com/transaction_processor/models 0.103s  coverage: 97.1% of statements 

**Building and testing could be done :**
1. `go build .` to build the code , genarate executable to run
2. `go run .` to build and run the code 
3. `go test -cover ./...`  to run the test cases with coverage. 


**Bottlenecks:**

In the given code, there are a few potential bottlenecks in the design that could affect performance and efficiency. Here are a few areas to consider:

Ring Buffer Capacity: The capacity of the RingBuffer is determined by the config.Capacity value. Make sure this value is appropriately set based on the expected number of transactions to process. If the capacity is too low, it may result in dropped transactions if the buffer becomes full. If the capacity is too high, it may consume excessive memory.

File Reading: The code reads the transaction records from a file line by line using bufio.NewReader(). This approach could be inefficient for large files due to frequent I/O operations. Consider using buffered reading techniques like io.ReadFull() or reading the entire file into memory if the file size is manageable.

JSON Decoding: Each line read from the file is being decoded into a Transaction object using json.Unmarshal(). This operation can be expensive, especially if there are many records. Consider using a streaming JSON decoder like json.Decoder to avoid loading the entire file into memory before decoding.

Concurrency and Goroutine Management: The code uses goroutines to parallelize transaction processing. However, there is a fixed number of goroutines created based on the config.Workers value. Depending on the workload and available system resources, the number of workers can be optimized to achieve better performance.

To further optimize the code, we may consider:

Using a more efficient data structure or algorithm for the ring buffer.
Implementing error handling and logging mechanisms to track and handle errors during processing.
Considering other concurrent patterns like worker pools or task scheduling frameworks to manage and distribute work efficiently.   
