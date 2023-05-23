package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"unicode"
)

const defaultLimit int = 20

var errDataNotFound = errors.New("no data")

type (
	config struct {
		filePath   string
		wordsCount int
	}

	rawData [][]byte

	word struct {
		value []byte
		count int
	}

	wordList []word
)

func (r rawData) getWordList() wordList {

	words := make([]word, 0, len(r))

	for i := range r {
		if len(words) == 0 || !words[len(words)-1].isEqual(r[i]) {
			words = append(words, word{
				value: r[i],
				count: 1,
			})
			continue
		}
		words[len(words)-1].increase()
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count > words[j].count
	})
	return words
}

func (w *word) increase() {
	w.count++
}

func (w *word) isEqual(data []byte) bool {
	return bytes.Equal(w.value, data)
}

func (w *word) countLength() int {
	return len(fmt.Sprintf("%d", w.count))
}

func (w *word) printWithOffset(offset int) string {
	return fmt.Sprintf("%"+strconv.Itoa(offset)+"d %s", w.count, w.value)
}

func (w wordList) print(limit int) {
	if len(w) == 0 {
		return
	}
	offset := w[0].countLength()

	for i := range w {
		if i >= limit {
			break
		}
		fmt.Println(w[i].printWithOffset(offset))
	}
}

func getConfig() *config {
	filePath := flag.String("file_path", "./assets/mobydick.txt", "file path")
	wordsCount := flag.Int("words_count", 20, "count of words")
	flag.Parse()
	if *wordsCount < 0 {
		*wordsCount = defaultLimit
	}
	return &config{
		filePath:   *filePath,
		wordsCount: *wordsCount,
	}
}

func getSortedRawData(buf *bytes.Buffer) (rawData, error) {
	data := bytes.FieldsFunc(buf.Bytes(), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	if len(data) == 0 {
		return nil, errDataNotFound
	}
	for i := range data {
		data[i] = bytes.ToLower(data[i])
	}
	sort.Slice(data, func(i, j int) bool {
		return bytes.Compare(data[i], data[j]) < 0
	})
	return data, nil
}

func main() {

	conf := getConfig()

	file, err := os.Open(conf.filePath)
	if err != nil {
		log.Printf("[ERROR] Open: %s", err.Error())
		return
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Printf("[ERROR] ReadFrom: %s", err.Error())
		return
	}

	data, err := getSortedRawData(buf)
	if err != nil {
		log.Printf("[ERROR] getSortedRawData: %s", err.Error())
		return
	}

	words := data.getWordList()

	words.print(conf.wordsCount)
}
