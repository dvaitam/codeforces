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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		var sStr, tStr string
		fmt.Fscan(in, &sStr)
		fmt.Fscan(in, &tStr)

		adj := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		s := make([]int, n)
		tgt := make([]int, n)
		for i := 0; i < n; i++ {
			if sStr[i] == '1' {
				s[i] = 1
			}
			if tStr[i] == '1' {
				tgt[i] = 1
			}
		}

		visited := make([]bool, n)
		color := make([]int, n)
		for i := range color {
			color[i] = -1
		}

		comps := make([][4]int, 0)
		ok := true

		for i := 0; i < n && ok; i++ {
			if visited[i] {
				continue
			}
			queue := []int{i}
			visited[i] = true
			color[i] = 0
			nodes := []int{i}
			bip := true

			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]

				for _, v := range adj[u] {
					if color[v] == -1 {
						color[v] = color[u] ^ 1
						visited[v] = true
						queue = append(queue, v)
						nodes = append(nodes, v)
					} else if color[v] == color[u] {
						bip = false
					}
				}
			}

			sumS := 0
			sumT := 0
			for _, v := range nodes {
				sumS += s[v]
				sumT += tgt[v]
			}
			if sumS != sumT {
				ok = false
				break
			}

			if bip {
				s0 := 0
				t0 := 0
				for _, v := range nodes {
					if color[v] == 0 {
						s0 += s[v]
						t0 += tgt[v]
					}
				}
				s1 := sumS - s0
				t1 := sumT - t0
				comps = append(comps, [4]int{s0, t0, s1, t1})
			}
		}

		if !ok {
			fmt.Fprintln(out, "No")
			continue
		}

		if len(comps) == 0 {
			fmt.Fprintln(out, "Yes")
			continue
		}

		// check even parity possibility
		canEven := true
		canOdd := true
		for _, info := range comps {
			s0, t0, s1, t1 := info[0], info[1], info[2], info[3]
			if s0 != t0 || s1 != t1 {
				canEven = false
			}
			if s0 != t1 || s1 != t0 {
				canOdd = false
			}
		}

		if canEven || canOdd {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
