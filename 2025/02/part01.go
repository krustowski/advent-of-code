package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	line, err := os.ReadFile("input02.txt")
	if err != nil {
		panic("no input file")
	}

	var invalid int

	rangesLF := strings.Split(string(line), "\n")
	ranges := strings.Split(rangesLF[0], ",")

	for _, rang := range ranges {
		pair := strings.Split(rang, "-")

		lower, err := strconv.Atoi(pair[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		upper, err := strconv.Atoi(pair[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		invalid += checkIdRanges(lower, upper)
	}

	fmt.Printf("invalid IDs: %d\n", invalid)
}

func checkIdRanges(lower, upper int) int {
	count := 0

	for i := lower; i <= upper; i++ {
		id := strconv.Itoa(i)

		if len(id)%2 != 0 {
			continue
		}

		s1 := id[0:(len(id) / 2)]
		s2 := id[(len(id) / 2):]

		if s1 == s2 {
			count += i
		}
	}

	return count
}
