package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAX = 1000000

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)

	present := make([]bool, MAX+1)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		present[x] = true
	}

	cnt := make([]int, MAX+1)
	gval := make([]int, MAX+1)

	for g := 1; g <= MAX; g++ {
		val := 0
		c := 0
		for m := g; m <= MAX; m += g {
			if present[m] {
				c++
				if val == 0 {
					val = m
				} else {
					val = gcd(val, m)
				}
			}
		}
		cnt[g] = c
		gval[g] = val
	}

	ans := 0
	for g := 1; g <= MAX; g++ {
		if cnt[g] >= 2 && gval[g] == g && !present[g] {
			ans++
		}
	}
	fmt.Fprintln(writer, ans)
}
