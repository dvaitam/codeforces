package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type test struct {
	input       string
	expectedYes []bool
	cases       []caseData
}

type caseData struct {
	n     int
	edges [][2]int
}

func prepareBinary(path, tag string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	bin := filepath.Join(os.TempDir(), tag)
	cmd := exec.Command("go", "build", "-o", bin, path)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", func() {}, fmt.Errorf("build %s: %v\n%s", path, err, out)
	}
	cleanup := func() { os.Remove(bin) }
	return bin, cleanup, nil
}

func runBinary(path, input string, timeout time.Duration) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	timer := time.AfterFunc(timeout, func() { cmd.Process.Kill() })
	err := cmd.Run()
	timer.Stop()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatInput(cases []caseData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for _, e := range cs.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func parseFeasible(refOut string, cases []caseData) ([]bool, error) {
	tokens := strings.Fields(refOut)
	idx := 0
	res := make([]bool, len(cases))
	for i, cs := range cases {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("reference output ended early on case %d", i+1)
		}
		word := strings.ToLower(tokens[idx])
		idx++
		if word[0] == 'n' {
			res[i] = false
			continue
		}
		if word[0] != 'y' {
			return nil, fmt.Errorf("invalid YES/NO token %q on case %d", word, i+1)
		}
		res[i] = true
		need := (cs.n - 1) * 2
		if idx+need > len(tokens) {
			return nil, fmt.Errorf("reference output missing edges for case %d", i+1)
		}
		idx += need
	}
	return res, nil
}

func verifyCase(cs caseData, wantYes bool, tokens []string, pos *int) error {
	if *pos >= len(tokens) {
		return fmt.Errorf("output ended early for case with n=%d", cs.n)
	}
	word := strings.ToLower(tokens[*pos])
	*pos++
	if word[0] == 'n' {
		if wantYes {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	}
	if word[0] != 'y' {
		return fmt.Errorf("invalid YES/NO token %q", word)
	}
	if !wantYes {
		return fmt.Errorf("expected NO but got YES")
	}

	if (*pos)+2*(cs.n-1) > len(tokens) {
		return fmt.Errorf("not enough edges provided")
	}

	edgeMap := make(map[int]map[int]bool)
	for _, e := range cs.edges {
		u, v := e[0], e[1]
		if edgeMap[u] == nil {
			edgeMap[u] = make(map[int]bool)
		}
		if edgeMap[v] == nil {
			edgeMap[v] = make(map[int]bool)
		}
		edgeMap[u][v] = true
		edgeMap[v][u] = true
	}

	outAdj := make([][]int, cs.n+1)
	used := make(map[[2]int]bool)
	for i := 0; i < cs.n-1; i++ {
		u := atoi(tokens[*pos])
		v := atoi(tokens[*pos+1])
		*pos += 2
		if u < 1 || u > cs.n || v < 1 || v > cs.n {
			return fmt.Errorf("edge vertices out of range")
		}
		if !edgeMap[u][v] {
			return fmt.Errorf("edge %d-%d not in original tree", u, v)
		}
		key := [2]int{u, v}
		rev := [2]int{v, u}
		if used[key] || used[rev] {
			return fmt.Errorf("edge %d-%d repeated", u, v)
		}
		used[key] = true
		outAdj[u] = append(outAdj[u], v)
	}
	if len(used) != cs.n-1 {
		return fmt.Errorf("not all edges used")
	}

	// detect cycles and build topological order
	color := make([]int, cs.n+1)
	order := make([]int, 0, cs.n)
	var dfs func(int) bool
	dfs = func(u int) bool {
		color[u] = 1
		for _, v := range outAdj[u] {
			if color[v] == 1 {
				return false
			}
			if color[v] == 0 && !dfs(v) {
				return false
			}
		}
		color[u] = 2
		order = append(order, u)
		return true
	}
	for i := 1; i <= cs.n; i++ {
		if color[i] == 0 {
			if !dfs(i) {
				return fmt.Errorf("cycle detected in candidate orientation")
			}
		}
	}

	reach := make([]int64, cs.n+1)
	for _, u := range order {
		reach[u] = 1
		for _, v := range outAdj[u] {
			reach[u] += reach[v]
		}
	}
	var good int64
	for i := 1; i <= cs.n; i++ {
		good += reach[i] - 1
	}
	if good != int64(cs.n) {
		return fmt.Errorf("good pair count %d != %d", good, cs.n)
	}
	return nil
}

func atoi(s string) int {
	var v int
	for i := 0; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v
}

func generateTree(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return edges
}

func generateTests() []test {
	var tests []test

	// small fixed cases
	tests = append(tests, test{
		cases: []caseData{
			{n: 2, edges: [][2]int{{1, 2}}},
		},
	})
	tests = append(tests, test{
		cases: []caseData{
			{n: 3, edges: [][2]int{{1, 2}, {2, 3}}},
		},
	})

	rng := rand.New(rand.NewSource(2112))
	for i := 0; i < 40; i++ {
		tcCount := rng.Intn(3) + 1
		cases := make([]caseData, tcCount)
		for j := 0; j < tcCount; j++ {
			n := rng.Intn(40) + 2
			cases[j] = caseData{n: n, edges: generateTree(n, rng)}
		}
		tests = append(tests, test{cases: cases})
	}

	// bigger stress
	cases := []caseData{
		{n: 200, edges: generateTree(200, rand.New(rand.NewSource(7)))},
		{n: 201, edges: generateTree(201, rand.New(rand.NewSource(8)))},
	}
	tests = append(tests, test{cases: cases})

	for i := range tests {
		tests[i].input = formatInput(tests[i].cases)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	candBin, candCleanup, err := prepareBinary(os.Args[1], "cand2112D")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "2112D.go")
	refBin, refCleanup, err := prepareBinary(refPath, "ref2112D")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runBinary(refBin, tc.input, 5*time.Second)
		if err != nil {
			fmt.Printf("Reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectedYes, err := parseFeasible(refOut, tc.cases)
		if err != nil {
			fmt.Printf("Reference parse error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		gotOut, err := runBinary(candBin, tc.input, 5*time.Second)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		tokens := strings.Fields(gotOut)
		pos := 0
		for idx, cs := range tc.cases {
			if err := verifyCase(cs, expectedYes[idx], tokens, &pos); err != nil {
				fmt.Printf("Wrong answer on test %d case %d: %v\nInput:\n%s\n", i+1, idx+1, err, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
