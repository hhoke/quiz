package main

// Finished
import (
	"bufio"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type quiz struct {
	totalQuestions int
	questionsRight int
}

func (q *quiz) grade(record []string, rcChan chan int) {

	reader := bufio.NewReader(os.Stdin)
	q.totalQuestions++
	log.Print(record[0])
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s = strings.Replace(s, "\n", "", -1)
	if s == record[1] {
		q.questionsRight++
		rcChan <- 0
	} else {
		log.Print(record[1])
		log.Print(s)
		rcChan <- 1
	}

}

// Main just reads in a csv and prints it, for now
func main() {
	log.Printf("vim-go")
	file := "./problems.csv"
	timeLimit := flag.Int("limit", 1000, "the time limit for the quiz in seconds")
	fileDescriptor, err := os.Open(file)
	if err != nil {
		log.Printf("error opening file %s: %s", file, err.Error())
	}
	reader := csv.NewReader(fileDescriptor)
	q := quiz{}
	start_time := time.Now()
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		rcChan := make(chan int)
		go q.grade(record, rcChan)

		select {
		case <-timer.C:
			log.Print(q.totalQuestions)
			log.Print(q.questionsRight)
			log.Print(*timeLimit)
			return
		case <-rcChan:
			continue
		}
	}
	elapsed := time.Since(start_time)
	log.Print(q.totalQuestions)
	log.Print(q.questionsRight)
	log.Print(elapsed)
	return
}
