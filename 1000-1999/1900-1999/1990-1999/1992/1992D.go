package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

type state struct {
	pos   int
	water bool
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		var s string
		fmt.Fscan(in, &s)
		river := make([]byte, n+2)
		river[0] = 'L'
		for i := 0; i < n; i++ {
			river[i+1] = s[i]
		}
		river[n+1] = 'L'

		const INF = int(1e9)
		surf := make([]int, n+2)
		wat := make([]int, n+2)
		for i := 0; i < n+2; i++ {
			surf[i] = INF
			wat[i] = INF
		}
		dq := list.New()
		surf[0] = 0
		dq.PushFront(state{0, false})
		for dq.Len() > 0 {
			cur := dq.Remove(dq.Front()).(state)
			var d int
			if cur.water {
				d = wat[cur.pos]
			} else {
				d = surf[cur.pos]
			}
			if cur.water {
				nxt := cur.pos + 1
				if nxt <= n+1 && river[nxt] != 'C' {
					nw := river[nxt] == 'W'
					nd := d + 1
					if nw {
						if nd < wat[nxt] {
							wat[nxt] = nd
							dq.PushBack(state{nxt, true})
						}
					} else {
						if nd < surf[nxt] {
							surf[nxt] = nd
							dq.PushBack(state{nxt, false})
						}
					}
				}
			} else {
				for j := cur.pos + 1; j <= n+1 && j <= cur.pos+m; j++ {
					if river[j] == 'C' {
						continue
					}
					nw := river[j] == 'W'
					if nw {
						if d < wat[j] {
							wat[j] = d
							dq.PushFront(state{j, true})
						}
					} else {
						if d < surf[j] {
							surf[j] = d
							dq.PushFront(state{j, false})
						}
					}
				}
			}
		}
		ans := surf[n+1]
		if wat[n+1] < ans {
			ans = wat[n+1]
		}
		if ans <= k {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
