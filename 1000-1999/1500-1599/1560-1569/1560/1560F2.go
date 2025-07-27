package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

// Solution for problem F2 (see problemF2.txt).
// Given n and k (k <= 10), find the smallest integer >= n whose
// decimal representation uses no more than k distinct digits.

func countDistinct(mask int) int {
	return bits.OnesCount(uint(mask))
}

func nextBeautiful(n int, k int) int {
	s := strconv.Itoa(n)
	L := len(s)

	prefixMask := make([]int, L+1)
	for i := 0; i < L; i++ {
		d := int(s[i] - '0')
		prefixMask[i+1] = prefixMask[i] | (1 << d)
	}

	if countDistinct(prefixMask[L]) <= k {
		return n
	}

	best := int64(1<<63 - 1)

	for j := L - 1; j >= 0; j-- {
		mask := prefixMask[j]
		if countDistinct(mask) > k {
			continue
		}
		start := int(s[j]-'0') + 1
		for d := start; d < 10; d++ {
			m2 := mask | (1 << d)
			if countDistinct(m2) > k {
				continue
			}
			fill := byte('0')
			for t := 0; t < 10; t++ {
				m3 := m2
				if (m3 & (1 << t)) == 0 {
					m3 |= 1 << t
				}
				if countDistinct(m3) <= k {
					fill = byte('0' + t)
					break
				}
			}
			cand := s[:j] + string('0'+byte(d)) + strings.Repeat(string(fill), L-j-1)
			if val, err := strconv.ParseInt(cand, 10, 64); err == nil {
				if val < best {
					best = val
				}
			}
		}
	}

	if best != int64(1<<63-1) {
		return int(best)
	}

	if k == 1 {
		ans := 0
		for i := 0; i < L+1; i++ {
			ans = ans*10 + 1
		}
		return ans
	}

	pow := 1
	for i := 0; i < L; i++ {
		pow *= 10
	}
	return pow
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		fmt.Fprintln(out, nextBeautiful(n, k))
	}
}
