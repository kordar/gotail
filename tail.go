package gotail

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
)

type Request struct {
	filter  string
	handler func(line string)
}

type Tail struct {
	filename string
	done     chan bool
	offset   int64
	workers  chan Request
	poolSize int
}

func NewTail(filename string, poolSize int) Tail {
	return Tail{
		filename: filename,
		done:     make(chan bool),
		offset:   0,
		workers:  make(chan Request, poolSize),
		poolSize: poolSize,
	}
}

func (t *Tail) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Printf("err = %v", err))
		}
	}()
	close(t.done)
}

func (t *Tail) ToEnd() {
	size := ReadSize(t.filename)
	t.offset = size
}

func (t *Tail) startPoolWork() {
	for {
		select {
		case <-t.done:
			return
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-t.workers:
			t.ReadData(request.filter, request.handler)
		}
	}
}

func (t *Tail) ReadData(filter string, read func(string)) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Printf("err = %v", err))
			return
		}
	}()
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
		log.Println("NewWatcher Failed: ", err)
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

	go t.startPoolWork()

	for {
		select {
		case <-t.done:
			log.Println("--------------------end!!!!")
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) {
				request := Request{filter: filter, handler: read}
				if len(t.workers) < t.poolSize {
					t.workers <- request
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}

}
