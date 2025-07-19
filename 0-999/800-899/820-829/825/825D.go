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

	var s, t string
	fmt.Fscan(reader, &s, &t)
	n := len(s)
	m := len(t)

	const C = 26
	tc := make([]int, C)
	for i := 0; i < m; i++ {
		tc[int(t[i]-'a')]++
	}
	sc := make([]int, C)
	scQ := 0
	for i := 0; i < n; i++ {
		if s[i] == '?' {
			scQ++
		} else {
			sc[int(s[i]-'a')]++
		}
	}

	poss := func(suit int) bool {
		tot := 0
		for i := 0; i < C; i++ {
			need := suit*tc[i] - sc[i]
			if need > 0 {
				tot += need
			}
			if tot > scQ {
				return false
			}
		}
		return true
	}

	lo, hi := 0, n
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if poss(mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}

	// calculate deficits
	deficit := make([]int, C)
	for i := 0; i < C; i++ {
		need := lo*tc[i] - sc[i]
		if need > 0 {
			deficit[i] = need
		}
	}

	// build result
	res := []byte(s)
	j := 0
	for i := 0; i < n; i++ {
		if res[i] == '?' {
			for j < C && deficit[j] == 0 {
				j++
			}
			if j < C {
				res[i] = byte('a' + j)
				deficit[j]--
			} else {
				res[i] = byte('t')
			}
		}
	}
	fmt.Fprintln(writer, string(res))
}
