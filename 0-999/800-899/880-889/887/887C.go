package main

import (
	"bufio"
	"fmt"
	"os"
)

func rotate(c []int, cycles [][]int) []int {
	res := make([]int, len(c))
	copy(res, c)
	for _, cyc := range cycles {
		tmp := c[cyc[0]]
		for i := 0; i < len(cyc)-1; i++ {
			res[cyc[i]] = c[cyc[i+1]]
		}
		res[cyc[len(cyc)-1]] = tmp
	}
	return res
}

func solved(c []int) bool {
	for i := 0; i < 24; i += 4 {
		if !(c[i] == c[i+1] && c[i] == c[i+2] && c[i] == c[i+3]) {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	cube := make([]int, 24)
	for i := 0; i < 24; i++ {
		if _, err := fmt.Fscan(in, &cube[i]); err != nil {
			return
		}
	}

	// face indices: U0-3,F4-7,R8-11,B12-15,L16-19,D20-23
	rotations := [][][]int{
		{{0, 1, 3, 2}, {4, 5, 8, 9, 12, 13, 16, 17}},       // U
		{{20, 21, 23, 22}, {6, 7, 18, 19, 14, 15, 10, 11}}, // D
		{{4, 5, 7, 6}, {2, 3, 16, 18, 21, 20, 9, 8}},       // F
		{{12, 13, 15, 14}, {0, 1, 11, 9, 22, 23, 18, 16}},  // B
		{{16, 17, 19, 18}, {0, 2, 4, 6, 20, 22, 15, 13}},   // L
		{{8, 9, 11, 10}, {1, 3, 14, 12, 23, 21, 7, 5}},     // R
	}

	for _, rot := range rotations {
		// clockwise
		if solved(rotate(cube, rot)) {
			fmt.Println("YES")
			return
		}
		// counter-clockwise: apply 3 times
		cc := rotate(cube, rot)
		cc = rotate(cc, rot)
		cc = rotate(cc, rot)
		if solved(cc) {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}
