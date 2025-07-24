package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			if u == v {
				continue
			}
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		for i := 1; i <= n; i++ {
			sort.Ints(adj[i])
		}
		groups := make(map[string][]int)
		for i := 1; i <= n; i++ {
			var sb strings.Builder
			for _, v := range adj[i] {
				sb.WriteString(strconv.Itoa(v))
				sb.WriteByte(',')
			}
			key := sb.String()
			groups[key] = append(groups[key], i)
		}
		found := false
		for _, g := range groups {
			if len(g) >= 3 && len(g)%3 == 0 && len(g) < n {
				fmt.Fprintln(out, "Yes")
				fmt.Fprintln(out, len(g))
				for idx, v := range g {
					if idx > 0 {
						fmt.Fprint(out, " ")
					}
					fmt.Fprint(out, v)
				}
				fmt.Fprintln(out)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintln(out, "No")
		}
	}
}
