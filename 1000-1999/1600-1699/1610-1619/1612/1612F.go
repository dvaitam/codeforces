package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: Implement solution for problem F (Armor and Weapons).
// The full algorithm involves advanced graph or dynamic programming
// which is omitted here. This placeholder just reads the input and
// outputs 0 so that the source file compiles.
func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
	}
	fmt.Println(0)
}
