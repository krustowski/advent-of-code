package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	bytes, err := os.ReadFile("input02t.txt")
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
			// count += i
			nums[i] = make([]int, 2)
			nums[i][0] = num
			nums[i][1] = matches

			if matches > max {
				max = matches
			}

			i++
		}
	}

	for _, val := range nums {
		if len(val) == 2 && val[1] == max {
			count = val[0]

			fmt.Printf("%d: \t %d\n", val[0], val[1])
		}
	}

	return count
}

func inspectSlice(str string) (bool, int) {
	buf := make([]rune, 0, len(str))
	matches := 0

	for j, r := range str {
		if j == 0 {
			buf = append(buf, r)
			continue
		}

		if j+len(buf) <= len(str) && j-len(buf) >= 0 {
			if str[j-1] == str[j] {
				buf = append(buf, r)
				matches++
				continue
			}

			if str[j:j+len(buf)] == str[j-len(buf):j] {
				buf = append(buf, r)
				matches++
				continue
			}

			if j+len(buf) <= len(str) && string(buf) == str[j:j+len(buf)] {
				buf = append(buf, r)
				matches++
				continue
			}

			buf = append(buf, r)
		}
	}

	if matches > 0 && len(str)%(matches) == 0 {
		return true, matches
	}

	return false, 0
}
