package main

import (
	"bufio"
	"fmt"
	"os"
)

func simulate(start int, steps int64) int {
	pos := start
	for i := int64(0); i < steps; i++ {
		if i%2 == 0 {
			if pos == 0 {
				pos = 1
			} else if pos == 1 {
				pos = 0
			}
		} else {
			if pos == 1 {
				pos = 2
			} else if pos == 2 {
				pos = 1
			}
		}
	}
	return pos
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	var x int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &x); err != nil {
		return
	}

	steps := n % 6
	for start := 0; start < 3; start++ {
		if simulate(start, steps) == x {
			fmt.Fprintln(out, start)
			return
		}
	}
}
