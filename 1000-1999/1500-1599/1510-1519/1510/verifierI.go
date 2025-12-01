package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "./1510I.go"

type testCase struct {
	input string
	info  testInfo
}

type testInfo struct {
	n, m   int
	guess  []string
	actual []byte
	best   int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests, err := buildTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build tests:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refOut, err := runExecutable(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(refOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "reference invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(candOut, tc.info); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\nInput:\n%sCandidate output:\n%s\n", i+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1510I-ref-*")
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

func buildTests() ([]testCase, error) {
	var tests []testCase
	baseInputs := []string{
		"3 4\n000\n1\n100\n1\n001\n0\n111\n1\n",
		"5 5\n01010\n0\n11111\n1\n00000\n0\n10101\n1\n11100\n1\n",
	}
	for _, inp := range baseInputs {
		info, err := prepareTest(inp)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: inp, info: info})
	}

	randomConfigs := []struct {
		n, m int
		seed int64
		perf bool
	}{
		{5, 20, 1, true},
		{10, 40, 2, false},
		{25, 120, 3, true},
		{60, 300, 4, false},
		{200, 800, 5, true},
		{500, 1200, 6, false},
		{1000, 1500, 7, true},
		{1000, 2000, time.Now().UnixNano(), false},
	}

	for _, cfg := range randomConfigs {
		input := randomTest(cfg.n, cfg.m, cfg.seed, cfg.perf)
		info, err := prepareTest(input)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{input: input, info: info})
	}

	return tests, nil
}

func randomTest(n, m int, seed int64, ensurePerfect bool) string {
	if n < 1 {
		n = 1
	}
	if m < 1 {
		m = 1
	}
	r := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		actual := byte('0' + r.Intn(2))
		var line strings.Builder
		for j := 0; j < n; j++ {
			var ch byte
			switch {
			case ensurePerfect && j == 0:
				// Participant 1 always right.
				ch = actual
			case !ensurePerfect && j == 0:
				// Participant 1 mostly right but not perfect.
				if r.Float64() < 0.7 {
					ch = actual
				} else {
					ch = flip(actual)
				}
			case ensurePerfect && j == 1 && n > 1:
				// Participant 2 always wrong to vary best.
				ch = flip(actual)
			default:
				if r.Float64() < 0.5 {
					ch = '0'
				} else {
					ch = '1'
				}
			}
			line.WriteByte(ch)
		}
		sb.WriteString(line.String())
		sb.WriteByte('\n')
		sb.WriteByte(actual)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func flip(b byte) byte {
	if b == '0' {
		return '1'
	}
	return '0'
}

func prepareTest(input string) (testInfo, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return testInfo{}, fmt.Errorf("failed to read n and m: %v", err)
	}
	if n <= 0 || m <= 0 {
		return testInfo{}, fmt.Errorf("n and m must be positive")
	}
	guesses := make([]string, m)
	actual := make([]byte, m)
	for i := 0; i < m; i++ {
		var line string
		if _, err := fmt.Fscan(reader, &line); err != nil {
			return testInfo{}, fmt.Errorf("failed to read prediction %d: %v", i+1, err)
		}
		if len(line) != n {
			return testInfo{}, fmt.Errorf("prediction %d has length %d, expected %d", i+1, len(line), n)
		}
		for _, ch := range line {
			if ch != '0' && ch != '1' {
				return testInfo{}, fmt.Errorf("prediction %d has invalid character %q", i+1, ch)
			}
		}
		guesses[i] = line

		var res string
		if _, err := fmt.Fscan(reader, &res); err != nil {
			return testInfo{}, fmt.Errorf("failed to read result %d: %v", i+1, err)
		}
		if len(res) != 1 || (res[0] != '0' && res[0] != '1') {
			return testInfo{}, fmt.Errorf("result %d is invalid: %q", i+1, res)
		}
		actual[i] = res[0]
	}

	best := computeBest(guesses, actual, n)

	return testInfo{
		n:      n,
		m:      m,
		guess:  guesses,
		actual: actual,
		best:   best,
	}, nil
}

func computeBest(guesses []string, actual []byte, n int) int {
	mistakes := make([]int, n)
	for round := 0; round < len(guesses); round++ {
		for i := 0; i < n; i++ {
			if guesses[round][i] != actual[round] {
				mistakes[i]++
			}
		}
	}
	best := mistakes[0]
	for _, v := range mistakes {
		if v < best {
			best = v
		}
	}
	return best
}

func checkOutput(output string, info testInfo) error {
	tokens := strings.Fields(output)
	if len(tokens) != info.m {
		return fmt.Errorf("expected %d answers, got %d", info.m, len(tokens))
	}
	mistakes := 0
	for i := 0; i < info.m; i++ {
		tok := tokens[i]
		if tok != "0" && tok != "1" {
			return fmt.Errorf("answer %d is not 0/1: %q", i+1, tok)
		}
		if tok[0] != info.actual[i] {
			mistakes++
		}
	}
	limit := 1.3*float64(info.best) + 100
	if float64(mistakes) > limit+1e-9 {
		return fmt.Errorf("too many mistakes: %d, allowed %.2f (best=%d)", mistakes, limit, info.best)
	}
	return nil
}
