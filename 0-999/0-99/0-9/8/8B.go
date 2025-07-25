package main

import (
	"bufio"
	"fmt"
	"os"
)

// pt represents a point on the grid
type pt struct{ x, y int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	// track positions visited in order
	positions := make([]pt, 0, len(s)+1)
	positions = append(positions, pt{0, 0})
	visited := make(map[pt]int)
	visited[positions[0]] = 0
	x, y := 0, 0
	// build path and check for revisits
	for i, ch := range s {
		switch ch {
		case 'L':
			x--
		case 'R':
			x++
		case 'U':
			y++
		case 'D':
			y--
		}
		curr := pt{x, y}
		if _, ok := visited[curr]; ok {
			fmt.Println("BUG")
			return
		}
		positions = append(positions, curr)
		visited[curr] = i + 1
	}
	// directions for adjacency check
	dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	// ensure no shortcuts between non-consecutive points
	for i, p := range positions {
		for _, d := range dirs {
			np := pt{p.x + d.x, p.y + d.y}
			if j, ok := visited[np]; ok {
				if j != i-1 && j != i+1 {
					fmt.Println("BUG")
					return
				}
			}
		}
	}
	fmt.Println("OK")
}
