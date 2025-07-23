package main

import (
	"bufio"
	"fmt"
	"os"
)

func grundyOdd(x int64) int {
	switch x {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 0
	case 3:
		return 1
	case 4:
		return 2
	}
	if x%2 == 1 {
		return 0
	}
	g := grundyOdd(x / 2)
	if g == 1 {
		return 2
	}
	return 1
}

func grundyEven(x int64) int {
	if x == 1 {
		return 1
	}
	if x == 2 {
		return 2
	}
	if x%2 == 1 {
		return 0
	}
	return 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	res := 0
	for i := 0; i < n; i++ {
		var a int64
		fmt.Fscan(in, &a)
		var g int
		if k%2 == 0 {
			g = grundyEven(a)
		} else {
			g = grundyOdd(a)
		}
		res ^= g
	}
	if res != 0 {
		fmt.Fprintln(out, "Kevin")
	} else {
		fmt.Fprintln(out, "Nicky")
	}
}
