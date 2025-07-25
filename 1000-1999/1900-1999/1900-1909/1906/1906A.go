package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	grid := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
	}

	dirs := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	best := "~~~" // char '~~' is after 'Z'
	for x1 := 0; x1 < 3; x1++ {
		for y1 := 0; y1 < 3; y1++ {
			for _, d1 := range dirs {
				x2, y2 := x1+d1[0], y1+d1[1]
				if x2 < 0 || x2 >= 3 || y2 < 0 || y2 >= 3 {
					continue
				}
				for _, d2 := range dirs {
					x3, y3 := x2+d2[0], y2+d2[1]
					if x3 < 0 || x3 >= 3 || y3 < 0 || y3 >= 3 {
						continue
					}
					if x3 == x1 && y3 == y1 || x3 == x2 && y3 == y2 {
						continue
					}
					word := string([]byte{grid[x1][y1], grid[x2][y2], grid[x3][y3]})
					if word < best {
						best = word
					}
				}
			}
		}
	}

	fmt.Fprintln(out, best)
}
