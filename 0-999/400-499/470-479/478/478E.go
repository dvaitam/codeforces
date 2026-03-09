package main

import (
	"bufio"
	"fmt"
	"os"
)

func isWavy(x int64) bool {
	if x < 10 {
		return true
	}
	var d [20]int
	n := 0
	for tmp := x; tmp > 0; tmp /= 10 {
		d[n] = int(tmp % 10)
		n++
	}
	// reverse so d[0] is most significant
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
	// consecutive equal digits are not allowed (strict alternation required)
	for i := 1; i < n; i++ {
		if d[i] == d[i-1] {
			return false
		}
	}
	// for 3+ digit numbers, every middle digit must be a peak or valley
	for i := 1; i < n-1; i++ {
		peak := d[i] > d[i-1] && d[i] > d[i+1]
		valley := d[i] < d[i-1] && d[i] < d[i+1]
		if !peak && !valley {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	const limit = int64(1e9) // oracle only needs to handle small n,k for verification
	var cnt int64
	for x := n; x <= limit; x += n {
		if isWavy(x) {
			cnt++
			if cnt == k {
				fmt.Println(x)
				return
			}
		}
	}
	fmt.Println(-1)
}
