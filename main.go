package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

func main() {
	log.SetFlags(0)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to read current working dir :", err)
	}

	file, err := os.Open(filepath.Join(wd, "words.txt"))
	if err != nil {
		log.Fatalln("failed to open file :", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln("failed to close file :", err)
		}
	}()

	wordCount := CountWords(file)
	fmt.Println("Word Count =>", wordCount)
}

func countWords(data []byte) int {
	return len(bytes.Fields(data))
}

func CountWords(reader io.Reader) int {
	wordCount := 0
	bufferSize := 4096
	isInsideWord := false
	buf := make([]byte, bufferSize)

	for {
		size, err := reader.Read(buf)
		if err != nil {
			break
		}

		isInsideWord = !unicode.IsSpace(rune(buf[0])) && isInsideWord

		bufferCount := countWords(buf[:size])
		if isInsideWord {
			bufferCount--
		}

		wordCount += bufferCount
		isInsideWord = !unicode.IsSpace(rune(buf[size-1]))
	}

	return wordCount
}
