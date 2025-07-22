package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

func generateBus(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(6) + 2 // 2..7
	edges := make([]edge, n-1)
	for i := 1; i < n; i++ {
		edges[i-1] = edge{i, i + 1}
	}
	return n, edges
}

func generateRing(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(6) + 3 // 3..8
	edges := make([]edge, n)
	for i := 1; i < n; i++ {
		edges[i-1] = edge{i, i + 1}
	}
	edges[n-1] = edge{n, 1}
	return n, edges
}

func generateStar(rng *rand.Rand) (int, []edge) {
	n := rng.Intn(6) + 3 // 3..8
	edges := make([]edge, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = edge{1, i}
	}
	return n, edges
}

func generateUnknown(rng *rand.Rand) (int, []edge) {
	// start from one topology and add extra edge
	typ := rng.Intn(3)
	var n int
	var e []edge
	switch typ {
	case 0:
		n, e = generateBus(rng)
	case 1:
		n, e = generateRing(rng)
	default:
		n, e = generateStar(rng)
	}
	// add extra edge to break topology
	// choose random pair not already present
	exist := make(map[[2]int]bool)
	for _, ed := range e {
		a, b := ed.u, ed.v
		if a > b {
			a, b = b, a
		}
		exist[[2]int{a, b}] = true
	}
	for {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		if !exist[[2]int{a, b}] {
			e = append(e, edge{u, v})
			break
		}
	}
	return n, e
}

func topology(n int, edges []edge) string {
	m := len(edges)
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
	}
	cnt1, cnt2, cntn1 := 0, 0, 0
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			cnt1++
		}
		if deg[i] == 2 {
			cnt2++
		}
		if deg[i] == n-1 {
			cntn1++
		}
	}
	switch {
	case m == n-1 && cnt1 == 2 && cnt2 == n-2:
		return "bus topology"
	case m == n && cnt2 == n:
		return "ring topology"
	case m == n-1 && cntn1 == 1 && cnt1 == n-1:
		return "star topology"
	default:
		return "unknown topology"
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	typ := rng.Intn(4)
	var n int
	var e []edge
	var expected string
	switch typ {
	case 0:
		n, e = generateBus(rng)
	case 1:
		n, e = generateRing(rng)
	case 2:
		n, e = generateStar(rng)
	default:
		n, e = generateUnknown(rng)
	}
	expected = topology(n, e)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(e))
	for _, ed := range e {
		fmt.Fprintf(&sb, "%d %d\n", ed.u, ed.v)
	}
	return sb.String(), expected
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
