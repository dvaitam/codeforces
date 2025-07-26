package main

import (
	"bufio"
	"fmt"
	"os"
)

type node struct{ r, a int }

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		circle := make([][360]bool, 21) // circle[r][a]
		radial := make([][360]bool, 20) // radial[r][a]
		for i := 0; i < n; i++ {
			var typ string
			fmt.Fscan(in, &typ)
			if typ == "C" {
				var r, t1, t2 int
				fmt.Fscan(in, &r, &t1, &t2)
				for a := t1; a != t2; a = (a + 1) % 360 {
					circle[r][a] = true
				}
			} else {
				var r1, r2, theta int
				fmt.Fscan(in, &r1, &r2, &theta)
				for r := r1; r < r2; r++ {
					radial[r][theta] = true
				}
			}
		}
		visited := make([][360]bool, 21)
		q := []node{{0, 0}}
		visited[0][0] = true
		success := false
		for len(q) > 0 && !success {
			cur := q[0]
			q = q[1:]
			r, a := cur.r, cur.a
			if r == 20 {
				success = true
				break
			}
			// outward
			if r < 20 && !circle[r+1][a] && !visited[r+1][a] {
				visited[r+1][a] = true
				q = append(q, node{r + 1, a})
			}
			// inward
			if r > 0 && !circle[r][a] && !visited[r-1][a] {
				visited[r-1][a] = true
				q = append(q, node{r - 1, a})
			}
			// clockwise
			na := (a + 1) % 360
			if !radial[r][na] && !visited[r][na] {
				visited[r][na] = true
				q = append(q, node{r, na})
			}
			// counter-clockwise
			nb := (a + 359) % 360
			if !radial[r][a] && !visited[r][nb] {
				visited[r][nb] = true
				q = append(q, node{r, nb})
			}
		}
		if success {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}
