package gotail

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
)

type Tail struct {
	filename string
	done     chan bool
	offset   int64
}

func NewTail(filename string) Tail {
	return Tail{
		filename: filename,
		done:     make(chan bool),
		offset:   0,
	}
}

func (t *Tail) Close() {
	close(t.done)
}

func (t *Tail) ToEnd() {
	size := ReadSize(t.filename)
	t.offset = size
}

func (t *Tail) ReadData(filter string, read func(string)) {
	t.offset = ReadByLine(t.filename, t.offset, func(line string) {
		if filter != "" {
			if strings.Contains(line, filter) {
				read(line)
			}
		} else {
			read(line)
		}
	})
}

func (t *Tail) TailLine(filter string, read func(string)) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("NewWatcher failed: ", err)
	}

	defer func() {
		if err := recover(); err != nil {

		}
		watcher.Close()
		t.Close()
	}()

	err = watcher.Add(t.filename)
	if err != nil {
		log.Println("Add failed:", err)
	}

	for {
		select {
		case <-t.done:
			log.Println("--------------------end!!!!")
			return
		case _, ok := <-watcher.Events:
			if !ok {
				return
			}
			t.ReadData(filter, read)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}
