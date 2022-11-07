package gotail

import (
	"log"
	"testing"
	"time"
)

func TestTailLog(t *testing.T) {
	filename := "/Users/mac/Downloads/log.txt"
	tail := NewTail(filename)
	//tail.ReadData("", func(s string) {
	//	log.Println(s)
	//})
	tail.ToEnd()
	go tail.TailLine("12", func(s string) {
		log.Println(s)
	})
	time.Sleep(time.Duration(30) * time.Second)
	tail.Close()
}
