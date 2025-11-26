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
	inputFile string
)

func init() {
	flag.StringVar(&inputFile, "inputFile", "", "a string, a valid path to the input data file")

	flag.Parse()

	if inputFile == "" {
		fmt.Println("The inputFile flag cannot be blank")
		os.Exit(1)
	}
}

func runLevelsDiffAnalysis(diffs []int) (result reportResult, stopIdx, problemCount int) {
	fmt.Printf("--- new diffs analysis: %v\n", diffs)

	var (
		decCount int
		incCount int
		//interimResult reportResult
	)

	for j := 0; j < len(diffs); j++ {
		// Zero or (absolutely) too high diff.
		if diffs[j] == 0 || diffs[j] > 3 || diffs[j] < -3 {
			fmt.Printf("problem at index: %d\n", j)

			if problemCount == 0 {
				problemCount++
				stopIdx = j
			}
			continue
		}

		if diffs[j] < 0 && result == INCREASING {
			fmt.Printf("sign change at index: %d\n", j)
			decCount++
			result = DECREASING

			if problemCount == 0 {
				problemCount++
				stopIdx = j - 1
			}
			continue
		}

		if diffs[j] > 0 && result == DECREASING {
			fmt.Printf("sign change at index: %d\n", j)
			incCount++
			result = INCREASING

			if problemCount == 0 {
				problemCount++
				stopIdx = j - 1
			}
			continue
		}

		if diffs[j] < 0 && (result == UNKNOWN || result == DECREASING) {
			decCount++
			result = DECREASING
			continue
		}

		if diffs[j] > 0 && (result == UNKNOWN || result == INCREASING) {
			incCount++
			result = INCREASING
			continue
		}
	}

	fmt.Printf("inc: %d, dec: %d\n", incCount, decCount)

	result = func() reportResult {
		if problemCount > 1 {
			return MIXED
		}

		if (incCount > 1 && decCount > 1) || incCount == decCount {
			//problemCount++
			return MIXED
		}

		if incCount > 0 {
			return INCREASING
		}

		if decCount > 0 {
			return DECREASING
		}

		return MIXED
	}()

	return
}

func main() {
	var (
		// Part2 condition (set to false to get results for the part1 only).
		problemDampenerEnabled bool = true

		reportsSafe   int
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

		var (
			diffs  []int
			report []int
		)

		fmt.Printf("\n--- new report: %v\n", levels)

		// Convert string elements into integers.
		for i := 0; i < len(levels); i++ {
			if level, err := strconv.Atoi(levels[i]); err == nil {
				report = append(report, level)
			} else {
				fmt.Printf("strconv: %s\n", err.Error())
				os.Exit(1)
			}

			if i > 0 {
				diffs = append(diffs, report[i]-report[i-1])
			}
		}

		var (
			idx    int
			idxCpy int
			result reportResult
			tries  int
		)

		// Run multiple analysis (if needed as the Problem Dampener is enabled).
		for {
			tries++
			if (problemDampenerEnabled && tries > 4) || (!problemDampenerEnabled && tries > 1) {
				result = MIXED
				break
			}

			diffsCpy := func() []int {
				if result == UNKNOWN {
					return diffs
				}

				var (
					newDiffs  []int
					newReport []int
				)

				for k, val := range report {
					if k == idxCpy+3-(tries-1) {
						fmt.Printf("removing level: %d\n", val)
						continue
					}

					newReport = append(newReport, val)
				}

				for l := 1; l < len(newReport); l++ {
					newDiffs = append(newDiffs, newReport[l]-newReport[l-1])
				}

				return newDiffs
			}()

			var problemCount int

			if result, idx, problemCount = runLevelsDiffAnalysis(diffsCpy); problemCount == 0 {
				break
			}

			/*if problemCount > 1 {
				result = MIXED
				break
			}*/

			if tries == 1 {
				idxCpy = idx
			}
		}

		if result == INCREASING || result == DECREASING {
			fmt.Print("marked as safe\n")
			reportsSafe++
		} else {
			reportsUnsafe++
		}
	}

	fmt.Printf("\n--- results\n")
	fmt.Printf("safe: \t\t%d\n", reportsSafe)
	fmt.Printf("unsafe: \t%d\n", reportsUnsafe)
}
