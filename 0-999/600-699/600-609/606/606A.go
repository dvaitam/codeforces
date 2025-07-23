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

	var a, b, c int
	var x, y, z int
	if _, err := fmt.Fscan(in, &a, &b, &c); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &x, &y, &z); err != nil {
		return
	}

	spare := 0
	need := 0

	if a >= x {
		spare += (a - x) / 2
	} else {
		need += x - a
	}
	if b >= y {
		spare += (b - y) / 2
	} else {
		need += y - b
	}
	if c >= z {
		spare += (c - z) / 2
	} else {
		need += z - c
	}

	if spare >= need {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
