package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	// Skip MIC values which are irrelevant for this naive generator
	var a, b, c float64
	var d, e float64
	var P int
	fmt.Fscan(in, &a, &b)     // right direction mean, std
	fmt.Fscan(in, &c, &d)     // transition mean, std
	fmt.Fscan(in, &e, &b, &P) // correlation mean, std, window size

	// Start node connects to node 1 (if exists)
	if n > 0 {
		fmt.Println(1)
	} else {
		// no internal nodes, directly end
		fmt.Println(0)
		return
	}

	extra := k - n
	if extra < 0 {
		extra = 0
	}

	// Node 1: self loop on left to consume extra steps, then go right to node 2
	for i := 1; i <= n; i++ {
		if i == 1 {
			left := 1
			right := 0
			if n > 1 {
				right = 2
			}
			fmt.Printf("%d %d ", left, right)
			// vector: extra times 0 (left loops), then 1 to go right
			for j := 0; j < extra; j++ {
				fmt.Print("0")
			}
			fmt.Println("1")
		} else if i < n {
			left := 0
			right := i + 1
			fmt.Printf("%d %d 1\n", left, right)
		} else {
			// last node
			fmt.Printf("0 0 1\n")
		}
	}
}
