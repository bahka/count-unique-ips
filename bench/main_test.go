package bench

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/axiomhq/hyperloglog"
)

func BenchmarkFile(b *testing.B) {
	filePath := "ips.txt"
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReaderSize(file, 1<<20)
	for n := 0; n < b.N; n++ {
		sketch, _ := hyperloglog.NewSketch(18, false)
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			slog.Error("error resetting file reader", slog.String("error", err.Error()))
		}

		//s := cap(rawIp) / 8
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				slog.Error("error reading line", slog.String("error", err.Error()))
			}
			line = bytes.TrimSpace(line)
			parsedIp := net.ParseIP(string(line))
			if parsedIp != nil && (parsedIp.To4() != nil || parsedIp.To16() != nil) {
				sketch.Insert(parsedIp)
			}
		}
		fmt.Println(sketch.Estimate())
	}
}
