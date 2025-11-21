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
	name  string
	input string
	count int
}

type singleCase struct {
	swords   []int64
	monsters []int64
	rewards  []int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2164C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2164C.go")
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
	return strings.TrimSpace(stdout.String()), nil
}

func formatInput(cases []singleCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, cs := range cases {
		n := len(cs.swords)
		m := len(cs.monsters)
		if len(cs.rewards) != m {
			panic("rewards length mismatch")
		}
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(m))
		sb.WriteByte('\n')
		for i, v := range cs.swords {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range cs.monsters {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range cs.rewards {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	const limit int64 = 1_000_000_000
	return []testCase{
		{
			name: "single_can_kill",
			input: formatInput([]singleCase{
				{
					swords:   []int64{5},
					monsters: []int64{4},
					rewards:  []int64{0},
				},
			}),
			count: 1,
		},
		{
			name: "single_cannot_kill",
			input: formatInput([]singleCase{
				{
					swords:   []int64{3},
					monsters: []int64{5},
					rewards:  []int64{0},
				},
			}),
			count: 1,
		},
		{
			name: "upgrade_sequence",
			input: formatInput([]singleCase{
				{
					swords:   []int64{4},
					monsters: []int64{4, 8},
					rewards:  []int64{6, 0},
				},
				{
					swords:   []int64{2, 5},
					monsters: []int64{3, 6, 7},
					rewards:  []int64{0, 10, 0},
				},
				{
					swords:   []int64{limit},
					monsters: []int64{limit, limit},
					rewards:  []int64{limit, limit},
				},
			}),
			count: 3,
		},
		{
			name: "no_rewards",
			input: formatInput([]singleCase{
				{
					swords:   []int64{1, 2, 3},
					monsters: []int64{2, 3, 4, 5},
					rewards:  []int64{0, 0, 0, 0},
				},
			}),
			count: 1,
		},
	}
}

func randomCase(rng *rand.Rand, maxN, maxM int, maxVal int64) singleCase {
	n := rng.Intn(maxN) + 1
	m := rng.Intn(maxM) + 1
	swords := make([]int64, n)
	for i := 0; i < n; i++ {
		swords[i] = rng.Int63n(maxVal) + 1
	}
	monsters := make([]int64, m)
	rewards := make([]int64, m)
	for i := 0; i < m; i++ {
		monsters[i] = rng.Int63n(maxVal) + 1
		if rng.Intn(3) == 0 {
			rewards[i] = 0
		} else {
			rewards[i] = rng.Int63n(maxVal) + 1
		}
	}
	return singleCase{swords: swords, monsters: monsters, rewards: rewards}
}

func generateTests() []testCase {
	tests := deterministicTests()
	const maxVal int64 = 1_000_000_000
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 40; i++ {
		t := rng.Intn(5) + 1
		cases := make([]singleCase, 0, t)
		sumN := 0
		sumM := 0
		for len(cases) < t {
			remainingN := 200000 - sumN
			remainingM := 200000 - sumM
			if remainingN <= 0 || remainingM <= 0 {
				break
			}
			maxN := remainingN
			if maxN > 2000 {
				maxN = 2000
			}
			maxM := remainingM
			if maxM > 2000 {
				maxM = 2000
			}
			cs := randomCase(rng, maxN, maxM, maxVal)
			sumN += len(cs.swords)
			sumM += len(cs.monsters)
			cases = append(cases, cs)
		}
		if len(cases) == 0 {
			continue
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_batch_%d", i+1),
			input: formatInput(cases),
			count: len(cases),
		})
	}

	// Stress test hitting constraint limits.
	maxCase := singleCase{
		swords:   make([]int64, 200000),
		monsters: make([]int64, 200000),
		rewards:  make([]int64, 200000),
	}
	for i := range maxCase.swords {
		maxCase.swords[i] = int64((i % 1_000) + 1)
	}
	for i := range maxCase.monsters {
		maxCase.monsters[i] = int64((i % 1_000) + 1)
		if i%4 == 0 {
			maxCase.rewards[i] = 0
		} else {
			val := int64((i*17)%1_000_000_000 + 1)
			maxCase.rewards[i] = val
		}
	}
	tests = append(tests, testCase{
		name:  "max_constraints",
		input: formatInput([]singleCase{maxCase}),
		count: 1,
	})

	return tests
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d answers, got %d", count, len(act))
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at case %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expected, err := runBinary(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		actual, err := runBinary(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual, tc.count); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, tc.name, err, tc.input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
