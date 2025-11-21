package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSourceC1  = "0-999/100-199/170-179/177/177C1.go"
	randomTrials = 200
)

type pair struct {
	u int
	v int
}

type genDSU struct {
	parent []int
	size   []int
}

func newGenDSU(n int) *genDSU {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &genDSU{parent: parent, size: size}
}

func (d *genDSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *genDSU) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceC1)
	if err != nil {
		fatal("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	fixed := deterministicCases()
	tests := append([]string{}, fixed...)
	tests = append(tests, makeHeavyCase(2000, 100000, 0))
	tests = append(tests, makeHeavyCase(2000, 100000, 80000))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := len(tests) + randomTrials
	for len(tests) < target {
		tests = append(tests, randomCase(rng))
	}

	for idx, input := range tests {
		expect, err := runBinary(refBin, input)
		if err != nil {
			fatal("reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fatal("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		if expect != got {
			fatal("case %d mismatch\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "177C1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicCases() []string {
	result := make([]string, 0, 5)
	// Simple connected pair without dislikes.
	result = append(result, formatCase(2, []pair{{1, 2}}, nil))
	// Two components where the larger one is invalidated by a dislike.
	result = append(result, formatCase(4, []pair{{1, 2}, {2, 3}}, []pair{{1, 3}}))
	// Multiple components without dislikes.
	result = append(result, formatCase(5, []pair{{1, 2}, {2, 3}, {4, 5}}, nil))
	// Dislike across components should not matter.
	result = append(result, formatCase(6, []pair{{1, 2}, {3, 4}}, []pair{{2, 3}}))
	// No friendships but several dislikes.
	result = append(result, formatCase(6, nil, []pair{{1, 2}, {3, 4}, {5, 6}}))
	return result
}

func makeHeavyCase(n, k, m int) string {
	maxPairs := n * (n - 1) / 2
	if k > maxPairs {
		k = maxPairs
	}
	if m > maxPairs-k {
		m = maxPairs - k
	}
	used := make(map[pair]struct{}, k+m)
	friends := make([]pair, 0, k)
	for u := 1; u <= n && len(friends) < k; u++ {
		for v := u + 1; v <= n && len(friends) < k; v++ {
			p := pair{u: u, v: v}
			friends = append(friends, p)
			used[p] = struct{}{}
		}
	}
	dislikes := make([]pair, 0, m)
	for u := n; u >= 1 && len(dislikes) < m; u-- {
		for v := u - 1; v >= 1 && len(dislikes) < m; v-- {
			p := pair{u: v, v: u}
			if _, ok := used[p]; ok {
				continue
			}
			used[p] = struct{}{}
			dislikes = append(dislikes, p)
		}
	}
	return formatCase(n, friends, dislikes)
}

func randomCase(rng *rand.Rand) string {
	n := randomN(rng)
	maxPairs := n * (n - 1) / 2
	friendCap := maxPairs
	if n > 8 {
		friendCap = min(friendCap, 1500)
		if rng.Intn(5) == 0 {
			friendCap = min(friendCap, 20000)
		}
		if rng.Intn(12) == 0 {
			friendCap = min(friendCap, 100000)
		}
	}
	k := 0
	if friendCap > 0 {
		k = rng.Intn(friendCap + 1)
	}
	allUsed := make(map[pair]struct{})
	friends := make([]pair, 0, k)
	dsu := newGenDSU(n)
	for len(friends) < k {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := pair{u: u, v: v}
		if _, ok := allUsed[p]; ok {
			continue
		}
		allUsed[p] = struct{}{}
		friends = append(friends, p)
		dsu.union(u, v)
	}

	mCap := maxPairs - k
	if n > 8 {
		mCap = min(mCap, 1500)
		if rng.Intn(5) == 0 {
			mCap = min(mCap, 20000)
		}
		if rng.Intn(12) == 0 {
			mCap = min(mCap, 100000)
		}
	}
	m := 0
	if mCap > 0 {
		m = rng.Intn(mCap + 1)
	}
	dislikes := make([]pair, 0, m)
	if m > 0 && rng.Intn(3) == 0 {
		if p, ok := forceInternalDislike(n, dsu, allUsed, rng); ok {
			dislikes = append(dislikes, p)
		}
	}
	for len(dislikes) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := pair{u: u, v: v}
		if _, ok := allUsed[p]; ok {
			continue
		}
		allUsed[p] = struct{}{}
		dislikes = append(dislikes, p)
	}
	return formatCase(n, friends, dislikes)
}

func forceInternalDislike(n int, dsu *genDSU, used map[pair]struct{}, rng *rand.Rand) (pair, bool) {
	for attempt := 0; attempt < 1000; attempt++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if dsu.find(u) != dsu.find(v) {
			continue
		}
		if u > v {
			u, v = v, u
		}
		p := pair{u: u, v: v}
		if _, ok := used[p]; ok {
			continue
		}
		used[p] = struct{}{}
		return p, true
	}
	return pair{}, false
}

func formatCase(n int, friends, dislikes []pair) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	fmt.Fprintf(&sb, "%d\n", len(friends))
	for _, e := range friends {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	fmt.Fprintf(&sb, "%d\n", len(dislikes))
	for _, e := range dislikes {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func randomN(rng *rand.Rand) int {
	bucket := rng.Intn(100)
	switch {
	case bucket < 25:
		return 2 + rng.Intn(13) // 2..14
	case bucket < 60:
		return 15 + rng.Intn(185) // 15..199
	case bucket < 85:
		return 200 + rng.Intn(800) // 200..999
	default:
		return 1000 + rng.Intn(1001) // 1000..2000
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
