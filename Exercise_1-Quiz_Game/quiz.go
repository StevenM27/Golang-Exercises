package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}


func main() {

	// Create our flags for the file name and the time limit
	filePtr := flag.String("file", "problems.csv", "file path/name")
	timePtr := flag.Int("time", 30, "time limit for quiz")
	flag.Parse()

	fmt.Println(*filePtr)

	// Get raw data and check for errors
	data, err := os.ReadFile(string(*filePtr))
	check(err)

	// Create csv reader to read raw data and parse into usable key/value pairs
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	// Create scanner to read standard input
	scanner := bufio.NewScanner(os.Stdin)
	
	runningScore := 0

	// Ask the user if they are ready to start. Then, start the timer
	fmt.Print("Press Enter to start: ")
	scanner.Scan()


	// Create goroutine to time the user
	done := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(*timePtr) * time.Second)
		done <- true
	}()


	// Create channel for user input
	input := make(chan string, 1)


	questionLoop:
	for _, line := range records {
		question, answer := line[0], line[1]
		fmt.Print(question, ": ")

		// Goroutine to listen to user input
		go func() {
			scanner.Scan()
			input <- scanner.Text()
		}()

		// Listens for either the user's input or the timer, whichever comes first
		select {
		case <-done:
			fmt.Println("\nTime's up!")
			break questionLoop
		case <-input:
			if scanner.Text() == answer {
				runningScore++
			}
		}

		
	}

	fmt.Println("Your final score is ", runningScore, "/", len(records))


}