package gotail

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestTailLog(t *testing.T) {
	filename := "/Users/mac/Downloads/log.txt"
	tail := NewTail(filename, 2)
	//tail.ReadData("", func(s string) {
	//	log.Println(s)
	//})
	tail.ToEnd()
	go tail.TailLine("哈哈", func(s string) {
		log.Println(s)
	})

	f := func(string2 string) {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		// 及时关闭file句柄
		defer file.Close()
		// 写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		for i := 0; i < 5; i++ {
			write.WriteString(string2)
			write.Flush()
		}
		// Flush将缓存的文件真正写入到文件中
	}

	go f("Hello，C。 \r\n")
	go f("Hello，JAVA。哈哈 \r\n")
	go f("Hello，Python。 \r\n")
	go f("Hello，C#。 \r\n")

	time.Sleep(time.Duration(5) * time.Second)
	tail.Close()
}
