package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var x [3]int64
	var y [3]int64
	for i := 0; i < 3; i++ {
		if _, err := fmt.Fscan(in, &x[i], &y[i]); err != nil {
			return
		}
	}

	// Check if all points lie on the same vertical or horizontal line.
	if (x[0] == x[1] && x[1] == x[2]) || (y[0] == y[1] && y[1] == y[2]) {
		fmt.Println(1)
		return
	}

	// Check if there exists a point that can serve as the single corner
	// of an L-shaped polyline passing through the other two points.
	for k := 0; k < 3; k++ {
		i := (k + 1) % 3
		j := (k + 2) % 3
		if x[i] == x[k] && y[j] == y[k] {
			fmt.Println(2)
			return
		}
		if y[i] == y[k] && x[j] == x[k] {
			fmt.Println(2)
			return
		}
	}

	// If at least one pair shares x or y coordinate, we need three segments.
	if x[0] == x[1] || x[0] == x[2] || x[1] == x[2] ||
		y[0] == y[1] || y[0] == y[2] || y[1] == y[2] {
		fmt.Println(3)
	} else {
		fmt.Println(4)
	}
}
