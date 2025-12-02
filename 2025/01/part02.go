package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const dialSize = 100

var debug = false

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

				for i := 1; i < dist; i++ {
					if (dialSize+(position-i))%dialSize == 0 {
						counter++
						debugf("=> new zero!\n")
					}
				}
			}
		case "R":
			{
				op = 1

				for i := 1; i < dist; i++ {
					if (position+i)%dialSize == 0 {
						counter++
						debugf("=> new zero!\n")
					}
				}
			}
		}

		position += op * dist

		if position < 0 {
			position = (dialSize + position) % dialSize
		} else {
			position %= dialSize
		}

		if position == 0 {
			counter++
			debugf("=> new zero!\n")
		}

		debugf("%s to position: %d\n", string(bb), position)
	}

	fmt.Printf("counter: %d\n", counter)
}

func debugf(f string, arg ...any) {
	if debug {
		fmt.Printf(f, arg...)
	}
}
