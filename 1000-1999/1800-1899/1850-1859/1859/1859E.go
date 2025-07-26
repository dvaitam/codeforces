package main

import (
	"bufio"
	"fmt"
	"os"
)

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	const INF int64 = -1 << 60

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n+1)
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &b[i])
		}

		diff := make([]int64, n+1)
		sum := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			diff[i] = b[i] - a[i]
			sum[i] = a[i] + b[i]
		}

		dpPrev := make([]int64, k+1)
		dpCur := make([]int64, k+1)
		for j := 1; j <= k; j++ {
			dpPrev[j] = INF
		}
		size := k + n + 5
		off := n + 2
		M1 := make([]int64, size)
		M2 := make([]int64, size)
		M3 := make([]int64, size)
		M4 := make([]int64, size)
		for i := 0; i < size; i++ {
			M1[i] = INF
			M2[i] = INF
			M3[i] = INF
			M4[i] = INF
		}

		for r := 1; r <= n; r++ {
			// incorporate new starting point l = r
			for p := 0; p <= k; p++ {
				val := dpPrev[p]
				if val == INF {
					continue
				}
				q := p - (r - 1)
				idx := q + off
				if idx < 0 || idx >= size {
					continue
				}
				v := val + diff[r]
				if v > M1[idx] {
					M1[idx] = v
				}
				v = val - diff[r]
				if v > M2[idx] {
					M2[idx] = v
				}
				v = val + sum[r]
				if v > M3[idx] {
					M3[idx] = v
				}
				v = val - sum[r]
				if v > M4[idx] {
					M4[idx] = v
				}
			}

			for j := 0; j <= k; j++ {
				best := dpPrev[j]
				q := j - r
				idx := q + off
				if idx >= 0 && idx < size {
					if M1[idx] != INF {
						if tmp := M1[idx] + diff[r]; tmp > best {
							best = tmp
						}
					}
					if M2[idx] != INF {
						if tmp := M2[idx] - diff[r]; tmp > best {
							best = tmp
						}
					}
					if M3[idx] != INF {
						if tmp := M3[idx] - sum[r]; tmp > best {
							best = tmp
						}
					}
					if M4[idx] != INF {
						if tmp := M4[idx] + sum[r]; tmp > best {
							best = tmp
						}
					}
				}
				dpCur[j] = best
			}
			// prepare for next iteration
			for j := 0; j <= k; j++ {
				dpPrev[j] = dpCur[j]
				dpCur[j] = INF
			}
		}
		fmt.Fprintln(out, dpPrev[k])
	}
}
