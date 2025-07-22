package main

import (
	"bufio"
	"fmt"
	"os"
)

type preData struct {
	prefBC []int
	prefB  []int
	prefC  []int
	tailA  []int
}

func preprocess(s string) preData {
	n := len(s)
	prefBC := make([]int, n+1)
	prefB := make([]int, n+1)
	prefC := make([]int, n+1)
	tail := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefBC[i] = prefBC[i-1]
		prefB[i] = prefB[i-1]
		prefC[i] = prefC[i-1]
		tail[i] = 0
		switch s[i-1] {
		case 'B':
			prefBC[i]++
			prefB[i]++
		case 'C':
			prefBC[i]++
			prefC[i]++
		case 'A':
			tail[i] = tail[i-1] + 1
		}
	}
	return preData{prefBC, prefB, prefC, tail}
}

func getTail(data preData, l, r int) int {
	t := data.tailA[r]
	if t > r-l+1 {
		t = r - l + 1
	}
	return t
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &t)

	ds := preprocess(s)
	dt := preprocess(t)

	var Q int
	fmt.Fscan(in, &Q)
	result := make([]byte, Q)
	for i := 0; i < Q; i++ {
		var a, b, c, d int
		fmt.Fscan(in, &a, &b, &c, &d)

		bc1 := ds.prefBC[b] - ds.prefBC[a-1]
		bc2 := dt.prefBC[d] - dt.prefBC[c-1]
		diff1 := (ds.prefB[b] - ds.prefB[a-1] - (ds.prefC[b] - ds.prefC[a-1])) & 1
		diff2 := (dt.prefB[d] - dt.prefB[c-1] - (dt.prefC[d] - dt.prefC[c-1])) & 1
		tail1 := getTail(ds, a, b)
		tail2 := getTail(dt, c, d)

		ok := true
		if diff1 != diff2 {
			ok = false
		} else if bc1 > bc2 {
			ok = false
		} else if (bc2-bc1)%2 != 0 {
			ok = false
		} else if tail1 < tail2 {
			ok = false
		} else {
			if bc1 == 0 && bc2 > 0 && tail1 == tail2 {
				ok = false
			} else if bc1 == bc2 && (tail1-tail2)%3 != 0 {
				ok = false
			}
		}
		if ok {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}
	fmt.Fprintln(out, string(result))
}
