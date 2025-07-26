package main

import (
	"bufio"
	"fmt"
	"os"
)

func sieve(limit int) []bool {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	return isPrime
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	const maxA = 1000000
	prime := sieve(maxA)

	for ; t > 0; t-- {
		var n, e int
		fmt.Fscan(in, &n, &e)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		var ans int64
		for r := 0; r < e; r++ {
			// collect sequence with indices r, r+e, r+2e, ...
			seq := make([]int, 0)
			for i := r; i < n; i += e {
				seq = append(seq, arr[i])
			}
			m := len(seq)
			if m == 0 {
				continue
			}
			left := make([]int, m)
			for i := 0; i < m; i++ {
				if seq[i] == 1 {
					if i > 0 {
						left[i] = left[i-1] + 1
					} else {
						left[i] = 1
					}
				} else {
					left[i] = 0
				}
			}
			right := make([]int, m)
			for i := m - 1; i >= 0; i-- {
				if seq[i] == 1 {
					if i+1 < m {
						right[i] = right[i+1] + 1
					} else {
						right[i] = 1
					}
				} else {
					right[i] = 0
				}
			}
			for i := 0; i < m; i++ {
				val := seq[i]
				if val <= maxA && prime[val] {
					l := 0
					if i > 0 {
						l = left[i-1]
					}
					rCnt := 0
					if i+1 < m {
						rCnt = right[i+1]
					}
					ans += int64(l*rCnt + l + rCnt)
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
