package Netpbm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	var magicNumber string
	var width int
	var height int
	var counter int
	var headersize int
	var splitfile []string
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(string(file), "\r") {
		splitfile = strings.SplitN(string(file), "\r\n", -1)
	} else {
		splitfile = strings.SplitN(string(file), "\n", -1)
	}

	// Range each line, to insert information into the struct
	for i, _ := range splitfile {
		if strings.Contains(splitfile[i], "P1") {
			magicNumber = "P1"
		} else if strings.Contains(splitfile[i], "P4") {
			magicNumber = "P4"
		}

		if strings.HasPrefix(splitfile[i], "#") && height != 0 {
			headersize = counter
		}

		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			height, err = strconv.Atoi(splitl[1])
			headersize = counter
		}
		counter++
	}

	data := make([][]bool, height)

	for j := 0; j < height; j++ {
		data[j] = make([]bool, width)
	}
	var splitdata []string

	if counter > headersize {
		// Split the data each time there is a space to insert it into the data
		for i := 0; i < height; i++ {
			data[i] = make([]bool, width)
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			for k := 0; k < width; k++ {
				if splitdata[k] == "1" {
					data[i][k] = true
				}
				if splitdata[k] == "0" {
					data[i][k] = false
				}
			}
		}
	}
	pbm := &PBM{data: data, magicNumber: magicNumber, width: width, height: height}

	return pbm, err
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, pbm.magicNumber)
	if err != nil {
		return err
	}

	fmt.Fprintln(file, pbm.width, pbm.height)

	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			if pbm.data[y][x] {
				fmt.Fprint(file, "1", " ")
			} else {
				fmt.Fprint(file, "0", " ")
			}
		}
		fmt.Fprintln(file)
	}
	return err
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := 0; i < len(pbm.data); i++ {
		for j := 0; j < len(pbm.data[i]); j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	if len(pbm.data[0]) > 0 {
		for i := 0; i < len(pbm.data); i++ {
			for j := 0; j < len(pbm.data[i])/2; j++ {
				startdata := pbm.data[i][j]
				pbm.data[i][j] = pbm.data[i][len(pbm.data[i])-1-j]
				pbm.data[i][len(pbm.data[i])-1-j] = startdata
			}
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	if len(pbm.data) > 0 {
		for i := 0; i < len(pbm.data)/2; i++ {
			startdata := pbm.data[i]
			pbm.data[i] = pbm.data[len(pbm.data)-1-i]
			pbm.data[len(pbm.data)-1-i] = startdata
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
