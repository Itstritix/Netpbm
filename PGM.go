package Netpbm

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	var err error
	var magicNumber string = ""
	var width int
	var height int
	var maxval int
	var counter int
	var headersize int
	var intdata int
	var uintval uint8
	var splitfile []string
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(file), "\r") {
		splitfile = strings.SplitN(string(file), "\r\n", -1)
	} else {
		splitfile = strings.SplitN(string(file), "\n", -1)
	}
	// Range each line, to insert information into the struct
	for i, _ := range splitfile {
		if strings.Contains(splitfile[i], "P2") {
			magicNumber = "P2"
		} else if strings.Contains(splitfile[i], "P5") {
			magicNumber = "P5"
		}
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}
		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			height, err = strconv.Atoi(splitl[1])
			headersize = counter
		}
		if maxval == 0 && width != 0 {
			maxval, err = strconv.Atoi(splitfile[i])
			headersize = counter
		}
		counter++

	}

	data := make([][]uint8, height)

	var splitdata []string

	if counter > headersize {
		for i := 0; i < height; i++ {
			// Split the data each time there is a space to insert it into the data
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			data[i] = make([]uint8, width)
			for k := 0; k < width; k++ {
				intdata, err = strconv.Atoi(splitdata[k])
				if err != nil {
					return nil, err
				}
				uintval = uint8(intdata)
				if intdata < maxval {
					data[i][k] = uintval
				} else {
					data[i][k] = uint8(maxval)

				}
			}
		}
	}
	return &PGM{data: data, width: width, height: height, magicNumber: magicNumber, max: uint8(maxval)}, err
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		file.Close()
		return err
	}
	// Write the magic number
	_, err = fmt.Fprintln(file, pgm.magicNumber)
	if err != nil {
		file.Close()
		return err
	}
	// Write the width and height
	_, err = fmt.Fprintln(file, pgm.width, pgm.height)
	if err != nil {
		file.Close()
		return err
	}
	// Write the max value
	_, err = fmt.Fprintln(file, pgm.max)
	if err != nil {
		file.Close()
		return err
	}

	// Insert the data into the file
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			if pgm.data[y][x] >= 0 && pgm.data[y][x] <= pgm.max {
				fmt.Fprint(file, pgm.data[y][x], " ")
			} else {
				return errors.New("wrong data")
			}
		}
		fmt.Fprintln(file)
	}
	return err
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = pgm.max - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	if len(pgm.data[0]) > 0 {
		for i := 0; i < len(pgm.data); i++ {
			for j := 0; j < len(pgm.data[i])/2; j++ {
				startdata := pgm.data[i][j]
				pgm.data[i][j] = pgm.data[i][len(pgm.data[i])-1-j]
				pgm.data[i][len(pgm.data[i])-1-j] = startdata
			}
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	if len(pgm.data) > 0 {
		for i := 0; i < len(pgm.data)/2; i++ {
			startdata := pgm.data[i]
			pgm.data[i] = pgm.data[len(pgm.data)-1-i]
			pgm.data[len(pgm.data)-1-i] = startdata
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	oldMax := pgm.max
	pgm.max = maxValue
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = pgm.data[i][j] * uint8(5) / oldMax
		}
	}
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	newData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.height; i++ {
		newData[i] = make([]uint8, pgm.height)
		for j := 0; j < pgm.width; j++ {
			newData[i][j] = pgm.data[j][i]
		}
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			prevdata := newData[i][j]
			newData[i][j] = newData[i][pgm.height-j-1]
			newData[i][pgm.height-j-1] = prevdata
		}
	}

	pgm.data = newData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	data := make([][]bool, pgm.height)
	pbmnb := ""
	// Convert PGM data to PBM data
	for i := 0; i < pgm.height; i++ {
		data[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			if pgm.data[i][j] < pgm.max/2 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}

	// Change the magic number to the correct one
	if pgm.magicNumber == "P2" {
		pbmnb = "P1"
	} else {
		pbmnb = "P4"
	}

	pbmstruct := &PBM{data: data, magicNumber: pbmnb, width: pgm.width, height: pgm.height}

	return pbmstruct
}
