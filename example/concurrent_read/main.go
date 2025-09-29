/* Try concurrent read example: go run example/concurrent_read/main.go */

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	h "github.com/viktorxia/hgnc-go"
)

func main() {
	// Load HGNC database
	fmt.Println("Loading HGNC Database...")
	hgnc, err := h.LoadTsv("data/hgnc_complete_set.txt.gz", true)
	if err != nil {
		log.Fatalf("Failed to load HGNC data: %v", err)
	}
	fmt.Println("Database loaded successfully\n")

	// Test different scales
	fmt.Println("=== Concurrent Read Stress Tests ===")
	fmt.Printf("Available CPU cores: %d\n\n", runtime.NumCPU())

	// Small scale warmup
	fmt.Println("🔥 Warmup Test (100K operations)")
	stressTest(hgnc, 100000, []int{1, 2, 4})

	// Medium scale
	fmt.Println("\n🚀 Medium Scale Test (1M operations)")
	stressTest(hgnc, 1000000, []int{1, 2, 4, 6, 8, 12})

	// Large scale
	fmt.Println("\n💥 Large Scale Test (5M operations)")
	stressTest(hgnc, 5000000, []int{1, 2, 4, 6, 8, 12})

	fmt.Println("\n✅ All stress tests completed successfully!")
}

// getAllSymbolsFromFile reads the TSV file directly and extracts symbols from the second column
func getAllSymbolsFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filename, err)
		return []string{"TP53", "BRCA1", "EGFR", "KRAS"} // fallback
	}
	defer file.Close()

	// Handle gzip compression
	var scanner *bufio.Scanner
	if strings.HasSuffix(filename, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			log.Printf("Failed to create gzip reader: %v", err)
			return []string{"TP53", "BRCA1", "EGFR", "KRAS"} // fallback
		}
		defer gzReader.Close()
		scanner = bufio.NewScanner(gzReader)
	} else {
		scanner = bufio.NewScanner(file)
	}

	var symbols []string
	symbolSet := make(map[string]bool)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Skip header line
		if lineCount == 1 {
			continue
		}

		// Split by tab and get the second field (index 1)
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 {
			symbol := strings.TrimSpace(fields[1])
			if symbol != "" && symbol != "-" {
				if !symbolSet[symbol] {
					symbolSet[symbol] = true
					symbols = append(symbols, symbol)
				}
			}
		}
	}

	if len(symbols) == 0 {
		return []string{"TP53", "BRCA1", "EGFR", "KRAS"} // fallback
	}

	return symbols
}

// stressTest performs stress testing with different goroutine counts
func stressTest(hgnc *h.HGNC, totalOps int, goroutineCounts []int) {
	allSymbols := getAllSymbolsFromFile("data/hgnc_complete_set.txt.gz")
	fmt.Printf("Using %d gene symbols for testing\n", len(allSymbols))

	fmt.Printf("%-12s %-15s %-15s %-12s %-8s\n",
		"Goroutines", "Time", "Ops/sec", "Avg(ns)", "Speedup")
	fmt.Printf("%-12s %-15s %-15s %-12s %-8s\n",
		"----------", "----", "-------", "--------", "-------")

	var baselineTime time.Duration

	for i, numGoroutines := range goroutineCounts {
		// Skip if requesting more goroutines than CPU cores (except for 1)
		if numGoroutines > runtime.NumCPU() && numGoroutines != 1 {
			continue
		}

		duration := runStressTest(hgnc, allSymbols, totalOps, numGoroutines)

		opsPerSec := float64(totalOps) / duration.Seconds()
		avgNs := duration.Nanoseconds() / int64(totalOps)

		var speedupStr string
		if i == 0 {
			baselineTime = duration
			speedupStr = "1.00x"
		} else {
			speedup := baselineTime.Seconds() / duration.Seconds()
			speedupStr = fmt.Sprintf("%.2fx", speedup)
		}

		// Format time nicely
		timeStr := formatDuration(duration)

		fmt.Printf("%-12d %-15s %-15.0f %-12d %-8s\n",
			numGoroutines, timeStr, opsPerSec, avgNs, speedupStr)
	}
}

// runStressTest executes the actual stress test
func runStressTest(hgnc *h.HGNC, symbols []string, totalOps, numGoroutines int) time.Duration {
	var wg sync.WaitGroup
	var completedOps int64

	opsPerGoroutine := totalOps / numGoroutines
	remainingOps := totalOps % numGoroutines

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		// Distribute remaining operations to first few goroutines
		ops := opsPerGoroutine
		if i < remainingOps {
			ops++
		}

		go func(goroutineID, numOps int) {
			defer wg.Done()

			for j := 0; j < numOps; j++ {
				// Use different symbols for better coverage
				symbolIdx := (goroutineID*numOps + j) % len(symbols)
				gene := symbols[symbolIdx]

				// Mix different operations for comprehensive testing
				switch j % 7 {
				case 0:
					hgnc.IsCodingGene(gene)
				case 1:
					hgnc.SymbolToEntrezID(gene)
				case 2:
					hgnc.GetManeSelect(gene)
				case 3:
					hgnc.Fetch(gene, h.FIELD_SYMBOL)
				case 4:
					hgnc.Lookup(gene, h.FIELD_SYMBOL, h.FIELD_ENTREZ_ID)
				case 5:
					hgnc.SymbolToEnsg(gene)
				case 6:
					// Test another fetch operation
					hgnc.Fetch(gene, h.FIELD_ENTREZ_ID)
				}

				atomic.AddInt64(&completedOps, 1)
			}
		}(i, ops)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// Verify we completed all operations
	if atomic.LoadInt64(&completedOps) != int64(totalOps) {
		log.Printf("Warning: Expected %d operations, completed %d",
			totalOps, atomic.LoadInt64(&completedOps))
	}

	return duration
}

// formatDuration formats duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d >= time.Second {
		return fmt.Sprintf("%.3fs", d.Seconds())
	} else if d >= time.Millisecond {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
	} else if d >= time.Microsecond {
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1e3)
	} else {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
}
