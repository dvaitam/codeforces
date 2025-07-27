package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement efficient solution for problem G.
// Current implementation is a placeholder that reads the input
// and outputs 0. The full algorithm is non-trivial and should
// count the number of required time travels modulo 1e9+7.

const mod int64 = 1_000_000_007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i], &b[i])
	}

	var t int
	fmt.Fscan(in, &t)
	s := make([]int, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &s[i])
		s[i]--
	}

	// Placeholder: output 0 until proper algorithm is implemented
	fmt.Println(0)
}
