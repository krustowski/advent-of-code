package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	fhdl, err := os.Open("input03.txt")
	if err != nil {
		panic("no input file")
	}
	defer func() {
		err := fhdl.Close()
		if err != nil {
			panic("cannot close the file handler")
		}
	}()

	var totalJuice int

	scanner := bufio.NewScanner(fhdl)
	for scanner.Scan() {
		line := scanner.Text()

		var max1 int
		var pos1 int

		for i, r := range line {
			num, err := strconv.Atoi(string(r))
			if err != nil {
				fmt.Println(err)
				continue
			}

			if num > max1 {
				max1 = num
				pos1 = i
			}
		}

		var max2 int
		var pos2 int

		runes := []rune(line)
		runes[pos1] = '0'
		line = string(runes)

		for i, r := range line {
			num, err := strconv.Atoi(string(r))
			if err != nil {
				fmt.Println(err)
				continue
			}

			if num > max2 {
				max2 = num
				pos2 = i
			}
		}

		if pos1 < len(line) {
			runes := []rune(line)
			runes[pos1] = '0'
			line = string(runes[pos1:])
		}

		var max3 int
		// var pos3 int

		if len(line) > 1 {
			for _, r := range line {
				num, err := strconv.Atoi(string(r))
				if err != nil {
					fmt.Println(err)
					continue
				}

				if num > max3 {
					max3 = num
					// pos3 = i
				}
			}
		}

		var joltage []int

		if pos1 < pos2 {
			joltage = append(joltage, max1*10+max2)
		} else {
			joltage = append(joltage, max2*10+max1)
		}
		if max3 > 0 {
			joltage = append(joltage, max1*10+max3)
		} else {
			joltage = append(joltage, 0)
		}

		totalJuice += int(math.Max(float64(joltage[0]), float64(joltage[1])))

		// fmt.Printf("%v \n", joltage)
		// fmt.Printf("joltage: %d\n", int(math.Max(float64(joltage[0]), float64(joltage[1]))))
	}

	fmt.Printf("total joltage: %d\n", totalJuice)
}
