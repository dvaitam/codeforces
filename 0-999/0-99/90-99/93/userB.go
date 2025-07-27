package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const eps = 1e-8

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, w, m int
	if _, err := fmt.Fscan(in, &n, &w, &m); err != nil {
		return
	}

	if m > 2*n {
		fmt.Println("NO")
		return
	}

	s := float64(n*w) / float64(m) // required volume in each cup

	type pair struct {
		bottle int
		vol    float64
	}
	ans := make([][]pair, m)

	remain := float64(w) // milk left in current bottle
	cur := 0             // current bottle index (0-based)

	usedCnt := make([]int, n)  // how many cups each bottle was poured into
	lastCup := make([]int, n)  // last cup index this bottle was used for
	for i := 0; i < n; i++ {   // initialise to -1
		lastCup[i] = -1
	}

	for cup := 0; cup < m; cup++ {
		need := s
		for need > eps {
			if cur >= n { // no bottles left
				fmt.Println("NO")
				return
			}
			if remain < eps {
				cur++
				if cur >= n {
					fmt.Println("NO")
					return
				}
				remain = float64(w)
			}

			pour := math.Min(remain, need)

			if lastCup[cur] != cup {
				usedCnt[cur]++
				if usedCnt[cur] > 2 {
					fmt.Println("NO")
					return
				}
				lastCup[cur] = cup
			}

			ans[cup] = append(ans[cup], pair{cur + 1, pour})

			remain -= pour
			need -= pour
		}
	}

	fmt.Println("YES")
	out := bufio.NewWriter(os.Stdout)
	for _, cup := range ans {
		for i, p := range cup {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprintf(out, "%d %.6f", p.bottle, p.vol)
		}
		fmt.Fprintln(out)
	}
	out.Flush()
}