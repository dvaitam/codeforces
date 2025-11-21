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

type testCase struct {
	n, m int
	p    []int
	a    [][]int64
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierH.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)

	for idx, tc := range tests {
		input := buildInput(tc)

		oracleOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, oracleOut)
			os.Exit(1)
		}
		candOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}

		bestSeq, impossibleOracle, err := parseSolution(oracleOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, oracleOut)
			os.Exit(1)
		}
		candSeq, impossibleCand, err := parseSolution(candOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			os.Exit(1)
		}

		if impossibleOracle {
			if !impossibleCand {
				fmt.Fprintf(os.Stderr, "test %d: oracle says impossible but candidate provided a plan\ninput:\n%s", idx+1, input)
				os.Exit(1)
			}
			continue
		}
		if impossibleCand {
			fmt.Fprintf(os.Stderr, "test %d: candidate answered -1 but a solution exists\ninput:\n%s", idx+1, input)
			os.Exit(1)
		}

		bestScore, err := simulate(tc, bestSeq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle plan invalid on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		candScore, err := simulate(tc, candSeq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid candidate plan: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if candScore != bestScore {
			fmt.Fprintf(os.Stderr, "test %d: suboptimal score (got %d, expected %d)\ninput:\n%s", idx+1, candScore, bestScore, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2038H-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleH")
	cmd := exec.Command("go", "build", "-o", outPath, "2038H.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
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
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.a[i][j], 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2, m: 3,
			p: []int{2, 1, 2},
			a: [][]int64{
				{1, 2, 3},
				{4, 5, 6},
			},
		},
		{
			n: 3, m: 5,
			p: []int{1, 1, 1, 2, 1},
			a: [][]int64{
				{1, 1, 1, 1, 1},
				{10, 5, 7, 8, 15},
				{7, 10, 9, 8, 15},
			},
		},
		{
			n: 3, m: 5,
			p: []int{1, 1, 1, 1, 1},
			a: [][]int64{
				{1, 1, 1, 1, 1},
				{10, 5, 7, 8, 15},
				{7, 10, 9, 8, 15},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	for len(tests) < cap(tests) {
		n := rng.Intn(5) + 2 // 2..6
		m := rng.Intn(5) + 2 // 2..6
		tc := testCase{
			n: n,
			m: m,
			p: make([]int, m),
			a: make([][]int64, n),
		}
		for i := 0; i < m; i++ {
			tc.p[i] = rng.Intn(n) + 1
		}
		for i := 0; i < n; i++ {
			tc.a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				tc.a[i][j] = int64(rng.Intn(50) + 1)
			}
		}
		tests = append(tests, tc)
	}
	return tests
}

func parseSolution(out string, moves int) ([]int, bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return nil, false, fmt.Errorf("empty output")
	}
	if tokens[0] == "-1" {
		if len(tokens) > 1 {
			return nil, false, fmt.Errorf("extra tokens after -1")
		}
		return nil, true, nil
	}
	if len(tokens) != moves {
		return nil, false, fmt.Errorf("expected %d integers, got %d", moves, len(tokens))
	}
	seq := make([]int, moves)
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, false, fmt.Errorf("invalid integer %q", tok)
		}
		seq[i] = val
	}
	return seq, false, nil
}

func simulate(tc testCase, seq []int) (int64, error) {
	n := tc.n
	m := tc.m
	if len(seq) != m {
		return 0, fmt.Errorf("expected %d moves, got %d", m, len(seq))
	}
	counts := make([]int, n)
	ruling := -1
	score := int64(0)

	for turn := 0; turn < m; turn++ {
		choice := seq[turn] - 1
		if choice < 0 || choice >= n {
			return 0, fmt.Errorf("move %d: party %d outside [1,%d]", turn+1, seq[turn], n)
		}
		if ruling != -1 && choice == ruling {
			return 0, fmt.Errorf("move %d: supported current ruling party %d", turn+1, seq[turn])
		}
		counts[choice]++
		score += tc.a[choice][turn]

		newRuling := determineLeader(counts)
		if newRuling != tc.p[turn]-1 {
			return 0, fmt.Errorf("turn %d: ruling party became %d but required %d", turn+1, newRuling+1, tc.p[turn])
		}
		ruling = newRuling
	}

	return score, nil
}

func determineLeader(counts []int) int {
	lead := 0
	mx := counts[0]
	for i := 1; i < len(counts); i++ {
		if counts[i] > mx || (counts[i] == mx && i < lead) {
			mx = counts[i]
			lead = i
		}
	}
	return lead
}
