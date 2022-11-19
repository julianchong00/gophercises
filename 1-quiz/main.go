package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

const (
	// DefaultTimeLimit is the default time limit for the quiz in seconds
	DefaultTimeLimit = 30
)

type Config struct {
	// TimeLimit is the time limit for the quiz in seconds
	TimeLimit    int
	ProblemsFile string
}

func readProblems(csvFile string) []Problem {
	// Open the file
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Couldn't open the csv file: %v", err)
	}

	problemList := []Problem{}

	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		problemList = append(problemList, Problem{line[0], line[1]})
	}

	return problemList
}

func main() {
	// Set the default configuration
	config := Config{
		TimeLimit:    DefaultTimeLimit,
		ProblemsFile: "problems.csv",
	}

	// Parse the command line flags
	flag.IntVar(
		&config.TimeLimit,
		"limit",
		DefaultTimeLimit,
		"the time limit for the quiz in seconds",
	)
	flag.StringVar(
		&config.ProblemsFile,
		"csv",
		"problems.csv",
		"a csv file in the format of 'question,answer'",
	)
	flag.Parse()

	// Read the problems from the csv file
	problems := readProblems(config.ProblemsFile)

	// Wait for user to press enter before starting the quiz timer
	fmt.Print("Press enter to start the quiz...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Keep track of score
	score := 0

	// Start timer
	// Timer sends a message on the channel after the specified duration
	timer := time.NewTimer(time.Duration(config.TimeLimit) * time.Second)
	go func() {
		// Block until the timer receives a message on its channel after time limit
		<-timer.C
		fmt.Println("\nTime's up!")
		fmt.Printf("You scored %d out of %d.\n", score, len(problems))
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i, problem.Question)
		answer, _ := reader.ReadString('\n')

		if strings.TrimSpace(answer) == problem.Answer {
			score++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}
