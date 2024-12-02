package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type reportResult int

const (
	UNKNOWN reportResult = iota
	DECREASING
	INCREASING
	MIXED
)

var (
	debug     bool
	inputFile string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "a boolean, verbose mode switch")
	flag.StringVar(&inputFile, "inputFile", "", "a string, a valid path to the input data file")

	flag.Parse()

	if inputFile == "" {
		fmt.Println("The inputFile flag cannot be blank")
		os.Exit(1)
	}
}

func main() {
	var (
		reportsSafe int
		reportsUnsafe int
	)

	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)

	// Scan one line at time.
	for scanner.Scan() {
		levels := strings.Split(string(scanner.Bytes()), " ")

		var report []int
		var result reportResult

		fmt.Printf("\n--- new report: %v\n", levels)
	
		// Convert string elements into integers, run tests.
		for i := 0; i < len(levels); i++ {
			if level, err := strconv.Atoi(levels[i]); err == nil {
				report = append(report, level)
			} else {
				fmt.Printf("strconv: %s\n", err.Error())
				os.Exit(1)
			}

			// Check adjacent elements.
			if i > 0 {
				diff := report[i] - report[i-1]

				fmt.Printf("diff: %d\n", diff)

				if diff == 0 || diff > 3 || diff < -3 {
					// Not increasing nor decreasing => unsafe.
					// Diff > 3 || < 3 => unsafe.
					fmt.Printf("Diff too high or zero\n")
					result = MIXED
					reportsUnsafe++
					break
				}

				if (diff > 0 && result == DECREASING) || (diff < 0 && result == INCREASING) {
					fmt.Printf("MIXED condition\n")
					result = MIXED
					reportsUnsafe++
					break
				}

				if diff > 0 && (result == UNKNOWN || result == INCREASING) {
					result = INCREASING
					continue
				}

				if diff < 0 && (result == UNKNOWN || result == DECREASING) {
					result = DECREASING
					continue
				}
			}
		}

		if result == INCREASING || result == DECREASING {
			reportsSafe++
		}
	}

	fmt.Printf("\n--- results\n")
	fmt.Printf("safe: \t\t%d\n", reportsSafe)
	fmt.Printf("unsafe: \t%d\n", reportsUnsafe)
}
