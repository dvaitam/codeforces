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

type edge struct{ u, v int }

func expectedColor(n int, colors []int, edges []edge) int {
	neigh := make(map[int]map[int]struct{})
	present := make(map[int]struct{})
	for _, c := range colors {
		present[c] = struct{}{}
	}
	for _, e := range edges {
		cu := colors[e.u-1]
		cv := colors[e.v-1]
		if cu == cv {
			continue
		}
		if neigh[cu] == nil {
			neigh[cu] = make(map[int]struct{})
		}
		if neigh[cv] == nil {
			neigh[cv] = make(map[int]struct{})
		}
		neigh[cu][cv] = struct{}{}
		neigh[cv][cu] = struct{}{}
	}
	bestColor := 0
	bestCount := -1
	colorsList := make([]int, 0, len(present))
	for c := range present {
		colorsList = append(colorsList, c)
	}
	sort.Ints(colorsList)
	for _, c := range colorsList {
		cnt := len(neigh[c])
		if cnt > bestCount {
			bestCount = cnt
			bestColor = c
		}
	}
	return bestColor
}

func runCase(bin string, n int, colors []int, edges []edge) error {
	var sb strings.Builder
	m := len(edges)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(c))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("non-integer output %q", gotStr)
	}
	expect := expectedColor(n, colors, edges)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		colors := make([]int, n)
		for j := range colors {
			colors[j] = rng.Intn(5) + 1
		}
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges + 1)
		used := make(map[[2]int]struct{})
		edges := make([]edge, 0, m)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			edges = append(edges, edge{u, v})
		}
		if err := runCase(bin, n, colors, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
