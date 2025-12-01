package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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
		counter  int16
		position = 50
	)

	scanner := bufio.NewScanner(fhdl)
	for scanner.Scan() {
		bb := scanner.Bytes()
		op := 1

		switch string(bb[0]) {
		case "L":
			{
				op = -1
			}
		case "R":
			{
				op = 1
			}
		}

		distRaw := string(bb[1:])
		dist, err := strconv.Atoi(distRaw)
		if err != nil {
			panic("cannot read the instruction")
		}

		position += op * dist
		position %= 100

		if position == 0 {
			counter++
		}
	}

	fmt.Printf("counter: %d\n", counter)
}
