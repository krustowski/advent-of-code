package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
		leftVals  []int
		rightVals []int
	)

	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(string(scanner.Bytes()), "   ")

		if debug {
			fmt.Printf("%v, %d\n", parts, len(parts))
		}

		if left, err := strconv.Atoi(parts[0]); err == nil {
			leftVals = append(leftVals, left)
		} else {
			fmt.Printf("strconv: %s\n", err.Error())
		}

		if right, err := strconv.Atoi(parts[1]); err == nil {
			rightVals = append(rightVals, right)
		} else {
			fmt.Printf("strconv: %s\n", err.Error())
		}
	}

	// Check for any possible error.
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner: %s\n", err.Error())
		return
	}

	slices.Sort(leftVals)
	slices.Sort(rightVals)

	if len(leftVals) != len(rightVals) {
		fmt.Println("Value slices are not of the same len")
		os.Exit(2)
	}

	var totalDist int

	// Iterate over sorted list(s) and calculate an atomic distance between the same-indexed elements.
	for i := 0; i < len(leftVals); i++ {
		var atomicDist int

		if atomicDist = leftVals[i] - rightVals[i]; atomicDist < 0 {
			atomicDist *= -1
		}

		if debug {
			fmt.Printf("%d\n", atomicDist)
		}

		totalDist += atomicDist
	}

	fmt.Printf("Total distance: %d\n", totalDist)

	os.Exit(0)
}
