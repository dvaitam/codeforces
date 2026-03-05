package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const totalTests = 80

type testCase struct {
	n     int
	edges [][2]int
}

// solveReference implements the C++ algorithm inline:
//  1. DFS from node 1 to find one diameter endpoint r.
//  2. DFS2 from r: for each node compute max-child-height t,
//     do c[t]--, c[t+1]++, return t+1.
//     (The missing return in the C++ source returns t+1 in practice.)
//  3. Sweep c[] from n down to 1 to build the answer array.
func solveReference(tc testCase) []int {
	n := tc.n
	adj := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// Step 1: find one endpoint of the diameter.
	dist := make([]int, n+1)
	r := 1
	var dfs1 func(fa, u int)
	dfs1 = func(fa, u int) {
		if dist[r] < dist[u] {
			r = u
		}
		for _, v := range adj[u] {
			if v != fa {
				dist[v] = dist[u] + 1
				dfs1(u, v)
			}
		}
	}
	dfs1(0, 1)

	// Step 2: DFS from r to fill c[].
	cArr := make([]int, n+2)
	var dfs2 func(fa, u int) int
	dfs2 = func(fa, u int) int {
		t := 0
		for _, v := range adj[u] {
			if v != fa {
				if h := dfs2(u, v); h > t {
					t = h
				}
			}
		}
		cArr[t]--
		cArr[t+1]++
		return t + 1
	}
	dfs2(0, r)

	// Step 3: build answer array.
	result := make([]int, n)
	result[0] = 1
	rr := 0
	s := 0
	for i := n; i >= 1; i-- {
		for j := 0; j < cArr[i]; j++ {
			s += i
			rr++
			result[rr] = s
		}
	}
	for i := rr + 1; i < n; i++ {
		result[i] = n
	}
	return result
}

func prepareProgram(path string) (string, func(), error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".go" && ext != ".cpp" && ext != ".cc" && ext != ".cxx" {
		return path, nil, nil
	}
	dir, err := os.MkdirTemp("", "verifier958B2-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "cand")
	var cmd *exec.Cmd
	if ext == ".go" {
		cmd = exec.Command("go", "build", "-o", bin, path)
	} else {
		// The canonical solution for this problem has a missing return statement
		// in dfs2 (undefined behaviour). On Codeforces' GCC 7.3 -O2, t+1 happens
		// to remain in eax as the implicit return value. On newer GCC the garbage
		// return causes an out-of-bounds access and a segfault.
		// We patch the source to add the explicit return before compiling.
		// Use a regex so variations in whitespace (or statements on separate lines)
		// are handled correctly; error if the pattern is not found at all.
		src, err := os.ReadFile(path)
		if err != nil {
			os.RemoveAll(dir)
			return "", nil, err
		}
		dfs2Re := regexp.MustCompile(`c\[t\]--;\s*c\[t\+1\]\+\+;`)
		original := string(src)
		patched := dfs2Re.ReplaceAllString(original, "${0}return t+1;")
		if patched == original {
			os.RemoveAll(dir)
			return "", nil, fmt.Errorf("cannot patch dfs2 missing return: pattern 'c[t]--;...c[t+1]++;' not found in %s", path)
		}
		patchedSrc := filepath.Join(dir, "src.cpp")
		if err := os.WriteFile(patchedSrc, []byte(patched), 0644); err != nil {
			os.RemoveAll(dir)
			return "", nil, err
		}
		cmd = exec.Command("g++", "-O2", "-std=c++17", "-o", bin, patchedSrc)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("compile error: %v\n%s", err, out.String())
	}
	if runtime.GOOS == "windows" {
		if _, err2 := os.Stat(bin); err2 != nil {
			if _, e3 := os.Stat(bin + ".exe"); e3 == nil {
				bin += ".exe"
			}
		}
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return []byte(sb.String())
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	res := make([]int, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildStar(3),
		buildLine(5),
		buildLine(10),
		buildRandom(20, rand.New(rand.NewSource(1))),
		buildRandom(40, rand.New(rand.NewSource(2))),
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-5 {
		n := rnd.Intn(80) + 2
		tests = append(tests, buildRandom(n, rnd))
	}
	tests = append(tests,
		buildRandom(200, rand.New(rand.NewSource(3))),
		buildRandom(500, rand.New(rand.NewSource(4))),
		buildRandom(1000, rand.New(rand.NewSource(5))),
		buildRandom(5000, rand.New(rand.NewSource(6))),
		buildRandom(10000, rand.New(rand.NewSource(7))),
	)
	return tests
}

func buildStar(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return testCase{n: n, edges: edges}
}

func buildLine(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testCase{n: n, edges: edges}
}

func buildRandom(n int, rnd *rand.Rand) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rnd.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return testCase{n: n, edges: edges}
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}

	candidate, cleanup, err := prepareProgram(os.Args[1])
	if err != nil {
		fmt.Println("failed to prepare candidate:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refVals := solveReference(tc)

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d values, got %d\n", idx+1, len(refVals), len(candVals))
			printInput(input)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Printf("test %d failed at position %d: expected %d, got %d\n", idx+1, i+1, refVals[i], candVals[i])
				printInput(input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
