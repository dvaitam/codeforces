package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ u, v int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	edges := make(map[pair]byte)
	opp, same := 0, 0

	for i := 0; i < m; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "+" {
			var u, v int
			var c string
			fmt.Fscan(in, &u, &v, &c)
			edges[pair{u, v}] = c[0]
			if ch, ok := edges[pair{v, u}]; ok {
				opp++
				if ch == c[0] {
					same++
				}
			}
		} else if op == "-" {
			var u, v int
			fmt.Fscan(in, &u, &v)
			c := edges[pair{u, v}]
			delete(edges, pair{u, v})
			if ch, ok := edges[pair{v, u}]; ok {
				opp--
				if ch == c {
					same--
				}
			}
		} else {
			var k int
			fmt.Fscan(in, &k)
			if k%2 == 1 {
				if opp > 0 {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			} else {
				if same > 0 {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			}
		}
	}
}
