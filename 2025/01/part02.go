package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const dialSize = 100

func main() {
	fhdl, err := os.Open("input01.txt")
	if err != nil {
		panic("no input file")
	}
	defer func() {
		err := fhdl.Close()
		if err != nil {
			panic("cannot close the file handler")
		}
	}()

	var (
		counter  int
		position = 50
	)

	scanner := bufio.NewScanner(fhdl)
	for scanner.Scan() {
		bb := scanner.Bytes()
		op := 1

		distRaw := string(bb[1:])
		dist, err := strconv.Atoi(distRaw)
		if err != nil {
			panic("cannot read the instruction")
		}

		switch string(bb[0]) {
		case "L":
			{
				op = -1
				if position-dist < 0 && position-dist > -dialSize && position != 0 {
					counter++
				}
				if position-dist <= -dialSize {
					counter -= (position - dist) / dialSize
				}
			}
		case "R":
			{
				op = 1
				if position+dist >= dialSize {
					counter += (position + dist) / dialSize
				}
			}
		}

		position += op * dist
		if position < 0 {
			position = (dialSize + position%dialSize) % dialSize
		}
		position %= dialSize

		if position == 0 {
			counter++
		}

		fmt.Printf("position: %d\n", position)
	}

	fmt.Printf("counter: %d\n", counter)
}
