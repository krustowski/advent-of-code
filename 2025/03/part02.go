package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fhdl, err := os.Open("input03t.txt")
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
		og := scanner.Text()

		if len(og) < 12 {
			continue
		}

		line := strings.Clone(og)

		//

		poss := make([]int, 0, 12)

		for i := 0; len(poss) < 12; i++ {
			if len(poss) > 0 {
				if i > len(poss) {
					i = 1
				}

				for j := 0; j < len(line)-poss[i-1]+1; j++ {
					if len(poss) == 12 {
						break
					}

					pos, err := getNextMaxPosition(&line, (poss[i-1]+j)%(len(line)-1))
					if err != nil {
						fmt.Println(err)
						continue
					}

					poss = append(poss, pos)
				}
			} else {
				pos, err := getNextMaxPosition(&line, 0)
				if err != nil {
					fmt.Println(err)
					continue
				}

				poss = append(poss, pos)
			}
		}

		//

		fmt.Printf("%v \n", poss)

		totalJuice += assembleBankJoltage(poss, og)
	}

	fmt.Printf("joltage: %d\n", totalJuice)
}

func getNextMaxPosition(line *string, offset int) (int, error) {
	if len(*line) <= offset {
		return 0, errors.New("offset out of bounds")
	}

	pos, err := extractPosFromString((*line)[offset:])
	if err != nil {
		return 0, err
	}

	pos += offset
	*line = shadeBatteryInBank(*line, pos)

	return pos, nil
}

func assembleBankJoltage(poss []int, og string) int {
	slices.Sort(poss)

	var total int

	for i, p := range poss {
		num, err := strconv.Atoi(string([]rune(og)[p]))
		if err != nil {
			fmt.Println(err)
			continue
		}

		total += int(float64(num) * math.Pow(10, float64(11-i)))
	}

	fmt.Printf("%d\n", total)
	return total
}

func extractPosFromString(str string) (int, error) {
	var max int
	var pos int

	for i, r := range str {
		num, err := strconv.Atoi(string(r))
		if err != nil {
			return 0, err
		}

		if num > max {
			max = num
			pos = i
		}
	}

	if max > 0 {
		return pos, nil
	}

	return 0, errors.New("attempt to read a zero value")
}

func shadeBatteryInBank(bank string, pos int) string {
	runes := []rune(bank)
	runes[pos] = '0'
	return string(runes)
}
