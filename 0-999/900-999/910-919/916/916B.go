package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func minCount(n int64, y int, limit int64) int64 {
	var count int64
	for i := 0; i < 61; i++ {
		if (n>>uint(i))&1 == 1 {
			if i > y {
				diff := i - y
				if diff >= 61 {
					return limit + 1
				}
				add := int64(1) << uint(diff)
				count += add
			} else {
				count++
			}
			if count > limit {
				return count
			}
		}
	}
	return count
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	pop := int64(bits.OnesCount64(uint64(n)))
	if k < pop {
		fmt.Fprintln(out, "No")
		return
	}

	low := -60
	high := 60
	for high-low > 1 {
		mid := (low + high) / 2
		if minCount(n, mid, k) <= k {
			high = mid
		} else {
			low = mid
		}
	}
	y := high

	cnt := make(map[int]int64)
	for i := 0; n>>uint(i) > 0; i++ {
		if (n>>uint(i))&1 == 1 {
			cnt[i]++
		}
	}

	for e := 60; e > y; e-- {
		if c, ok := cnt[e]; ok && c > 0 {
			cnt[e-1] += c * 2
			delete(cnt, e)
		}
	}

	total := int64(0)
	minExp := y
	for e, c := range cnt {
		total += c
		if e < minExp {
			minExp = e
		}
	}

	for total < k {
		for cnt[minExp] == 0 {
			minExp--
		}
		cnt[minExp]--
		cnt[minExp-1] += 2
		total++
		if minExp-1 < minExp {
			minExp = minExp - 1
		}
	}

	exps := make([]int, 0, k)
	for e := y; len(exps) < int(k); e-- {
		if c, ok := cnt[e]; ok {
			for i := int64(0); i < c; i++ {
				exps = append(exps, e)
			}
		}
	}

	fmt.Fprintln(out, "Yes")
	for i, v := range exps {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
