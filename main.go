package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Fake a long and difficult work.
func processRequest(duration int) {
	time.Sleep(time.Duration(duration * int(time.Second)))
}

func main() {
	// Usage: ./GoroutinesThrottling -c=2 7 6 4 3 7 3 4 6 4
	// This will launch a program with 9 goroutines with only 2 goroutines at any time.

	// Parsing the flags from command line.
	maxGoroutinesAtOnce := flag.Int("c", 5, "the number of goroutines that will run at once")
	flag.Parse()
	durationArgs := flag.Args()
	jobsCount := len(flag.Args())

	// Channels for managing the jobs
	finishChannel := make(chan bool)
	interruptChannel := make(chan os.Signal, 1)
	waitForJobsChannel := make(chan bool)
	goroutinesChannel := make(chan struct{}, *maxGoroutinesAtOnce)

	fmt.Printf("Ready to go !\n")
	fmt.Printf("Number of jobs : %d \n", jobsCount)
	fmt.Printf("Number of goroutines at once : %d \n", *maxGoroutinesAtOnce)

	for i := 0; i < *maxGoroutinesAtOnce; i++ {
		goroutinesChannel <- struct{}{}
	}

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// When job is finished, let another job to run.
	go func() {

		for i := 0; i < jobsCount; i++ {
			<-finishChannel
			goroutinesChannel <- struct{}{}
			if len(interruptChannel) > 0 {
				fmt.Println("\r- Ctrl+C interrupt pressed in Terminal")
				break
			}
		}
		waitForJobsChannel <- true
	}()

	// The job is waiting to launch until the goroutineChannel will be filled again
	// Running processRequest function with duration from command-line arguments.
	for i := 0; i < jobsCount; i++ {
		if len(interruptChannel) == 0 {
			fmt.Printf("Job ID: %v: Job waiting\n", i)
			<-goroutinesChannel
			if len(interruptChannel) == 0 {
				fmt.Printf("Job ID: %v: Job is running\n", i)
				go func(id int) {
					duration, _ := strconv.Atoi(durationArgs[id])
					processRequest(duration)
					fmt.Printf("Job ID: %v: Job finished!\n", id)
					finishChannel <- true
				}(i)
			}
		}
	}
	<-waitForJobsChannel
}
