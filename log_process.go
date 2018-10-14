package main

import (
	"strings"
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
)

type LogProcess struct {
	rc chan []byte
	wc chan string
	reader  Reader
	writer  Writer
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Schema, Status string
	UpstreamTime, RequestTime    float64
}

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type ReadFromFile struct {
	path string // 读取文件的路径
}

type WriteToInfluxDB struct {
	influxDBDsn string // influx db source
}

func (r *ReadFromFile) Read(rc chan []byte)  {
	// 读取模块
	f, err := os.Open(r.path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500*time.Millisecond)
			continue
		} else if err != nil {
			panic(err)
		}

		rc <- line[:len(line) - 1]
	}

}

func (l *LogProcess) Process()  {
	// 解析模块
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}
}

func (w *WriteToInfluxDB) Write(wc chan string)  {
	for v := range wc {
		fmt.Println(v)
	}
}

func main() {
	lp := &LogProcess{
		rc: make(chan []byte, 200),
		wc: make(chan string, 200),
		reader: &ReadFromFile{
			path: "./access.log",
		},
		writer: &WriteToInfluxDB{
			influxDBDsn: "username&password...",
		},
	}

	go lp.reader.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)

	time.Sleep(30*time.Second)
}
