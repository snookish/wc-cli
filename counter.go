package counter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/iamBelugax/wc-cli/display"
)

const (
	bufSize = 4096
)

type Counts struct {
	Words int
	Lines int
	Bytes int
}

func (c *Counts) Print(w io.Writer, opts display.Options, suffixes ...string) {
	result := make([]string, 0)

	if opts.ShowWords() {
		result = append(result, strconv.Itoa(c.Words))
	}

	if opts.ShowLines() {
		result = append(result, strconv.Itoa(c.Lines))
	}

	if opts.ShowBytes() {
		result = append(result, strconv.Itoa(c.Bytes))
	}

	line := strings.Join(result, "\t") + "\t"
	fmt.Fprint(w, line)

	suffixStr := strings.Join(suffixes, " ")
	if suffixStr != "" {
		fmt.Fprint(w, " ", suffixStr)
	}

	fmt.Fprintln(w)
}

func (c *Counts) Add(other Counts) {
	c.Bytes += other.Bytes
	c.Lines += other.Lines
	c.Words += other.Words
}

func CountAll(r io.Reader) Counts {
	var counts Counts
	var isInsideWord bool
	br := bufio.NewReaderSize(r, bufSize)

	for {
		ch, size, err := br.ReadRune()
		if err != nil {
			break
		}

		if !unicode.IsSpace(ch) && !isInsideWord {
			counts.Words++
		}

		if ch == '\n' {
			counts.Lines++
		}

		counts.Bytes += size
		isInsideWord = !unicode.IsSpace(ch)
	}

	return counts
}

func CountFile(path string) (Counts, error) {
	f, err := os.Open(path)
	if err != nil {
		return Counts{}, err
	}
	defer f.Close()

	return CountAll(f), nil
}

// CountLines counts newline characters.
func CountLines(r io.Reader) int {
	br := bufio.NewReaderSize(r, bufSize)
	n := 0

	for {
		ch, _, err := br.ReadRune()
		if err != nil {
			break
		}

		if ch == '\n' {
			n++
		}
	}

	return n
}

// CountBytes counts the total number of bytes in the reader.
func CountBytes(r io.Reader) int {
	n, _ := io.Copy(io.Discard, r)
	return int(n)
}

// CountWords counts words using a buffered scanner.
func CountWords(r io.Reader) int {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	n := 0
	for s.Scan() {
		_ = s.Text()
		n++
	}

	return n
}

// CountWordsBuf counts words using bufio.Reader and checking for whitespace.
func CountWordsBuf(r io.Reader) int {
	br := bufio.NewReaderSize(r, bufSize)

	n := 0
	isInsideWord := false

	for {
		ch, _, err := br.ReadRune()
		if err != nil {
			break
		}

		// If the rune is not a space and we're not inside a word, it's the start of a new word
		if !unicode.IsSpace(ch) && !isInsideWord {
			n++
		}

		isInsideWord = !unicode.IsSpace(ch)
	}

	return n
}

// CountWordsRaw is a custom implementation for counting words.
func CountWordsRaw(r io.Reader) int {
	buf := make([]byte, bufSize)
	leftover := make([]byte, 0)

	n := 0
	isInsideWord := false

	// Buffers for reading input and storing leftover bytes.
	for {
		k, err := r.Read(buf)
		if err != nil {
			break
		}

		// Append any leftover bytes from previous reads to the current buffer.
		data := append(leftover, buf[:k]...)

		for len(data) > 0 {
			ch, size := utf8.DecodeRune(data)
			if ch == utf8.RuneError {
				break
			}

			// If the rune is not whitespace and we're not inside a word, it's a new word.
			if !unicode.IsSpace(ch) && !isInsideWord {
				n++
			}

			isInsideWord = !unicode.IsSpace(ch)
			data = data[size:]
		}

		// Store any leftover bytes that didn't form a complete rune.
		leftover = leftover[:0]
		leftover = append(leftover, data...)
	}

	return n
}

func CountAllTeaReader(r io.Reader) Counts {
	buf1 := bytes.Buffer{}
	buf2 := bytes.Buffer{}

	linesReader := io.TeeReader(r, &buf1)
	wordsReader := io.TeeReader(&buf1, &buf2)
	bytesReader := &buf2

	linesCount := CountLines(linesReader)
	wordsCount := CountWords(wordsReader)
	bytesCount := CountBytes(bytesReader)

	return Counts{
		Lines: linesCount,
		Words: wordsCount,
		Bytes: bytesCount,
	}
}

func CountAllIOPipe(r io.Reader) Counts {
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()

	linesReader := io.TeeReader(r, pw1)
	wordsReader := io.TeeReader(pr1, pw2)
	bytesReader := pr2

	// var linesCount, wordsCount, bytesCount int
	// var wg sync.WaitGroup
	// wg.Add(3)

	// go func() {
	// 	defer func() {
	// 		wg.Done()
	// 		pw1.Close()
	// 	}()
	// 	linesCount = CountLines(linesReader)
	// }()

	// go func() {
	// 	defer func() {
	// 		wg.Done()
	// 		pw2.Close()
	// 	}()
	// 	wordsCount = CountWords(wordsReader)
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	bytesCount = CountBytes(bytesReader)
	// }()

	linesCh := make(chan int)
	wordsCh := make(chan int)
	bytesCh := make(chan int)

	go func() {
		defer pw1.Close()
		defer close(linesCh)
		linesCh <- CountLines(linesReader)
	}()

	go func() {
		defer pw2.Close()
		defer close(wordsCh)
		wordsCh <- CountWords(wordsReader)
	}()

	go func() {
		defer close(bytesCh)
		bytesCh <- CountBytes(bytesReader)
	}()

	linesCount := <-linesCh
	wordsCount := <-wordsCh
	bytesCount := <-bytesCh

	return Counts{
		Words: wordsCount,
		Lines: linesCount,
		Bytes: bytesCount,
	}
}

func CountAllMultiWriter(r io.Reader) Counts {
	linesReader, linesWriter := io.Pipe()
	wordsReader, wordsWriter := io.Pipe()
	bytesReader, bytesWriter := io.Pipe()
	w := io.MultiWriter(linesWriter, wordsWriter, bytesWriter)

	linesCh := make(chan int)
	wordsCh := make(chan int)
	bytesCh := make(chan int)

	go func() {
		defer close(linesCh)
		linesCh <- CountLines(linesReader)
	}()

	go func() {
		defer close(wordsCh)
		wordsCh <- CountWords(wordsReader)
	}()

	go func() {
		defer close(bytesCh)
		bytesCh <- CountBytes(bytesReader)
	}()

	io.Copy(w, r)
	wordsWriter.Close()
	linesWriter.Close()
	bytesWriter.Close()

	linesCount := <-linesCh
	wordsCount := <-wordsCh
	bytesCount := <-bytesCh

	return Counts{
		Words: wordsCount,
		Lines: linesCount,
		Bytes: bytesCount,
	}
}
