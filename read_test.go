package gotail

import (
	"log"
	"testing"
)

func TestReadByLine(t *testing.T) {
	filename := "/Users/mac/Downloads/log.txt"
	var offset int64 = 0
	ReadByLine(filename, offset, func(line string) {
		log.Println(line)
	})
}

func TestReadByBytes(t *testing.T) {
	filename := "/Users/mac/Downloads/log.txt"
	var offset int64 = 0
	ReadByBytes(filename, offset, func(bytes []byte, i int) {
		log.Printf(string(bytes))

	})
}
