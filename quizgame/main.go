package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]

	// Set default timer
	quizTime := 30
	if len(args) >= 2 {
		var err error
		quizTime, err = strconv.Atoi(os.Args[2])
		if err != nil {
			exit("Timer must be a number: " + err.Error())
		}
	}

	// Check if the user has provided a csv file
	fileName := "problems.csv"
	if len(args) >= 1 {
		fileName = os.Args[1]
	}

	// Read the csv file
	file, err := os.Open(fileName)
	if err != nil {
		exit("Failed to open the CSV file: " + fileName)
	}

	r := csv.NewReader(file)

	records, err := r.ReadAll()

	if err != nil {
		exit(err.Error())
	}

	totalQuestions := len(records)
	correctAnswers := 0

	fmt.Println("Starting Quiz...")

	timer := time.NewTimer(time.Duration(quizTime) * time.Second)

	for i, record := range records {
		fmt.Println("Question #", i+1, ": ", strings.TrimSpace(record[0]))
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Printf("Answer: ")
			fmt.Fscanf(os.Stdin, "%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			printSummary(totalQuestions, correctAnswers)
			return
		case answer := <-answerCh:
			if answer == record[1] {
				correctAnswers++
			}
		}
	}

	// Summary
	printSummary(totalQuestions, correctAnswers)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func printSummary(totalQuestions, correctAnswers int) {
	fmt.Println("Total Questions: ", totalQuestions)
	fmt.Println("Score: ", correctAnswers, "/", totalQuestions)
}

// func startTimer(seconds int) {

// 	var counter int = seconds

// 	defer timer.Stop()
// 	done := make(chan bool)
// 	go func() {
// 		time.Sleep(time.Duration(seconds) * time.Second)
// 		done <- true
// 	}()
// 	for {
// 		select {
// 		case <-done:
// 			exit("Time over!")
// 		case <-timer.C:
// 			fmt.Println("\nSeconds left: ", counter)
// 			counter--
// 		}

// 	}
// }
