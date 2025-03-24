package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

func main() {
	file, errOpen := os.OpenFile("messages.txt", os.O_RDONLY, fs.ModeType)
	if errOpen != nil {
		log.Fatalf("error on open file err: %s", errOpen)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
	fmt.Printf("read: %s\n", "end")

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)

	go func() {
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("err while reading closing file err: %s", err)
			}
		}()

		var err error
		var parts []string
		for !errors.Is(err, io.EOF) {
			data := make([]byte, 8)
			_, err = f.Read(data)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Printf("err while reading from file err: %s", err)
			}

			nLine := "\n"
			dataStr := string(data)

			lineSegments := strings.Split(dataStr, nLine)

			hasLineCut := len(lineSegments) <= 1
			if hasLineCut {
				parts = append(parts, dataStr)
				continue
			}

			parts = append(parts, lineSegments[0])
			line := strings.Join(parts, "")

			linesChan <- line

			parts = nil
			parts = append(parts, lineSegments[1:]...)
		}

		close(linesChan)
	}()

	return linesChan
}
