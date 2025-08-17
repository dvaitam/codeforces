package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// buildRef compiles the reference solution.
func buildRef() (string, error) {
	ref := "refL.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1252L.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return ref, nil
}

// runExe executes a binary with the provided input and returns stdout/stderr.
func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type TestCase string

// genTests creates random tests.
func genTests() []TestCase {
	rng := rand.New(rand.NewSource(12))
	tests := make([]TestCase, 0, 100)
	for i := 0; i < 100; i++ {
		N := rng.Intn(3) + 3 // 3..5
		K := rng.Intn(3) + 1 // 1..3 workers
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", N, K)

		perm := rng.Perm(N)
		A := make([]int, N)
		for j := 0; j < N; j++ {
			u := perm[j]
			v := perm[(j+1)%N]
			A[u] = v + 1
		}
		for j := 0; j < N; j++ {
			Ai := A[j]
			Mi := rng.Intn(3) + 1
			fmt.Fprintf(&sb, "%d %d ", Ai, Mi)
			vals := make(map[int]struct{})
			arr := make([]int, 0, Mi)
			for len(arr) < Mi {
				x := rng.Intn(5) + 1
				if _, ok := vals[x]; !ok {
					vals[x] = struct{}{}
					arr = append(arr, x)
				}
			}
			sort.Ints(arr)
			for _, x := range arr {
				fmt.Fprintf(&sb, "%d ", x)
			}
			sb.WriteByte('\n')
		}
		for j := 0; j < K; j++ {
			fmt.Fprintf(&sb, "%d ", rng.Intn(5)+1)
		}
		sb.WriteByte('\n')
		tests = append(tests, TestCase(sb.String()))
	}
	return tests
}

type instance struct {
	N, K    int
	workers []int64
	edgeIdx map[[2]int]int
	mats    []map[int64]struct{}
}

// parseInstance parses the input into a structured form for validation.
func parseInstance(input string) instance {
	r := strings.NewReader(input)
	var n, k int
	fmt.Fscan(r, &n, &k)
	edgeIdx := make(map[[2]int]int)
	mats := make([]map[int64]struct{}, n)
	for i := 0; i < n; i++ {
		var ai, mi int
		fmt.Fscan(r, &ai, &mi)
		ms := make(map[int64]struct{}, mi)
		for j := 0; j < mi; j++ {
			var x int64
			fmt.Fscan(r, &x)
			ms[x] = struct{}{}
		}
		mats[i] = ms
		u, v := i+1, ai
		if u > v {
			u, v = v, u
		}
		edgeIdx[[2]int{u, v}] = i
	}
	workers := make([]int64, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(r, &workers[i])
	}
	return instance{N: n, K: k, workers: workers, edgeIdx: edgeIdx, mats: mats}
}

type dsu struct {
	p []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{p: p}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a != b {
		d.p[a] = b
	}
}

// validate checks whether the candidate output is a valid solution for the instance.
func validate(inst instance, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	used := make(map[[2]int]bool)
	d := newDSU(inst.N + 1)
	for i := 0; i < inst.K; i++ {
		var u, v int
		if _, err := fmt.Fscan(reader, &u, &v); err != nil {
			return fmt.Errorf("output format error on worker %d: %v", i+1, err)
		}
		if u == 0 && v == 0 {
			continue
		}
		if u < 1 || u > inst.N || v < 1 || v > inst.N {
			return fmt.Errorf("invalid city in pair %d %d", u, v)
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		idx, ok := inst.edgeIdx[key]
		if !ok {
			return fmt.Errorf("edge %d %d not proposed", u, v)
		}
		if used[key] {
			return fmt.Errorf("edge %d %d used multiple times", u, v)
		}
		used[key] = true
		mat := inst.workers[i]
		if _, ok := inst.mats[idx][mat]; !ok {
			return fmt.Errorf("worker %d uses invalid material on edge %d %d", i+1, u, v)
		}
		d.union(u, v)
	}
	// ensure no extra tokens
	if _, err := fmt.Fscan(reader, new(string)); err == nil {
		return fmt.Errorf("extra data in output")
	}
	root := d.find(1)
	for v := 2; v <= inst.N; v++ {
		if d.find(v) != root {
			return fmt.Errorf("constructed roads do not connect all cities")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierL.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests := genTests()
	for i, t := range tests {
		input := string(t)
		inst := parseInstance(input)
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expTrim := strings.TrimSpace(exp)
		gotTrim := strings.TrimSpace(got)
		if expTrim == "-1" {
			if gotTrim != "-1" {
				fmt.Printf("Test %d failed\nInput:\n%sExpected:-1\nGot:%s\n", i+1, input, got)
				os.Exit(1)
			}
			continue
		}
		if gotTrim == "-1" {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:%sGot:-1\n", i+1, input, exp)
			os.Exit(1)
		}
		if err := validate(inst, got); err != nil {
			fmt.Printf("Test %d failed\nInput:\n%sError: %v\nCandidate Output:\n%s\n", i+1, input, err, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
