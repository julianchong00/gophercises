package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
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
	DefaultTimeLimit    = 30
	DefaultProblemsFile = "problems.csv"
	DefaultShuffle      = false
)

type Config struct {
	// TimeLimit is the time limit for the quiz in seconds
	TimeLimit    int
	ProblemsFile string
	Shuffle      bool
}

func readProblems(csvFile string, shuffle bool) []Problem {
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

	// Shuffle the problems if shuffle flag is true
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(
			len(problemList),
			func(i, j int) { problemList[i], problemList[j] = problemList[j], problemList[i] },
		)
	}

	return problemList
}

func main() {
	// Set the default configuration
	config := Config{
		TimeLimit:    DefaultTimeLimit,
		ProblemsFile: DefaultProblemsFile,
		Shuffle:      DefaultShuffle,
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
		DefaultProblemsFile,
		"a csv file in the format of 'question,answer'",
	)
	flag.BoolVar(&config.Shuffle, "shuffle", DefaultShuffle, "shuffle the problems")
	flag.Parse()

	// Read the problems from the csv file
	problems := readProblems(config.ProblemsFile, config.Shuffle)

	// Wait for user to press enter before starting the quiz timer
	fmt.Print("Press enter to start the quiz...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Start timer
	// Timer sends a message on the channel after the specified duration
	timer := time.NewTimer(time.Duration(config.TimeLimit) * time.Second)

	// Keep track of score
	score := 0

	// Create reader to get user input
	reader := bufio.NewReader(os.Stdin)
problemLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.Question)
		// Make channel for answer so that the program isn't stuck waiting for
		// the user to enter an answer.
		// This allows the program to quit when the time has run out, even though
		// the user has not answered yet.
		answerCh := make(chan string)
		go func() {
			answer, _ := reader.ReadString('\n')
			answerCh <- answer
		}()

		select {
		// Listen for message on timer channel
		case <-timer.C:
			fmt.Println("\nTime's up!")
			break problemLoop

		// Listen for answer on answer channel
		case answer := <-answerCh:
			if strings.TrimSpace(answer) == strings.ToLower(strings.TrimSpace(problem.Answer)) {
				score++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}
