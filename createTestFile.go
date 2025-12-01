package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
)

func main() {

	filePath := "bench/ips.txt"
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer f.Close()

	var cnt int64 = 0
	setRandom := true
	if setRandom {

		for ; cnt < 100_000_000; cnt++ {
			oct1 := 1 + rand.IntN(256)
			oct2 := 1 + rand.IntN(256)
			oct3 := 1 + rand.IntN(256)
			oct4 := 1 + rand.IntN(256)
			_, err := fmt.Fprintf(f, "%d.%d.%d.%d\n", oct1, oct2, oct3, oct4)
			if err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
		}
	} else {
		for oct1 := 1; oct1 < 256; oct1++ {
			for oct2 := 1; oct2 < 256; oct2++ {
				for oct3 := 1; oct3 < 256; oct3++ {
					for oct4 := 1; oct4 < 256; oct4++ {
						if cnt >= 100_000_000 {
							return
						}
						_, err := fmt.Fprintf(f, "%d.%d.%d.%d\n", oct1, oct2, oct3, oct4)
						if err != nil {
							log.Fatalf("Failed to write to file: %v", err)
						}
					}
				}
			}
		}
	}

	fmt.Println("Done")
}
