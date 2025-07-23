package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// It computes the minimal time to drive from the beginning of the street
// to house a on a street with houses on both sides as described in the
// statement.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a int
	if _, err := fmt.Fscan(in, &n, &a); err != nil {
		return
	}
	var result int
	if a%2 == 1 {
		result = (a + 1) / 2
	} else {
		result = (n-a)/2 + 1
	}
	fmt.Println(result)
}
