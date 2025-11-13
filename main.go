package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
	_ "net/http/pprof"
)

var filePath = flag.String("file", "nan", "a path to the file to be processed")

func processFile(filePath *string, rawIp chan [][]byte) error {
	file, err := os.Open(*filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 1MB buffer minimize system calls
	reader := bufio.NewReaderSize(file, 1<<20)
	//s := cap(rawIp) / 8
	batch := make([][]byte, 0, 64)
	for {
		ip, err := reader.ReadBytes('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		batch = append(batch, bytes.TrimSpace(ip))
		if len(batch) == cap(batch) {
			rawIp <- batch
			batch = make([][]byte, 0, cap(batch))
		}

	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	zerolog := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	logger := slog.New(
		slogzerolog.Option{Logger: &zerolog}.NewZerologHandler(),
	)
	flag.Parse()

	start := time.Now()
	rawIp := make(chan [][]byte, 2024)
	// read file line by line
	// parallel reading brings complexity
	go func() {
		defer close(rawIp)
		err := processFile(filePath, rawIp)
		if err != nil {
			logger.Error("error processing file", slog.Any("error", err))
			os.Exit(1)
		}
	}()
	var wg sync.WaitGroup
	goroutNum := 6
	sketches := make([]*hyperloglog.Sketch, goroutNum)
	for i := 0; i < goroutNum; i++ {
		wg.Add(1)
		sketches[i], _ = hyperloglog.NewSketch(18, false)
		go func() {
			defer wg.Done()
			for batch := range rawIp {
				for _, ip := range batch {
					sketches[i].Insert(ip)
				}
			}
		}()
	}

	// sparse > os.
	//sketch, _ := hyperloglog.NewSketch(14, false)
	//var wg sync.WaitGroup
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//	for ip := range rawIp {
	//		sketch.Insert(ip)
	//
	//		//logger.Info("Processing", slog.String("ip", string(s)))
	//	}
	//}()
	//sketch2, _ := hyperloglog.NewSketch(14, false)
	//go func() {
	//	defer wg.Done()
	//	for ip := range rawIp {
	//		sketch2.Insert(ip)
	//
	//		//logger.Info("Processing", slog.String("ip", string(s)))
	//	}
	//}()

	//for i := range rawIp {
	//	<-rawIp
	//}
	wg.Wait()
	//logger.Info("Start merge", slog.Any("elapsed", time.Since(start).Seconds()))
	for _, sketch := range sketches {
		sketches[0].Merge(sketch)
	}

	logger.Info("DONE", slog.Any("count", sketches[0].Estimate()), slog.Any("elapsed", time.Since(start).Seconds()))
} // 342 172 175
