package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log/slog"
	"os"
	"time"
)

//func isValidIp(s []byte) bool {
//	ip := net.ParseIP(string(s))
//	return ip != nil && (ip.To4() != nil || ip.To16() != nil)
//}

func mainy() {

	filePath := flag.String("file", "nan", "a path to the file to be processed")
	flag.Parse()

	start := time.Now()
	file, err := os.Open(*filePath)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	// 1MB buffer minimize system calls
	reader := bufio.NewReaderSize(file, 1<<20)
	//sketch, err := hyperloglog.NewSketch(18, false)
	cutSize := 0
	for {
		ip, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			os.Exit(1)
		}
		//if ip[cap(ip)-1] == '\n' {
		cutSize++
		//}
		//if ip[cap(ip)-2] == '\r' {
		//	cutSize++
		//}
		//ip = ip[:len(ip)-cutSize]
		ip = bytes.TrimSpace(ip)
		if !isValidIp(ip) {
			slog.Error("Invalid IP address", slog.Any("ip", ip))
			continue
		}
		//sketch.Insert(ip)
		cutSize = 0
	}

	slog.Info("DONE", slog.Any("count", 1), slog.Any("elapsed", time.Since(start).Seconds()))
}

// 342 172 175
// 341 704 755
// just reading - 11-12 sec
// with Triming = +1s
