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

			// fmt.Printf("buf len: %d\n", len(buf))
			continue
		}
	}

	return true, matches
}

func inspectSliceOld(str string) (bool, int) {
	buf := make([]rune, 0, len(str))
	matches := 0

	for i, r := range str {
		if i == 0 {
			buf = append(buf, r)
			continue
		}

		if i == len(str)-1 {
			if subInspectSlice(string(str[len(str)-1]), len(buf)-1, buf) {
				matches++
				break
			}
		}

		if subInspectSlice(str[i:], i-1, buf) {
			matches++
		} else {
			buf = append(buf, r)
		}
	}

	if matches > 0 && len(str)%(matches) == 0 {
		fmt.Printf("%v\n", string(buf))
		return true, matches
	}

	return false, 0
}

func subInspectSlice(str string, offset int, trend []rune) bool {
	if len(str) == 0 {
		return false
	}

	if len(str) == 1 && offset < len(trend) {
		if string(str[0]) == string(trend[offset]) {
			return true
		}
	}

	if len(str)%2 == 0 {
		if len(trend) == 2 {
			if string(trend) == string(str) {
				return true
			}
		}

		if subInspectSlice(str[1:], offset, trend) {
			return true
		}
	}

	if len(str) > 2 {
		if len(str)%len(trend) == 0 {
			if subInspectSlice(str[3:], offset, trend) {
				return true
			}
		}

		if subInspectSlice(str[1:], offset+1, trend) {
			return true
		}
	}

	return false
}
