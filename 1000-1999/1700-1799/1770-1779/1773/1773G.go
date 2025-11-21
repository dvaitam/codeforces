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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	size := 1 << m
	freq := make([]int64, size)

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		mask := 0
		for j := 0; j < m; j++ {
			if s[j] == '1' {
				mask |= 1 << j
			}
		}
		freq[mask]++
	}

	dpSubset := make([]int64, size)
	copy(dpSubset, freq)
	for b := 0; b < m; b++ {
		for mask := 0; mask < size; mask++ {
			if mask&(1<<b) != 0 {
				dpSubset[mask] += dpSubset[mask^(1<<b)]
			}
		}
	}

	dpSuperset := make([]int64, size)
	copy(dpSuperset, freq)
	for b := 0; b < m; b++ {
		for mask := 0; mask < size; mask++ {
			if mask&(1<<b) == 0 {
				dpSuperset[mask] += dpSuperset[mask|(1<<b)]
			}
		}
	}

	popcnt := make([]int, size)
	for mask := 1; mask < size; mask++ {
		popcnt[mask] = popcnt[mask&(mask-1)] + 1
	}

	leaderBit := 1
	f := make([]float64, size)
	H := make([]float64, size)
	fullMask := size - 1

	buckets := make([][]int, m+1)
	for mask := 1; mask < size; mask++ {
		if mask&leaderBit == 0 {
			continue
		}
		pc := popcnt[mask]
		buckets[pc] = append(buckets[pc], mask)
	}

	totalQuestions := int64(n)

	for pc := 1; pc <= m; pc++ {
		for _, mask := range buckets[pc] {
			maskWithoutLeader := mask &^ leaderBit

			sumLess := 0.0
			hPre := 0.0

			sub2 := maskWithoutLeader
			for {
				sub := sub2 | leaderBit
				if sub != mask {
					sumLess += float64(dpSuperset[sub]) * H[sub]
					sign := 1.0
					if ((pc - popcnt[sub]) & 1) != 0 {
						sign = -1.0
					}
					hPre += sign * f[sub]
				}
				if sub2 == 0 {
					break
				}
				sub2 = (sub2 - 1) & maskWithoutLeader
			}

			comp := fullMask ^ mask
			sumZero := dpSubset[comp]
			hMask := dpSuperset[mask]
			k := totalQuestions - sumZero - hMask

			if k == 0 {
				f[mask] = 1.0
			} else {
				numerator := sumLess + float64(hMask)*hPre
				f[mask] = numerator / float64(k)
			}
			H[mask] = hPre + f[mask]
		}
	}

	fmt.Fprintf(out, "%.15f\n", f[fullMask])
}
