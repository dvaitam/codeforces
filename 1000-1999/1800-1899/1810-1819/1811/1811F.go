package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		adj := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		k := int(math.Round(math.Sqrt(float64(n))))
		if k*k != n || k < 3 || m != n+k {
			fmt.Fprintln(writer, "NO")
			continue
		}

		deg := make([]int, n)
		for i := 0; i < n; i++ {
			deg[i] = len(adj[i])
		}
		central := []int{}
		isCentral := make([]bool, n)
		valid := true
		for i, d := range deg {
			if d == 4 {
				central = append(central, i)
				isCentral[i] = true
			} else if d != 2 {
				valid = false
				break
			}
		}
		if !valid || len(central) != k {
			fmt.Fprintln(writer, "NO")
			continue
		}

		centralEdges := map[[2]int]struct{}{}
		for _, u := range central {
			cnt := 0
			for _, v := range adj[u] {
				if isCentral[v] {
					cnt++
					if u < v {
						centralEdges[[2]int{u, v}] = struct{}{}
					}
				}
			}
			if cnt != 2 {
				valid = false
				break
			}
		}
		if !valid || len(centralEdges) != k {
			fmt.Fprintln(writer, "NO")
			continue
		}

		visited := make([]bool, n)
		queue := []int{central[0]}
		visited[central[0]] = true
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				if isCentral[v] && !visited[v] {
					visited[v] = true
					queue = append(queue, v)
				}
			}
		}
		cntVisited := 0
		for _, u := range central {
			if visited[u] {
				cntVisited++
			}
		}
		if cntVisited != k {
			fmt.Fprintln(writer, "NO")
			continue
		}

		visitedAll := make([]bool, n)
		compCount := 0
		for i := 0; i < n; i++ {
			if !visitedAll[i] {
				compCount++
				q := []int{i}
				visitedAll[i] = true
				size := 0
				centralCnt := 0
				for len(q) > 0 {
					u := q[0]
					q = q[1:]
					size++
					if isCentral[u] {
						centralCnt++
					}
					for _, v := range adj[u] {
						if isCentral[u] && isCentral[v] {
							continue
						}
						if !visitedAll[v] {
							visitedAll[v] = true
							q = append(q, v)
						}
					}
				}
				if centralCnt != 1 || size != k {
					valid = false
					break
				}
			}
		}
		if !valid || compCount != k {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
