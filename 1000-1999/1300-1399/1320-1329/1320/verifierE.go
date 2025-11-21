package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type virus struct {
	city  int
	speed int
}

type scenario struct {
	viruses   []virus
	important []int
}

type testCase struct {
	n         int
	edges     [][2]int
	scenarios []scenario
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1320E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "1320E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*6 + len(tc.scenarios)*32)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
		sb.WriteByte('\n')
	}
	sb.WriteString(strconv.Itoa(len(tc.scenarios)))
	sb.WriteByte('\n')
	for _, sc := range tc.scenarios {
		sb.WriteString(strconv.Itoa(len(sc.viruses)))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(len(sc.important)))
		sb.WriteByte('\n')
		for _, v := range sc.viruses {
			sb.WriteString(strconv.Itoa(v.city))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v.speed))
			sb.WriteByte('\n')
		}
		for i, city := range sc.important {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(city))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tc testCase) ([][]int, error) {
	tokens := strings.Fields(out)
	totalNeeded := 0
	for _, sc := range tc.scenarios {
		totalNeeded += len(sc.important)
	}
	if len(tokens) != totalNeeded {
		return nil, fmt.Errorf("expected %d integers, got %d", totalNeeded, len(tokens))
	}
	res := make([][]int, len(tc.scenarios))
	idx := 0
	for i, sc := range tc.scenarios {
		res[i] = make([]int, len(sc.important))
		for j := range sc.important {
			val, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %v", tokens[idx], err)
			}
			res[i][j] = val
			idx++
		}
	}
	return res, nil
}

func deterministicTests() []testCase {
	sample := testCase{
		n: 7,
		edges: [][2]int{
			{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {3, 7},
		},
		scenarios: []scenario{
			{
				viruses:   []virus{{city: 4, speed: 1}, {city: 7, speed: 1}},
				important: []int{1, 3},
			},
			{
				viruses:   []virus{{city: 4, speed: 3}, {city: 7, speed: 1}},
				important: []int{1, 3},
			},
			{
				viruses:   []virus{{city: 1, speed: 1}, {city: 4, speed: 100}, {city: 7, speed: 100}},
				important: []int{1, 2, 3},
			},
		},
	}

	singleNode := testCase{
		n:     1,
		edges: nil,
		scenarios: []scenario{
			{
				viruses:   []virus{{city: 1, speed: 5}},
				important: []int{1},
			},
		},
	}

	line := testCase{
		n:     5,
		edges: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
		scenarios: []scenario{
			{
				viruses:   []virus{{city: 1, speed: 1}, {city: 5, speed: 1}},
				important: []int{2, 3, 4},
			},
			{
				viruses:   []virus{{city: 2, speed: 2}},
				important: []int{5, 4, 3, 2, 1},
			},
		},
	}

	return []testCase{sample, singleNode, line}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 80)
	for len(tests) < cap(tests) {
		n := rng.Intn(80) + 1
		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			u := rng.Intn(v-1) + 1
			edges = append(edges, [2]int{u, v})
		}
		q := rng.Intn(8) + 1
		if rng.Intn(5) == 0 {
			q = rng.Intn(15) + 1
		}
		scenarios := make([]scenario, q)
		for i := 0; i < q; i++ {
			maxK := n
			if maxK > 12 {
				maxK = 12
			}
			k := rng.Intn(maxK) + 1
			if k > n {
				k = n
			}
			if rng.Intn(5) == 0 {
				k = 1
			}
			nodes := rng.Perm(n)
			viruses := make([]virus, k)
			for j := 0; j < k; j++ {
				speed := rng.Intn(1_000_000) + 1
				if rng.Intn(4) == 0 {
					speed = rng.Intn(20) + 1
				}
				viruses[j] = virus{city: nodes[j] + 1, speed: speed}
			}
			maxM := n
			if maxM > 20 {
				maxM = 20
			}
			m := rng.Intn(maxM) + 1
			if rng.Intn(6) == 0 {
				m = n
			}
			impPick := rng.Perm(n)
			important := make([]int, m)
			for j := 0; j < m; j++ {
				important[j] = impPick[j] + 1
			}
			scenarios[i] = scenario{viruses: viruses, important: important}
		}
		tests = append(tests, testCase{n: n, edges: edges, scenarios: scenarios})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedAns, err := parseOutput(expectedOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualAns, err := parseOutput(actualOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		for sIdx := range tc.scenarios {
			expVals := expectedAns[sIdx]
			actVals := actualAns[sIdx]
			for j := range expVals {
				if expVals[j] != actVals[j] {
					fmt.Fprintf(os.Stderr, "test %d scenario %d city %d mismatch: expected %d, got %d\ninput:\n%s", idx+1, sIdx+1, j+1, expVals[j], actVals[j], input)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Println("All tests passed.")
}
