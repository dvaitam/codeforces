package main

import (
	"bufio"
	"fmt"
	"os"
)

func beats(i, j, n int) bool {
	if i == 1 && j == n {
		return true
	}
	if j == 1 && i == n {
		return false
	}
	return i > j
}

func dfsState(aliceMask, bobMask int, n int, memo map[[2]int]bool) bool {
	state := [2]int{aliceMask, bobMask}
	if val, ok := memo[state]; ok {
		return val
	}
	if aliceMask == 0 {
		memo[state] = false
		return false
	}
	if bobMask == 0 {
		memo[state] = true
		return true
	}
	for i := 0; i < n; i++ {
		if (aliceMask>>i)&1 == 0 {
			continue
		}
		alNextMask := aliceMask & ^(1 << i)
		for j := 0; j < n; j++ {
			if (bobMask>>j)&1 == 0 {
				continue
			}
			bobNextMask := bobMask & ^(1 << j)
			if beats(i+1, j+1, n) {
				alNextMask |= 1<<i | 1<<j
			} else {
				bobNextMask |= 1<<i | 1<<j
			}
			if !dfsState(bobNextMask, alNextMask, n, memo) {
				memo[state] = true
				return true
			}
			alNextMask = alNextMask & ^(1 << i)
		}
	}
	memo[state] = false
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)

		aliceMask := 0
		bobMask := 0
		for i := 0; i < n; i++ {
			if s[i] == 'A' {
				aliceMask |= 1 << i
			} else {
				bobMask |= 1 << i
			}
		}

		memo := make(map[[2]int]bool)
		if dfsState(aliceMask, bobMask, n, memo) {
			fmt.Fprintln(writer, "Alice")
		} else {
			fmt.Fprintln(writer, "Bob")
		}
	}
}
