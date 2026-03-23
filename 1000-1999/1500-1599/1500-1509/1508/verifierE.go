package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// Embedded correct solver for 1508E
func solve1508E(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var out strings.Builder
	writer := bufio.NewWriter(&out)

	var n int
	fmt.Fscan(reader, &n)

	a := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		pos[a[i]] = i
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
	}

	minVal := make([]int, n+1)
	var dfsMin func(int)
	dfsMin = func(u int) {
		minVal[u] = a[u]
		for _, v := range adj[u] {
			dfsMin(v)
			if minVal[v] < minVal[u] {
				minVal[u] = minVal[v]
			}
		}
	}
	dfsMin(1)

	for u := 1; u <= n; u++ {
		sort.Slice(adj[u], func(i, j int) bool {
			return minVal[adj[u][i]] < minVal[adj[u][j]]
		})
	}

	aOrig := make([]int, n+1)
	origPos := make([]int, n+1)
	timer := 1
	var buildOrig func(int)
	buildOrig = func(u int) {
		aOrig[u] = timer
		origPos[timer] = u
		timer++
		for _, v := range adj[u] {
			buildOrig(v)
		}
	}
	buildOrig(1)

	bit := make([]int, n+2)
	add := func(idx, val int) {
		for ; idx <= n; idx += idx & -idx {
			bit[idx] += val
		}
	}
	query := func(idx int) int {
		sum := 0
		for ; idx > 0; idx -= idx & -idx {
			sum += bit[idx]
		}
		return sum
	}

	var D int64
	var countInversions func(int)
	countInversions = func(u int) {
		greater := query(n) - query(a[u])
		D += int64(greater)
		add(a[u], 1)
		for _, v := range adj[u] {
			countInversions(v)
		}
		add(a[u], -1)
	}
	countInversions(1)

	maxX := 0
	for u := 1; u <= n; u++ {
		if a[u] < aOrig[u] {
			if a[u] > maxX {
				maxX = a[u]
			}
		}
	}

	for y := 1; y < maxX; y++ {
		u := pos[y]
		for _, v := range adj[u] {
			if a[v] > y {
				fmt.Fprintln(writer, "NO")
				writer.Flush()
				return strings.TrimSpace(out.String())
			}
		}
	}

	for u := 1; u <= n; u++ {
		actualMin := a[u]
		var checkMin func(int)
		checkMin = func(curr int) {
			if a[curr] < actualMin {
				actualMin = a[curr]
			}
			for _, child := range adj[curr] {
				checkMin(child)
			}
		}
		checkMin(u)
		if actualMin != minVal[u] {
			fmt.Fprintln(writer, "NO")
			writer.Flush()
			return strings.TrimSpace(out.String())
		}
	}

	fmt.Fprintln(writer, "YES")
	fmt.Fprintln(writer, D)
	for i := 1; i <= n; i++ {
		fmt.Fprint(writer, aOrig[i])
		if i < n {
			fmt.Fprint(writer, " ")
		}
	}
	fmt.Fprintln(writer)
	writer.Flush()
	return strings.TrimSpace(out.String())
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2 // 2..6
	perm := r.Perm(n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", perm[i]+1)
	}
	sb.WriteByte('\n')
	for i := 2; i <= n; i++ {
		parent := r.Intn(i-1) + 1
		fmt.Fprintf(&sb, "%d %d\n", parent, i)
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := solve1508E(input)
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
