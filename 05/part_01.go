package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	debug     bool
	inputFile string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "a boolean, makes the output more verbose")
	flag.StringVar(&inputFile, "inputFile", "", "a string, a valid path to the input data file")

	flag.Parse()

	if inputFile == "" {
		fmt.Println("The inputFile flag cannot be blank")
		os.Exit(1)
	}
}

func debugf(format string, args ...interface{}) {
	if format == "" || !debug {
		return
	}

	fmt.Printf(format, args...)
}

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	var (
		pageRules       [][2]int
		pageUpdates     [][]int
		sectionBreakHit bool
	)

	// Use scanner to scan line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		if string(line) == "" {
			sectionBreakHit = true
		}

		switch sectionBreakHit {
		case false:
			rules := strings.Split(string(line), "|")

			var rulePair []int

			for _, ruleStr := range rules {
				rule, err := strconv.Atoi(ruleStr)
				if err != nil {
					fmt.Printf("Scanner: rules: %s\n", err.Error())
					os.Exit(1)
				}

				rulePair = append(rulePair, rule)
			}

			pageRules = append(pageRules, [2]int{rulePair[0], rulePair[1]})

		case true:
			updates := strings.Split(string(line), ",")

			// Skip a blank line (at the end usually).
			if string(line) == "" {
				continue
			}

			var updateList []int

			for _, updateStr := range updates {
				update, err := strconv.Atoi(updateStr)
				if err != nil {
					fmt.Printf("Scanner: updates: %s\n", err.Error())
					os.Exit(2)
				}

				updateList = append(updateList, update)
			}

			pageUpdates = append(pageUpdates, updateList)
		}
	}

	// Check for any possible error.
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner: %s\n", err.Error())
		return
	}

	debugf("Rules: %v\n", pageRules)
	debugf("Updates: %v\n", pageUpdates)

	//
	//
	//

	// Run the batch analysis of page updates matrix according to the rules matrix.
	middlePageSum, err := checkPageUpdateOrder(&pageRules, &pageUpdates)
	if err != nil {
		fmt.Printf("checkPageUpdateOrder: %s\n", err.Error())
		os.Exit(3)
	}

	fmt.Printf("\n--- results:\n")
	fmt.Printf("middlePageSum: %d\n", middlePageSum)
	os.Exit(0)
}

func checkPageUpdateOrder(rules *[][2]int, updates *[][]int) (middlesSum int, err error) {
	for _, update := range *updates {
		if ordered := checkSingleUpdate(rules, update); !ordered {
			continue
		}

		debugf("*** marked as ordered\n")

		// Add the middle page of such page update list to the running sum if ordered right.
		middlesSum += update[(len(update)-1)/2]
	}

	return
}

func checkSingleUpdate(rules *[][2]int, update []int) bool {
	if len(update) < 2 {
		return true
	}

	debugf("\n--- new update analysis\n")
	debugf("update: %v\n", update)

	// Iterate over the list of an update's pages. 
	for i := 0; i < len(update)-1; i++ {
		page := update[i]
		remnant := update[i+1:]

		debugf("page: %v\n", page)
		debugf("remnant: %v\n", remnant)

		// Iterate over the remnant of an update.
		for j := 0; j < len(remnant); j++ {
			for _, rulePair := range *rules {
				// Compare the 'page' and the current updateList item.
				if page == rulePair[0] && remnant[j] == rulePair[1] {
					continue
				}

				// Rule violation check.
				if page == rulePair[1] && remnant[j] == rulePair[0] {
					return false
				}
			}
		}
	}

	return true
}
