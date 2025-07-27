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
	for ; t > 0; t-- {
		var n, k, l int
		fmt.Fscan(in, &n, &k, &l)
		d := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
		}

		p := make([]int, 0, 2*k)
		for i := 0; i < k; i++ {
			p = append(p, i)
		}
		p = append(p, k)
		for i := k - 1; i >= 1; i-- {
			p = append(p, i)
		}
		cycle := len(p)

		type state struct{ pos, tm int }
		queue := make([]state, 0)
		vis := make([][]bool, n+2)
		for i := range vis {
			vis[i] = make([]bool, cycle)
		}
		push := func(s state) {
			if !vis[s.pos][s.tm] {
				vis[s.pos][s.tm] = true
				queue = append(queue, s)
			}
		}
		push(state{0, 0})
		found := false
		for len(queue) > 0 && !found {
			cur := queue[0]
			queue = queue[1:]
			pos, tm := cur.pos, cur.tm
			if pos == n+1 {
				found = true
				break
			}
			nt := (tm + 1) % cycle
			// wait at same position
			if pos == 0 || pos == n+1 {
				push(state{pos, nt})
			} else if d[pos-1]+p[nt] <= l {
				push(state{pos, nt})
			}
			// move forward
			np := pos + 1
			if np == n+1 {
				push(state{np, nt})
			} else if np <= n {
				if d[np-1]+p[nt] <= l {
					push(state{np, nt})
				}
			}
		}
		if found {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
