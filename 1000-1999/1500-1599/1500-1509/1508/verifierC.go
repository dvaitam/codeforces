package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Brute-force solver for 1508C (small n only).
// Complete graph: assigned edges have given weights, unassigned edges must be
// assigned non-negative weights so that XOR of ALL edge weights = 0, and the
// MST weight is minimised.
// Strategy: XOR of assigned weights = X.  If X=0, set all unassigned to 0.
// Otherwise, try putting X on each single unassigned edge (rest 0), compute MST.
func solve1508C(input string) string {
	data := []byte(input)
	ptr := 0
	nextInt := func() int {
		for ptr < len(data) && (data[ptr] < '0' || data[ptr] > '9') {
			ptr++
		}
		val := 0
		for ptr < len(data) && data[ptr] >= '0' && data[ptr] <= '9' {
			val = val*10 + int(data[ptr]-'0')
			ptr++
		}
		return val
	}

	n := nextInt()
	m := nextInt()

	type Edge struct {
		u, v, w int
	}

	assigned := make(map[[2]int]int) // key: (min,max) -> weight
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		w := nextInt()
		if u > v {
			u, v = v, u
		}
		assigned[[2]int{u, v}] = w
	}

	// Compute XOR of all assigned weights
	xorAll := 0
	for _, w := range assigned {
		xorAll ^= w
	}

	// Collect all edges of the complete graph
	type WEdge struct {
		u, v, w int
	}

	// Kruskal MST helper
	kruskalMST := func(allEdges []WEdge) int64 {
		sort.Slice(allEdges, func(i, j int) bool {
			return allEdges[i].w < allEdges[j].w
		})
		par := make([]int, n+1)
		for i := range par {
			par[i] = i
		}
		var find func(int) int
		find = func(x int) int {
			if par[x] != x {
				par[x] = find(par[x])
			}
			return par[x]
		}
		total := int64(0)
		cnt := 0
		for _, e := range allEdges {
			a, b := find(e.u), find(e.v)
			if a != b {
				par[a] = b
				total += int64(e.w)
				cnt++
				if cnt == n-1 {
					break
				}
			}
		}
		return total
	}

	// Build edge list for a given assignment scenario
	buildEdges := func(extraEdge [2]int, extraW int) []WEdge {
		var edges []WEdge
		for u := 1; u <= n; u++ {
			for v := u + 1; v <= n; v++ {
				key := [2]int{u, v}
				if w, ok := assigned[key]; ok {
					edges = append(edges, WEdge{u, v, w})
				} else if key == extraEdge {
					edges = append(edges, WEdge{u, v, extraW})
				} else {
					edges = append(edges, WEdge{u, v, 0})
				}
			}
		}
		return edges
	}

	best := int64(1<<62)

	if xorAll == 0 {
		// All unassigned edges get weight 0
		edges := buildEdges([2]int{0, 0}, 0)
		best = kruskalMST(edges)
	} else {
		// Try assigning xorAll to each unassigned edge
		for u := 1; u <= n; u++ {
			for v := u + 1; v <= n; v++ {
				key := [2]int{u, v}
				if _, ok := assigned[key]; ok {
					continue
				}
				edges := buildEdges(key, xorAll)
				cost := kruskalMST(edges)
				if cost < best {
					best = cost
				}
			}
		}
	}

	return strconv.FormatInt(best, 10)
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges)
	if m == maxEdges {
		m--
	}
	edgesMap := make(map[[2]int]struct{})
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for len(edgesMap) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := edgesMap[key]; ok {
			continue
		}
		edgesMap[key] = struct{}{}
		w := r.Intn(1000) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := solve1508C(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
