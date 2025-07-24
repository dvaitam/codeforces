package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxL = 21

var fact [maxL + 1]int64

func init() {
	fact[0] = 1
	for i := 1; i <= maxL; i++ {
		fact[i] = fact[i-1] * int64(i)
	}
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] / (fact[k] * fact[n-k])
}

func countFromCounts(cnt []int, L int) int64 {
	for v, c := range cnt {
		if v >= L && c > 0 {
			return 0
		}
	}
	used := 0
	ways := int64(1)
	for v := L - 1; v >= 0; v-- {
		c := cnt[v]
		if c == 0 {
			continue
		}
		allowed := L - v
		if v == 0 {
			allowed = L - 1
		}
		avail := allowed - used
		if c > avail {
			return 0
		}
		ways *= comb(int64(avail), int64(c))
		used += c
	}
	return ways
}

func countPermutations(digits []int) int64 {
	var cnt [maxL]int
	for _, d := range digits {
		cnt[d]++
	}
	L := len(digits) + 1
	total := countFromCounts(cnt[:], L)
	if cnt[0] == 0 {
		return total
	}
	cnt[0]--
	sub := countFromCounts(cnt[:], L-1)
	return total - sub
}

func main() {
	rdr := bufio.NewReader(os.Stdin)
	wtr := bufio.NewWriter(os.Stdout)
	defer wtr.Flush()

	var t int
	fmt.Fscan(rdr, &t)
	for ; t > 0; t-- {
		var k int64
		fmt.Fscan(rdr, &k)
		var digits []int
		div := int64(2)
		for x := k; x > 0; div++ {
			digits = append(digits, int(x%div))
			x /= div
		}
		ans := countPermutations(digits) - 1
		fmt.Fprintln(wtr, ans)
	}
}
