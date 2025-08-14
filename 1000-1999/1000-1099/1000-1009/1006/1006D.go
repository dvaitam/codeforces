package main

import (
	"bufio"
	"fmt"
	"os"
)

// get returns the minimum number of changes required to make the four
// characters form two equal pairs. Only the first and third characters
// (belonging to the first string) may be changed.
func get(a, b, c, d byte) int {
	if b == d {
		if a == c {
			return 0
		}
		return 1
	}
	cost1 := 0
	if a != b {
		cost1++
	}
	if c != d {
		cost1++
	}
	cost2 := 0
	if a != d {
		cost2++
	}
	if c != b {
		cost2++
	}
	if cost1 < cost2 {
		return cost1
	}
	return cost2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	// read n and two strings
	fmt.Fscan(reader, &n)
	s0 := make([]byte, 0, n)
	s1 := make([]byte, 0, n)
	// scan strings
	var str0, str1 string
	fmt.Fscan(reader, &str0, &str1)
	s0 = []byte(str0)
	s1 = []byte(str1)
	ans := 0
	// process pairs
	for i := 0; i < n/2; i++ {
		ans += get(s0[i], s1[i], s0[n-i-1], s1[n-i-1])
	}
	// middle column if odd
	if n%2 == 1 && s0[n/2] != s1[n/2] {
		ans++
	}
	fmt.Println(ans)
}
