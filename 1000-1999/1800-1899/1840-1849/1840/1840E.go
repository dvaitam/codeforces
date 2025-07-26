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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s1, s2 string
		fmt.Fscan(in, &s1)
		fmt.Fscan(in, &s2)
		n := len(s1)
		b1 := []byte(s1)
		b2 := []byte(s2)

		var t, q int
		fmt.Fscan(in, &t, &q)

		events := make([][]int, q+t+5)
		blocked := make([]bool, n)
		mism := 0
		for i := 0; i < n; i++ {
			if b1[i] != b2[i] {
				mism++
			}
		}

		for time := 1; time <= q; time++ {
			for _, pos := range events[time] {
				if blocked[pos] {
					blocked[pos] = false
					if b1[pos] != b2[pos] {
						mism++
					}
				}
			}
			var typ int
			fmt.Fscan(in, &typ)
			switch typ {
			case 1:
				var pos int
				fmt.Fscan(in, &pos)
				pos--
				if !blocked[pos] {
					if b1[pos] != b2[pos] {
						mism--
					}
					blocked[pos] = true
					events[time+t] = append(events[time+t], pos)
				}
			case 2:
				var a1, p1, a2, p2 int
				fmt.Fscan(in, &a1, &p1, &a2, &p2)
				p1--
				p2--
				idxs := make(map[int]struct{}, 2)
				idxs[p1] = struct{}{}
				idxs[p2] = struct{}{}
				for idx := range idxs {
					if !blocked[idx] && b1[idx] != b2[idx] {
						mism--
					}
				}
				var c1, c2 byte
				if a1 == 1 {
					c1 = b1[p1]
				} else {
					c1 = b2[p1]
				}
				if a2 == 1 {
					c2 = b1[p2]
				} else {
					c2 = b2[p2]
				}
				if a1 == 1 {
					b1[p1] = c2
				} else {
					b2[p1] = c2
				}
				if a2 == 1 {
					b1[p2] = c1
				} else {
					b2[p2] = c1
				}
				for idx := range idxs {
					if !blocked[idx] && b1[idx] != b2[idx] {
						mism++
					}
				}
			case 3:
				if mism == 0 {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			}
		}
	}
}
