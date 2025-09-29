/* Try: go run example/cache_vs_nocache/main.go */

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	h "github.com/viktorxia/hgnc-go"
)

func main() {
	// Load HGNC database
	fmt.Println("Loading HGNC database...")
	hgnc, err := h.LoadTsv("data/hgnc_complete_set.txt.gz", true)
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}
	fmt.Println("Database loaded!")

	// Extract real test data from database file
	fmt.Println("Extracting test data from database file...")
	symbols, vegaIds, err := extractTestData("data/hgnc_complete_set.txt.gz")
	if err != nil {
		log.Fatalf("Failed to extract test data: %v", err)
	}
	fmt.Printf("Extracted %d symbol-vega_id pairs for testing\n", len(symbols))
	fmt.Println()

	// Show cache vs no-cache performance
	fmt.Println("=== Cache vs No-Cache Performance Comparison ===")
	compareCachePerformance(hgnc, symbols, vegaIds)
}

func compareCachePerformance(hgnc *h.HGNC, symbols, vegaIds []string) {
	// Limit test data size for reasonable performance
	maxTestSize := 500
	if len(symbols) > maxTestSize {
		symbols = symbols[:maxTestSize]
		vegaIds = vegaIds[:maxTestSize]
	}

	fmt.Printf("Using %d test records\n\n", len(symbols))

	// Reduced operations since non-cached is slow
	testConfigs := []struct {
		name       string
		operations int
		goroutines int
	}{
		{"Small scale", 1000, 1},
		{"Medium scale", 5000, 2},
		{"Large scale", 10000, 4},
	}

	for _, config := range testConfigs {
		fmt.Printf("%s (%d operations, %d goroutines)\n",
			config.name, config.operations, config.goroutines)
		fmt.Printf("%-25s %-12s %-15s %-12s\n",
			"Field", "Time", "Ops/sec", "Cache Status")
		fmt.Printf("%-25s %-12s %-15s %-12s\n",
			strings.Repeat("-", 25), strings.Repeat("-", 12),
			strings.Repeat("-", 15), strings.Repeat("-", 12))

		// Test CACHED field: SYMBOL (has index)
		symbolTime := runFieldTest(hgnc, symbols, h.FIELD_SYMBOL, config.operations, config.goroutines)
		symbolOps := float64(config.operations) / symbolTime.Seconds()
		fmt.Printf("%-25s %-12s %-15.0f %-12s\n",
			"SYMBOL", formatDuration(symbolTime), symbolOps, "✅ cached")

		// Test NON-CACHED field: VEGA_ID (no index)
		vegaTime := runFieldTest(hgnc, vegaIds, h.FIELD_VEGA_ID, config.operations, config.goroutines)
		vegaOps := float64(config.operations) / vegaTime.Seconds()
		fmt.Printf("%-25s %-12s %-15.0f %-12s\n",
			"VEGA_ID", formatDuration(vegaTime), vegaOps, "❌ no cache")

		// Show performance difference (avoid inf when symbolTime is very small)
		var speedupText string
		if symbolTime.Nanoseconds() < 1000 { // Less than 1 microsecond
			speedupText = ">10000x faster"
		} else {
			speedup := vegaTime.Seconds() / symbolTime.Seconds()
			if math.IsInf(speedup, 1) {
				speedupText = ">10000x faster"
			} else {
				speedupText = fmt.Sprintf("%.1fx faster", speedup)
			}
		}

		fmt.Printf("\n💡 Performance Summary:\n")
		fmt.Printf("   SYMBOL (cached):        %s\n", formatDuration(symbolTime))
		fmt.Printf("   VEGA_ID (no cache):     %s\n", formatDuration(vegaTime))
		fmt.Printf("   Cache makes queries %s!\n\n", speedupText)
	}

	// Show field usage recommendation using the GetAllIndexedFieldNames function
	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Printf("Field Usage Recommendation:\n")
	fmt.Printf("• SYMBOL field: ✅ Has pre-built hash map index - O(1) lookup\n")
	fmt.Printf("• VEGA_ID field: ❌ Linear scan through all records - O(n) lookup\n")
	fmt.Printf("• VEGA_ID is rarely used, so no index was built for it\n")
	fmt.Printf("• With ~45,000 gene records, cache provides significant speedup\n")
	fmt.Printf("\n📋 All cached fields (with indexes):\n")

	// Use the package function to get indexed field names
	indexedFieldNames := h.GetAllIndexedFieldNames()
	for _, fieldName := range indexedFieldNames {
		fmt.Printf("   • %s\n", fieldName)
	}

	fmt.Printf("\n💡 Always prefer cached fields when possible!\n")
}

// Run test for specific field type
func runFieldTest(hgnc *h.HGNC, testValues []string, field h.Field, totalOps, numGoroutines int) time.Duration {
	var wg sync.WaitGroup
	opsPerGoroutine := totalOps / numGoroutines

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < opsPerGoroutine; j++ {
				// Use real values from database to query the corresponding field
				value := testValues[(goroutineID*opsPerGoroutine+j)%len(testValues)]
				hgnc.Fetch(value, field)
			}
		}(i)
	}

	wg.Wait()
	return time.Since(start)
}

// Format duration nicely
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

// Extract symbols and vega_ids from the database file
func extractTestData(filepath string) ([]string, []string, error) {
	// Open file
	fh, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer fh.Close()

	gz, err := gzip.NewReader(fh)
	if err != nil {
		return nil, nil, err
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)

	// Read header to find column indexes
	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("failed reading header line")
	}

	headerLine := scanner.Text()
	headers := strings.Split(headerLine, "\t")

	var symbolIdx, vegaIdx int = -1, -1
	for i, header := range headers {
		if header == "symbol" {
			symbolIdx = i
		} else if header == "vega_id" {
			vegaIdx = i
		}
	}

	if symbolIdx == -1 || vegaIdx == -1 {
		return nil, nil, fmt.Errorf("could not find symbol or vega_id columns")
	}

	fmt.Printf("Found columns - symbol: %d, vega_id: %d\n", symbolIdx+1, vegaIdx+1)

	var symbols, vegaIds []string

	// Read data lines
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		if len(fields) <= symbolIdx || len(fields) <= vegaIdx {
			continue
		}

		symbol := strings.TrimSpace(strings.Trim(fields[symbolIdx], "\""))
		vegaId := strings.TrimSpace(strings.Trim(fields[vegaIdx], "\""))

		// Skip if vega_id is empty
		if symbol != "" && vegaId != "" {
			symbols = append(symbols, symbol)
			vegaIds = append(vegaIds, vegaId)
		}
	}

	return symbols, vegaIds, scanner.Err()
}
