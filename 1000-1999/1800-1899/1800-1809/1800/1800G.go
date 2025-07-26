package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	idMap map[string]int
	idCnt int
	adj   [][]int
)

func getID(children []int) int {
	sort.Ints(children)
	var sb strings.Builder
	sb.WriteByte('(')
	for i, v := range children {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte(')')
	key := sb.String()
	if id, ok := idMap[key]; ok {
		return id
	}
	idCnt++
	idMap[key] = idCnt
	return idCnt
}

func dfs(v, p int) (int, bool) {
	ids := make([]int, 0, len(adj[v]))
	symFlags := make([]bool, 0, len(adj[v]))
	for _, u := range adj[v] {
		if u == p {
			continue
		}
		id, sym := dfs(u, v)
		ids = append(ids, id)
		symFlags = append(symFlags, sym)
	}
	id := getID(ids)

	count := make(map[int]int)
	symCount := make(map[int]int)
	for i, cid := range ids {
		count[cid]++
		if symFlags[i] {
			symCount[cid]++
		}
	}
	oddGroups := 0
	for cid, c := range count {
		if c%2 == 1 {
			if symCount[cid] == 0 {
				return id, false
			}
			oddGroups++
			if oddGroups > 1 {
				return id, false
			}
		}
	}
	return id, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		adj = make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		idMap = make(map[string]int)
		idCnt = 0
		_, sym := dfs(1, 0)
		if sym {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
