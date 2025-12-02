package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	bytes, err := os.ReadFile("input02.txt")
	if err != nil {
		panic("no input file")
	}

	var invalid int

	rangesLF := strings.Split(string(bytes), "\n")
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
	var count int
	var max int
	var i int

	nums := make([][]int, upper-lower+1)

	for num := lower; num <= upper; num++ {
		id := strconv.Itoa(num)

		if found, matches := inspectSlice(id); found {
			if matches < 2 {
				continue
			}

			nums[i] = make([]int, 2)
			nums[i][0] = num
			nums[i][1] = matches

			count += num

			if matches > max {
				max = matches
			}

			fmt.Printf("%d: \t %d\n", num, matches)
			i++
		}
	}

	return count
}

func inspectSlice(str string) (bool, int) {
	buf := []rune(str[:len(str)/2])
	matches := 0

	for i := 0; i < len(str)/2; i++ {
		rgx, err := regexp.Compile(string(buf))
		if err != nil {
			fmt.Println(err)
			return false, 0
		}

		res := rgx.FindAllString(str, -1)
		if res == nil {
			return false, 0
		}

		if len(res) > 1 && len(str)%len(res) == 0 {
			if len(buf)*len(res) == len(str) {
				matches = int(math.Max(float64(matches), float64(len(res))))
			}
		}

		if len(buf) > 1 {
			buf = buf[:len(buf)-1]
			continue
		}
	}

	return true, matches
}
