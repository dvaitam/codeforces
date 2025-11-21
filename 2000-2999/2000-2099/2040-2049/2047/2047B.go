package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var fact [11]int64

func init() {
	fact[0] = 1
	for i := 1; i <= 10; i++ {
		fact[i] = fact[i-1] * int64(i)
	}
}

func permCount(b []byte) int64 {
	freq := [26]int{}
	for _, ch := range b {
		freq[ch-'a']++
	}
	ans := fact[len(b)]
	for _, f := range freq {
		if f > 1 {
			ans /= fact[f]
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		best := ""
		bestVal := int64(math.MaxInt64)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				bytes := []byte(s)
				bytes[i] = s[j]
				val := permCount(bytes)
				if val < bestVal {
					bestVal = val
					best = string(bytes)
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}
