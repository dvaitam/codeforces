package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit int64 = 1_000_000_000
const infLCM int64 = limit + 1

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lcmCap(a, b int64) int64 {
	if a == infLCM || b == infLCM {
		return infLCM
	}
	g := gcd(a, b)
	a /= g
	res := a * b
	if res > infLCM {
		return infLCM
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int64, n)
	present := make(map[int64]bool)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		present[arr[i]] = true
	}

	dp := make(map[int64]int)
	for _, x := range arr {
		next := make(map[int64]int)
		if next[x] < 1 {
			next[x] = 1
		}
		for l, lLen := range dp {
			nl := lcmCap(l, x)
			if lLen+1 > next[nl] {
				next[nl] = lLen + 1
			}
		}
		for l, lLen := range next {
			if lLen > dp[l] {
				dp[l] = lLen
			}
		}
	}

	ans := 0
	for l, lLen := range dp {
		if l > limit || !present[l] {
			if lLen > ans {
				ans = lLen
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
