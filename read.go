package gotail

import (
	"bufio"
	"io"
	"os"
)

func ReadByLine(filename string, offset int64, output func(string)) int64 {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	_, err = f.Seek(offset, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF { // 读到文件末尾
			break
		} else {
			output(string(line))
		}
	}

	return info.Size()
}

func ReadByBytes(filename string, offset int64, output func([]byte, int)) int64 {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	_, err = f.Seek(offset, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	r := bufio.NewReader(f)
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	b := make([]byte, 1024)
	for {
		i, err := r.Read(b)
		if err == io.EOF {
			break
		} else {
			output(b, i)
		}
		b = b[:]
	}

	return info.Size()
}

func ReadSize(filename string) int64 {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	return info.Size()
}
