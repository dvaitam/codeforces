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
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	m, _ := strconv.Atoi(scanner.Text())

	type Pin struct {
		side string
		idx  int
	}

	type String struct {
		u Pin
		v Pin
	}

	stringsList := make([]String, n+m)
	hasTB := false
	hasRL := false

	for i := 0; i < n+m; i++ {
		scanner.Scan()
		s1 := scanner.Text()
		scanner.Scan()
		s2 := scanner.Text()
		scanner.Scan()
		x1, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		x2, _ := strconv.Atoi(scanner.Text())

		stringsList[i] = String{Pin{s1, x1}, Pin{s2, x2}}
		
		var sides [2]string
		sides[0] = s1
		sides[1] = s2
		sort.Strings(sides[:])
		edgeType := sides[0] + sides[1]
		if edgeType == "BT" {
			hasTB = true
		}
		if edgeType == "LR" {
			hasRL = true
		}
	}

	if hasTB && hasRL {
		fmt.Println("No solution")
		return
	}

	getTypeRank := func(s1, s2 string) int {
		var sides [2]string
		sides[0] = s1
		sides[1] = s2
		sort.Strings(sides[:])
		edgeType := sides[0] + sides[1]

		if hasRL {
			switch edgeType {
			case "LT": return 0
			case "RT": return 1
			case "LR": return 2
			case "BL": return 3
			case "BR": return 4
			}
		} else {
			switch edgeType {
			case "LT": return 0
			case "BL": return 1
			case "BT": return 2
			case "RT": return 3
			case "BR": return 4
			}
		}
		return -1
	}

	pinRank := make(map[Pin]int)
	for _, s := range stringsList {
		rank := getTypeRank(s.u.side, s.v.side)
		pinRank[s.u] = rank
		pinRank[s.v] = rank
	}

	numNodes := n + m
	adj := make([][]int, numNodes)

	addBiEdge := func(u, v int) {
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	for _, s := range stringsList {
		var u, v int
		if s.u.side == "L" || s.u.side == "R" {
			u = s.u.idx - 1
		} else {
			u = n + s.u.idx - 1
		}
		if s.v.side == "L" || s.v.side == "R" {
			v = s.v.idx - 1
		} else {
			v = n + s.v.idx - 1
		}
		addBiEdge(u, v)
	}

	type StrictEdge struct {
		u, v int
	}
	var strictEdges []StrictEdge

	addStrictGroup := func(nodesByRank map[int][]int) {
		var ranks []int
		for r := range nodesByRank {
			ranks = append(ranks, r)
		}
		sort.Ints(ranks)
		for i := 0; i < len(ranks)-1; i++ {
			r1 := ranks[i]
			r2 := ranks[i+1]
			d1 := numNodes
			numNodes++
			d2 := numNodes
			numNodes++
			adj = append(adj, []int{})
			adj = append(adj, []int{})
			for _, u := range nodesByRank[r1] {
				adj[u] = append(adj[u], d1)
			}
			adj[d1] = append(adj[d1], d2)
			strictEdges = append(strictEdges, StrictEdge{d1, d2})
			for _, v := range nodesByRank[r2] {
				adj[d2] = append(adj[d2], v)
			}
		}
	}

	nodesByRankL := make(map[int][]int)
	nodesByRankR := make(map[int][]int)
	for i := 1; i <= n; i++ {
		rL := pinRank[Pin{"L", i}]
		rR := pinRank[Pin{"R", i}]
		nodesByRankL[rL] = append(nodesByRankL[rL], i-1)
		nodesByRankR[rR] = append(nodesByRankR[rR], i-1)
	}
	addStrictGroup(nodesByRankL)
	addStrictGroup(nodesByRankR)

	nodesByRankT := make(map[int][]int)
	nodesByRankB := make(map[int][]int)
	for j := 1; j <= m; j++ {
		rT := pinRank[Pin{"T", j}]
		rB := pinRank[Pin{"B", j}]
		nodesByRankT[rT] = append(nodesByRankT[rT], n+j-1)
		nodesByRankB[rB] = append(nodesByRankB[rB], n+j-1)
	}
	addStrictGroup(nodesByRankT)
	addStrictGroup(nodesByRankB)

	idx := 0
	stack := []int{}
	inStack := make([]bool, numNodes)
	dfn := make([]int, numNodes)
	low := make([]int, numNodes)
	scc := make([]int, numNodes)
	sccCount := 0

	for i := 0; i < numNodes; i++ {
		dfn[i] = -1
	}

	var tarjan func(int)
	tarjan = func(u int) {
		dfn[u] = idx
		low[u] = idx
		idx++
		stack = append(stack, u)
		inStack[u] = true

		for _, v := range adj[u] {
			if dfn[v] == -1 {
				tarjan(v)
				if low[v] < low[u] {
					low[u] = low[v]
				}
			} else if inStack[v] {
				if dfn[v] < low[u] {
					low[u] = dfn[v]
				}
			}
		}

		if low[u] == dfn[u] {
			for {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[v] = false
				scc[v] = sccCount
				if u == v {
					break
				}
			}
			sccCount++
		}
	}

	for i := 0; i < numNodes; i++ {
		if dfn[i] == -1 {
			tarjan(i)
		}
	}

	for _, e := range strictEdges {
		if scc[e.u] == scc[e.v] {
			fmt.Println("No solution")
			return
		}
	}

	rows := make([]int, n)
	for i := 0; i < n; i++ {
		rows[i] = i + 1
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return scc[rows[i]-1] > scc[rows[j]-1]
	})

	cols := make([]int, m)
	for i := 0; i < m; i++ {
		cols[i] = i + 1
	}
	sort.SliceStable(cols, func(i, j int) bool {
		return scc[n+cols[i]-1] > scc[n+cols[j]-1]
	})

	var out []string
	for _, r := range rows {
		out = append(out, strconv.Itoa(r))
	}
	fmt.Println(strings.Join(out, " "))

	out = []string{}
	for _, c := range cols {
		out = append(out, strconv.Itoa(c))
	}
	fmt.Println(strings.Join(out, " "))
}