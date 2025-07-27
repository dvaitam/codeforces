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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	r := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &r[i])
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	// count tasks solved only by Robo-Coder and only by BionicSolver
	a := 0
	c := 0
	for i := 0; i < n; i++ {
		if r[i] == 1 && b[i] == 0 {
			a++
		} else if r[i] == 0 && b[i] == 1 {
			c++
		}
	}

	if a == 0 {
		fmt.Fprintln(out, -1)
		return
	}
	// minimum maximum points per problem is floor(c/a) + 1
	ans := c/a + 1
	fmt.Fprintln(out, ans)
}
