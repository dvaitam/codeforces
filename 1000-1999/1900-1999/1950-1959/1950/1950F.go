package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func calc(a, b, c, cap int) int {
	height := 0
	for a+b+c > 0 {
		if cap <= 0 {
			return math.MaxInt32
		}
		height++
		useA := a
		if useA > cap {
			useA = cap
		}
		cap -= useA
		a -= useA
		newCap := useA * 2
		useB := b
		if useB > cap {
			useB = cap
		}
		cap -= useB
		b -= useB
		newCap += useB
		useC := c
		if useC > cap {
			useC = cap
		}
		cap -= useC
		c -= useC
		cap = newCap
	}
	return height
}

func minHeight(a, b, c int) int {
	if c != a+1 {
		return -1
	}
	best := math.MaxInt32
	if a > 0 {
		if h := calc(a-1, b, c, 2); h < best {
			best = h
		}
	}
	if b > 0 {
		if h := calc(a, b-1, c, 1); h < best {
			best = h
		}
	}
	if c > 0 {
		if h := calc(a, b, c-1, 0); h < best {
			best = h
		}
	}
	if best == math.MaxInt32 {
		return -1
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		fmt.Fprintln(writer, minHeight(a, b, c))
	}
}
