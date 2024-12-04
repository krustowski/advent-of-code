package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

var (
	cols      int
	iterID    int
	iterPath  = make(map[int][][2]int)
	hashPath  = make(map[string][][2]int)
	rows      int
	smatrix   [][]rune = make([][]rune, 140)
	xmasCount int
)

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	var (
		counter int
	)

	// Use scanner to scan line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := string(scanner.Bytes())
		smatrix[counter] = []rune(str)

		counter++
	}

	// Check for any possible error.
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner: %s\n", err.Error())
		return
	}

	rows = len(smatrix)
	cols = len(smatrix[0])

	// Should be 140x140, issa hardcoded above but yolo.
	fmt.Printf("Matrix data: %dx%d\n", rows, cols)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Print("\033[33m", string(smatrix[i][j]), "\033[0m ")
		}
		fmt.Printf("\n")
	}

	var (
		prog int
		ctrl int
	)

	for {
		ctrl++

		x := prog / rows
		y := prog % cols

		grid(x, y, iterID)

		if ctrl == rows*cols {
			fmt.Printf("The ctrl control reached (x: %3d, y: %3d)!\n", x, y)
			break
		}

		prog++
	}

	drawColorSolution(smatrix, hashPath)

	fmt.Printf("\n--- stats:\n")
	fmt.Printf("iterations       : %d\n", iterID)
	fmt.Printf("XMAS found (lmao): %d\n", xmasCount)
	fmt.Printf("full path count  : %d\n", len(hashPath))
}

func grid(x, y int, iID int) {
	iterID++

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Prepare relative coords.
			cX := x + i
			cY := y + j

			if i == 0 && j == 0 {
				continue
			}

			if (cX >= 0 && cY >= 0) && (cX < rows && cY < cols) {
				L := smatrix[cX][cY]

				// Fetch the path from iID.
				path, found := iterPath[iID]

				// How many items are already in.
				l := len(path)

				switch string(L) {
				case "X":
					point := [2]int{cX, cY}

					if found {
						continue
					}

					iterPath[iterID] = append(iterPath[iterID], point)

					grid(cX, cY, iterID)

				case "M":
					if !found || (found && l != 1) {
						continue
					}

					point := [2]int{cX, cY}
					iterPath[iterID] = append(path, point)

					grid(cX, cY, iterID)

				case "A":
					if !found || (found && l != 2) {
						continue
					}

					point := [2]int{cX, cY}
					newPath := append(path, point)

					if !isDirectional(newPath) {
						continue
					}

					iterPath[iterID] = newPath

					grid(cX, cY, iterID)

				case "S":
					if !found || (found && l != 3) {
						continue
					}

					point := [2]int{cX, cY}
					fullPath := append(path, point)

					if !isDirectional(fullPath) {
						continue
					}

					//fmt.Printf("Full path: %v\n", fullPath)
					saveToHashPath(&fullPath, iterID)

					xmasCount++
				}
			}
		}
	}
}

func saveToHashPath(pts *[][2]int, iID int) {
	var joined string

	for _, pair := range *pts {
		joined += fmt.Sprintf("%d,%d,", pair[0], pair[1])
	}

	hashPath[joined] = *pts
}

func isDirectional(pts [][2]int) bool {
	var dX, dY int

	for i := range pts {
		switch i {
		case 0:
			continue

		case 1:
			dX = pts[i][0] - pts[i-1][0]
			dY = pts[i][1] - pts[i-1][1]

		default:
			dx := pts[i][0] - pts[i-1][0]
			dy := pts[i][1] - pts[i-1][1]

			if dx != dX || dy != dY {
				return false
			}
		}
	}

	return true
}

func drawColorSolution(mx [][]rune, hp map[string][][2]int) {
	var M [][]string = make([][]string, len(mx))

	for i, ival := range mx {
		row := make([]string, len(ival))

		for j, jval := range ival {
			row[j] = string(jval)
		}

		M[i] = row
	}

	for _, path := range hp {
		for _, pair := range path {
			x := pair[0]
			y := pair[1]

			OG := M[x][y]
			M[x][y] = "\033[32m" + OG + "\033[0m"
		}
	}

	for i := 0; i < len(M); i++ {
		for j := 0; j < len(M[0]); j++ {
			fmt.Printf("%s ", M[i][j])
		}
		fmt.Printf("\n")
	}
}
