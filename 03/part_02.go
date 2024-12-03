package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var inputFile string

func init() {
	flag.StringVar(&inputFile, "inputFile", "", "a string, a valid path to the input data file")

	flag.Parse()

	if inputFile == "" {
		fmt.Println("The inputFile flag cannot be blank")
		os.Exit(1)
	}
}

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)

	var instructions []string

	for {
		rawLine, err := reader.ReadString(byte('\n'))
		if err != nil {
			break
		}

		regex, err := regexp.Compile("(mul\\((.{1,3},.{1,3}\\))|do\\((\\)))|don't\\((\\))")
		if err != nil {
			panic(err)
		}

		match := regex.FindAllString(rawLine, -1)

		instructions = append(instructions, match...)
	}

	fmt.Printf("instructions: %v\n", instructions)

	var (
		counter    int
		mulEnabled bool = true
		mulSum     int64
	)

	for _, instruction := range instructions {
		switch instruction {
		case "do()":
			mulEnabled = true
			continue
		case "don't()":
			mulEnabled = false
			continue
		default:
			if !mulEnabled {
				continue
			}
		}

		regex, err := regexp.Compile("\\d+,\\d+")
		if err != nil {
			fmt.Printf("Regexp compilation error: %s\n", err.Error())
			os.Exit(2)
		}

		joinedElems := regex.FindAllString(instruction, -1)

		elems := strings.Split(joinedElems[0], ",")

		if len(elems) != 2 {
			fmt.Printf("Unexpected len of elements in the extracted array: %d; %v\n", len(elems), elems)
			continue
		}

		val1, err := strconv.ParseInt(elems[0], 10, 0)
		if err != nil {
			fmt.Printf("String to int conversion error: %s\n", err.Error())
			os.Exit(3)
		}

		val2, err := strconv.ParseInt(elems[1], 10, 0)
		if err != nil {
			fmt.Printf("String to int conversion error: %s\n", err.Error())
			os.Exit(4)
		}

		counter++
		mulSum += val1 * val2
	}

	fmt.Printf("\n--- results\n")
	fmt.Printf("Cummulative `mul` instruction result: %d (mul count: %d)\n", mulSum, counter)
}
