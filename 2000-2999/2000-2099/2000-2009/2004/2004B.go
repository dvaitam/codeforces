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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	buckets := make([][]uint8, 101)

	for ; t > 0; t-- {
		var l, r, L, R int
		fmt.Fscan(in, &l, &r)
		fmt.Fscan(in, &L, &R)

		for i := range buckets {
			buckets[i] = buckets[i][:0]
		}

		for a := l; a <= r; a++ {
			for b := L; b <= R; b++ {
				if a == b {
					continue
				}
				if a < b {
					hi := b - 1
					buckets[hi] = append(buckets[hi], uint8(a))
				} else {
					hi := a - 1
					buckets[hi] = append(buckets[hi], uint8(b))
				}
			}
		}

		ans := 0
		chosen := -1
		for hi := 1; hi <= 99; hi++ {
			bucket := buckets[hi]
			if len(bucket) == 0 {
				continue
			}
			need := false
			for _, lo := range bucket {
				if chosen < int(lo) {
					need = true
					break
				}
			}
			if need {
				ans++
				chosen = hi
			}
		}

		fmt.Fprintln(out, ans)
	}
}
