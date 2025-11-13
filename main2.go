package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	zerolog := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	logger := slog.New(
		slogzerolog.Option{Logger: &zerolog}.NewZerologHandler(),
	)

	filePath := flag.String("file", "nan", "a path to the file to be processed")
	//useHLL := flag.Bool("accuracy", false, "if true - we calculate with 100% accuracy, else - we estimate number of unique IPs")
	flag.Parse()

	start := time.Now()
	//rawIp := make(chan []byte, 256)
	// read file line by line
	// parallel reading brings complexity
	file, err := os.Open(*filePath)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	// 1MB buffer minimize system calls
	reader := bufio.NewReaderSize(file, 1<<20)
	sketch, err := hyperloglog.NewSketch(18, false)
	cutSize := 0
	for {
		ip, err := reader.ReadBytes('\n')
		if err == io.EOF {
			logger.Info("EOF", slog.Any("IP", ip))
			break
		}
		if err != nil {
			os.Exit(1)
		}
		if ip[cap(ip)-1] == '\n' {
			cutSize++
		}
		if ip[cap(ip)-2] == '\r' {
			cutSize++
		}
		sketch.Insert(ip[:cap(ip)-cutSize])
		cutSize = 0
	}

	logger.Info("DONE", slog.Any("count", sketch.Estimate()), slog.Any("elapsed", time.Since(start).Seconds()))
}

// 342 172 175
// 341 704 755
// just reading - 11-12 sec
// with Triming = +1s
