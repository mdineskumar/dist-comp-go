package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type quoteReader struct {
	state int
}

type Status int

const (
	NeedMoreInput Status = iota
	Success
	BadInput
)

func (q *quoteReader) Init() {
	q.state = 0
}

func makeBufioReader(file io.Reader) func() rune {
	reader := bufio.NewReader(file)

	return func() rune {
		char, _, err := reader.ReadRune()

		if err != nil {
			return 0
		}
		return char
	}
}

func readString(readChar func() rune) bool {
	if readChar() != '"' {
		return false
	}

	var c rune
	for c != '"' {
		c = readChar()
		if c == '\\' {
			readChar()
		}
	}

	return true
}

func (q *quoteReader) ProcessChar(c rune) Status {
	switch q.state {
	case 0:
		if c != '"' {
			return BadInput
		}
		q.state = 1
	case 1:
		if c == '"' {
			return Success
		}

		if c == '\\' {
			q.state = 2
		} else {
			q.state = 1
		}
	case 2:
		q.state = 1
	}
	return NeedMoreInput
}

func main() {

	// testCases := []string{
	// 	`"hello"`,           // Valid
	// 	`"hello \"world\""`, // Valid (escaped quotes)
	// 	`"C:\\Path"`,        // Valid (escaped backslash)
	// 	`hello"`,            // Invalid (missing opening quote)
	// 	`"hello`,            // Invalid (missing closing quote)
	// 	`"hello\`,           // Invalid (hanging escape character)
	// }

	fileName := "test_large_file.txt"

	os.WriteFile(fileName, []byte(`"This could be a \massive string \reading off a disk!"`), 0644)
	//defer os.Remove(fileName)

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()

	myReaderChar := makeBufioReader(file)

	isValid := readString(myReaderChar)

	fmt.Printf("Was the string in the file parsed successfully? %v\n", isValid)

	//fmt.Println("--- Testing State Machine ---")
	// for _, text := range testCases {
	// 	qr := quoteReader{}
	// 	qr.Init()

	// 	var status Status
	// 	for _, char := range text {
	// 		status = qr.ProcessChar(char)
	// 		if status != NeedMoreInput {
	// 			break
	// 		}
	// 	}

	// 	fmt.Printf("Input: %-20s | Status: %v\n", text, status == Success)

	// }
}
