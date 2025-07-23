package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// Given the side lengths of a hexagon with equal 120 degree angles, it
// computes the number of unit equilateral triangles that fit inside.
func main() {
	in := bufio.NewReader(os.Stdin)
	var a [6]int
	for i := 0; i < 6; i++ {
		if _, err := fmt.Fscan(in, &a[i]); err != nil {
			return
		}
	}
	s := a[0] + a[1] + a[2]
	ans := s*s - a[0]*a[0] - a[2]*a[2] - a[4]*a[4]
	fmt.Println(ans)
}
