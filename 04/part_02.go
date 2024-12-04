package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//"time"
)

//
//  Warning: this code does not behave properly (yet). Therefore will not most prolly work with your input.
//

var (
	debug     bool
	inputFile string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "a boolean, makes the output more verbose")
	flag.StringVar(&inputFile, "inputFile", "", "a string, a valid path to the input data file")

	flag.Parse()

	if inputFile == "" {
		fmt.Println("The inputFile flag cannot be blank")
		os.Exit(1)
	}
}

func debugf(format string, args ...interface{}) {
	if format == "" || !debug {
		return
	}

	fmt.Printf(format, args...)
}

type Path struct {
	Coords    [][2]int
	Literae   []string
	FullShape bool
}

var (
	cols     int
	iterID   int
	iterPath          = make(map[int]*Path)
	hashPath          = make(map[string]*Path)
	matrix   [][]byte = make([][]byte, 140)
	rows     int
	xCount   int
)

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	var counter int

	// Use scanner to scan line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bb := make([]byte, 140)
		bb = scanner.Bytes()

		matrix[counter] = bb
		counter++
	}

	// Check for any possible error.
	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner: %s\n", err.Error())
		return
	}

	rows = len(matrix)
	cols = len(matrix[0])

	debugf("Matrix dims: rows: %d, cols: %d\n", rows, cols)

	var (
		ctrl int
		prog int
	)

	for {
		ctrl++

		x := prog / rows
		y := prog % cols

		// Run the recursion
		grid(x, y, iterID)

		if ctrl == rows*cols {
			debugf("The ctrl control reached (x: %3d, y: %3d)!\n", x, y)
			break
		}

		prog++
	}

	drawColorSolution(matrix, hashPath)

	fmt.Printf("\n--- stats:\n")
	fmt.Printf("iterations : %d\n", iterID)
	fmt.Printf("hashPath   : %d\n", len(hashPath))
	fmt.Printf("X count    : %d\n", xCount)
	fmt.Printf("fullShape  : %d\n", func() int {
		var ctr int

		for _, path := range iterPath {
			if !path.FullShape {
				continue
			}

			ctr++
		}

		return ctr
	}())
}

func grid(x, y int, iID int) {
	iterID++

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			cX := x + i
			cY := y + j

			// The very middle element = 'A'
			if i == 0 && j == 0 {
				continue
			}

			// Skip horizontal and vertical coords.
			if (i == 0 && j == -1) || (i == -1 && j == 0) || (i == 0 && j == 1) || (i == 1 && j == 0) {
				continue
			}

			if (cX >= 0 && cY >= 0) && (cX < rows && cY < cols) {
				path, found := iterPath[iID]

				switch string(matrix[cX][cY]) {
				case "A":
					if found {
						continue
					}

					//_ = int(time.Now().Unix())

					pts := [][2]int{{cX, cY}}
					iterPath[iterID] = &Path{Coords: pts, Literae: []string{"A"}}
					//iterPath[id] = &Path{Coords: pts, Literae: []string{"A"}}

					grid(cX, cY, iterID)
					//grid(cX, cY, id)

				case "M":
					if !found {
						continue
					}

					point := [2]int{cX, cY}
					newCoords := append(path.Coords, point)
					newPath := &Path{Coords: newCoords, Literae: append(path.Literae, "M")}

					if !checkXShape(newPath) {
						continue
					}

					if len(newPath.Coords) == 5 {
						newPath.FullShape = true
						xCount++
						iterPath[iID] = newPath

						//fmt.Printf("FullShape: Literae: %v\n", newPath.Literae)
						//fmt.Printf("FullShape: Coords: %v\n", newPath.Coords)

						saveToHashPath(newPath)

						return
					}

					iterPath[iID] = newPath

				case "S":
					if !found {
						continue
					}

					point := [2]int{cX, cY}
					newCoords := append(path.Coords, point)
					newPath := &Path{Coords: newCoords, Literae: append(path.Literae, "S")}

					if !checkXShape(newPath) {
						continue
					}

					if len(newPath.Coords) == 5 {
						newPath.FullShape = true
						xCount++
						iterPath[iID] = newPath

						//fmt.Printf("FullShape: Literae: %v\n", newPath.Literae)
						//fmt.Printf("FullShape: Coords: %v\n", newPath.Coords)

						saveToHashPath(newPath)

						return
					}

					iterPath[iID] = newPath
				}
			}
		}
	}
}

func saveToHashPath(path *Path) {
	var joined string

	for _, pair := range path.Coords {
		joined += fmt.Sprintf("%d,%d,", pair[0], pair[1])
	}

	hashPath[joined] = path
}

// checkXShape is called exclusively in cases when the letters 'M' or 'S' are hit to be around the core 'A' letter.
// The purpose of this function is to filter out wrong letter constellations: reject the 'plus' shape, embrace the X shape.
func checkXShape(path *Path) bool {
	var (
		diagonal   [][2]int
		passiveLtr string
	)

	pth := path.Coords
	ltr := path.Literae

	// Check the 'plus' constellation. The first letter in the Literae slice must be 'A'.
	if ltr[0] != "A" {
		debugf("The first letter/char of the Literae slice is not 'A': invalid path\n")
		return false
	}

	//
	//  The X-shape:
	//
	//  M . S . . S . S
	//  . A . . . , A ,
	//  M . S . . M . M
	//

	switch ltr[len(ltr)-1] {
	case "M":
		passiveLtr = "S"

	case "S":
		passiveLtr = "M"
	}

	// Check the slope of such constellation.
	dx := pth[len(pth)-1][0] - pth[0][0]
	dy := pth[len(pth)-1][1] - pth[0][1]

	// Zero indicates the slope is horizontal/vertical = the 'plus' shape to be omitted.
	if dx == 0 || dy == 0 {
		debugf("The 'plus' shape, bye\n")
		return false
	}

	// Coords for the remote letter diagonally.
	x := pth[0][0] - dx
	y := pth[0][1] - dy

	// Check bounds.
	if x < 0 || y < 0 || x >= rows || y >= cols {
		debugf("Out of bounds\n")
		return false
	}

	if string(matrix[x][y]) != passiveLtr {
		debugf("The letter located remote diagonally is unexpected\n")
		return false
	}

	diagonal = [][2]int{pth[len(pth)-1], pth[0], {x, y}}

	debugf("--- Diagonal ready: %v\n", diagonal)
	return true
}

func drawColorSolution(mx [][]byte, hp map[string]*Path) {
	var M [][]string = make([][]string, len(mx))

	for i, ival := range mx {
		row := make([]string, len(ival))

		for j, jval := range ival {
			row[j] = string(jval)
		}

		M[i] = row
	}

	for _, path := range hp {
		if !path.FullShape {
			continue
		}

		for _, pair := range path.Coords {

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
