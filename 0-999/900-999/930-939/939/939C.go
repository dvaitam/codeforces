package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var s, f int
	fmt.Fscan(reader, &s, &f)
	k := f - s

	prefix := make([]int64, n+k+1)
	for i := 1; i <= n+k; i++ {
		prefix[i] = prefix[i-1] + a[(i-1)%n]
	}

	bestSum := int64(-1)
	bestTime := 1
	for start := 1; start <= n; start++ {
		sum := prefix[start+k-1] - prefix[start-1]
		time := (s-start+n)%n + 1
		if sum > bestSum || (sum == bestSum && time < bestTime) {
			bestSum = sum
			bestTime = time
		}
	}

	fmt.Fprintln(writer, bestTime)
}
