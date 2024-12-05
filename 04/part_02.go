package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
	iterPath = make(map[int]*Path)
	hashPath = make(map[string]*Path)
	matrix   [][]byte
	rows     int
	xCount   int
)

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}

	// Use scanner to scan line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bb := scanner.Bytes()

		matrix = append(matrix, bb)
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

	drawColouredInput(matrix, hashPath)

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
			// Skip the very middle element ('A').
			if i == 0 && j == 0 {
				continue
			}

			//
			//  Skip horizontal and vertical coords relative to the X shape center ('A').
			//
			//  i/j | -1  0  1
			//  ----|---------
			//    -1|  .  x  .
			//     0|  x  .  x
			//     1|  .  x  .
			//
			if (i == 0 && j == -1) || (i == -1 && j == 0) || (i == 0 && j == 1) || (i == 1 && j == 0) {
				continue
			}

			// Recalculate the position with relative coords.
			cX := x + i
			cY := y + j

			// Only run within the matrix' bounds.
			if cX >= 0 && cY >= 0 && cX < rows && cY < cols {
				path, found := iterPath[iID]

				var cornerLetter string

				switch string(matrix[cX][cY]) {
				case "A":
					if found {
						continue
					}

					// Prepare the contents for a new Path struct instance.
					pts := [][2]int{{cX, cY}}
					iterPath[iID] = &Path{Coords: pts, Literae: []string{"A"}}

					grid(cX, cY, iID)
					continue

				case "M":
					if !found {
						continue
					}

					cornerLetter = "M"

				case "S":
					if !found {
						continue
					}

					cornerLetter = "S"

				default:
					continue
				}

				point := [2]int{cX, cY}
				newCoords := append(path.Coords, point)
				newPath := &Path{Coords: newCoords, Literae: append(path.Literae, cornerLetter)}

				if !checkXShape(newPath) {
					continue
				}

				if len(newPath.Coords) != 5 {
					iterPath[iID] = newPath
					continue
				}

				if saved := saveToHashPath(newPath); saved {
					xCount++
				}

				newPath.FullShape = true
				iterPath[iID] = newPath
				return
			}
		}
	}
}

func saveToHashPath(path *Path) bool {
	var joined string

	for _, pair := range path.Coords {
		joined += fmt.Sprintf("%d,%d,", pair[0], pair[1])
	}

	// Check if such "hash" has been already saved (do not overwrite).
	if _, found := hashPath[joined]; found {
		return false
	}

	hashPath[joined] = path
	return true
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
		debugf("The remote letter relative diagonally is unexpected (not S when input is M and vice versa))\n")
		return false
	}

	diagonal = [][2]int{pth[len(pth)-1], pth[0], {x, y}}

	debugf("--- The diagonal is ready: %v\n", diagonal)
	return true
}

// drawColouredInput uses the hashPath map <hp> to create a coloured string matrix M based on the initial input matrix <mx>.
func drawColouredInput(mx [][]byte, hp map[string]*Path) {
	var M [][]string = make([][]string, len(mx))

	// At first, copy/convert the input matrix to a string matrix.
	for i, ival := range mx {
		row := make([]string, len(ival))

		for j, jval := range ival {
			row[j] = string(jval)
		}

		M[i] = row
	}

	// Then loop over ther hashPath map to colour the string matrix' elements.
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

	// And finally: print the matrix to stdout.
	for i := 0; i < len(M); i++ {
		for j := 0; j < len(M[0]); j++ {
			fmt.Printf("%s ", M[i][j])
		}
		fmt.Printf("\n")
	}
}
