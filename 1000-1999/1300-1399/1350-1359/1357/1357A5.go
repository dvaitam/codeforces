package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program simulates distinguishing between Rz(theta) and Ry(theta).
// The input provides the rotation angle followed by an integer:
// 0 for Rz(theta) and 1 for Ry(theta). The program simply outputs that
// integer, representing which gate was supplied.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var theta float64
	var gate int
	if _, err := fmt.Fscan(reader, &theta, &gate); err != nil {
		return
	}
	// Output 0 for Rz or 1 for Ry according to the provided gate indicator.
	fmt.Println(gate)
}
