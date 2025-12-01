package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// refSource points to the local reference solution to avoid GOPATH resolution.
	refSource        = "2020C.go"
	maxA      uint64 = 1 << 61
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, input := range tests {
		refOut, err := runExecutable(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		refAnswers := strings.Fields(strings.TrimSpace(refOut))
		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candAnswers := strings.Fields(strings.TrimSpace(candOut))
		if len(refAnswers) != len(candAnswers) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: mismatched line count\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}

		tcases, cases := parseCases(input)
		for idx := 0; idx < tcases; idx++ {
			if idx >= len(refAnswers) || idx >= len(candAnswers) {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d: missing outputs\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
				os.Exit(1)
			}
			got := candAnswers[idx]
			caseData := cases[idx]
			if refAnswers[idx] == "-1" {
				if got != "-1" {
					if !validateSolution(caseData.b, caseData.c, caseData.d, got) {
						fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected no solution\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, idx+1, input, refOut, candOut)
						os.Exit(1)
					}
				}
			} else {
				if got == "-1" {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: candidate claims no solution\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, idx+1, input, refOut, candOut)
					os.Exit(1)
				}
				if !validateSolution(caseData.b, caseData.c, caseData.d, got) {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: invalid solution\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, idx+1, input, refOut, candOut)
					os.Exit(1)
				}
			}
		}
		if len(candAnswers) != len(refAnswers) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: extra outputs\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2020C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runExecutable(path, input string) (string, error) {
	cmd := exec.Command(path)
	return execute(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return execute(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func execute(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type testCase struct {
	b *big.Int
	c *big.Int
	d *big.Int
}

func parseCases(input string) (int, []testCase) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, nil
	}
	t := 0
	fmt.Sscanf(fields[0], "%d", &t)
	cases := make([]testCase, 0, t)
	idx := 1
	for i := 0; i < t && idx+2 < len(fields); i++ {
		b := new(big.Int)
		c := new(big.Int)
		d := new(big.Int)
		b.SetString(fields[idx], 10)
		c.SetString(fields[idx+1], 10)
		d.SetString(fields[idx+2], 10)
		idx += 3
		cases = append(cases, testCase{b: b, c: c, d: d})
	}
	return len(cases), cases
}

func validateSolution(b, c, d *big.Int, ans string) bool {
	a := new(big.Int)
	if _, ok := a.SetString(ans, 10); !ok {
		return false
	}

	limitBig := new(big.Int).SetUint64(maxA)
	if a.Sign() < 0 || a.Cmp(limitBig) > 0 {
		return false
	}

	tmp := new(big.Int)
	orVal := new(big.Int).Or(a, b)
	andVal := new(big.Int).And(a, c)
	tmp.Sub(orVal, andVal)
	return tmp.Cmp(d) == 0
}

func buildTests() []string {
	tests := []string{
		"3\n2 2 2\n4 2 6\n10 2 14\n",
		"5\n0 0 0\n1 1 0\n1 0 1\n0 1 1\n123456789 987654321 111111111\n",
	}

	randomConfigs := []struct {
		t    int
		seed int64
	}{
		{10, 1},
		{50, 2},
		{100, 3},
		{500, 4},
		{1000, 5},
		{2000, time.Now().UnixNano()},
	}
	for _, cfg := range randomConfigs {
		tests = append(tests, randomTest(cfg.t, cfg.seed))
	}
	return tests
}

func randomTest(t int, seed int64) string {
	if t < 1 {
		t = 1
	}
	if t > 100000 {
		t = 100000
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		b := randomValue(r)
		c := randomValue(r)
		d := randomValue(r)
		if r.Intn(5) == 0 {
			dx := randomValue(r)
			cx := randomValue(r)
			a := randomValue(r) & ((1 << 61) - 1)
			if r.Intn(2) == 0 {
				b = dx
				c = cx
				d = (a | b) - (a & c)
			} else {
				b = r.Uint64() & ((1 << 61) - 1)
				c = r.Uint64() & ((1 << 61) - 1)
				d = (a | b) - (a & c)
			}
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", b, c, d))
	}
	return sb.String()
}

func randomValue(r *rand.Rand) uint64 {
	mode := r.Intn(5)
	switch mode {
	case 0:
		return uint64(r.Intn(16))
	case 1:
		return uint64(r.Intn(1_000_000))
	case 2:
		return uint64(r.Intn(1_000_000_000))
	case 3:
		return uint64(r.Int63n(1<<60) + 1)
	default:
		return r.Uint64() & ((1 << 61) - 1)
	}
}
