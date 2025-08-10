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
	var n, k int
	fmt.Fscan(in, &n, &k)
	a := make([]int, n)
	b := make([]int, n)
	var sumR, sumB int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i], &b[i])
		sumR += int64(a[i])
		sumB += int64(b[i])
	}
	total := (sumR + sumB) / int64(k)
	if total == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	dp := make([]bool, k)
	dp[0] = true
	for i := 0; i < n; i++ {
		possible := make([]int, 0, k)
		for r := 0; r < k; r++ {
			if r <= a[i] && (k-r)%k <= b[i] {
				possible = append(possible, r)
			}
		}
		next := make([]bool, k)
		for prev := 0; prev < k; prev++ {
			if !dp[prev] {
				continue
			}
			for _, r := range possible {
				next[(prev+r)%k] = true
			}
		}
		dp = next
	}
	remR := int(sumR % int64(k))
	remB := int(sumB % int64(k))
	can := false
	for remX := 0; remX < k; remX++ {
		if dp[remX] {
			remY := (k - remX) % k
			remRPool := (remR - remX + k) % k
			remBPool := (remB - remY + k) % k
			if remRPool+remBPool < k {
				can = true
				break
			}
		}
	}
	if can {
		fmt.Fprintln(out, total)
	} else {
		fmt.Fprintln(out, total-1)
	}
}
