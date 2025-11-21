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

const refSource = "2000-2999/2000-2099/2090-2099/2094/2094D.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2094D-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func parseT(input string) (int, error) {
	var t int
	if _, err := fmt.Fscan(strings.NewReader(input), &t); err != nil {
		return 0, err
	}
	return t, nil
}

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	for i := range fields {
		fields[i] = strings.ToUpper(fields[i])
	}
	return fields, nil
}

func fixedTests() []string {
	return []string{
		"5\nR\nR\nLR\nLLRR\nLL\nLLLR\nRRRLL\nRRLLLL\nLRLR\nLLRLRR\n",
		"4\nL\nLL\nR\nLL\nLRLRLR\nLRLRLR\nRRR\nRRRRRR\n",
	}
}

func randomString(rng *rand.Rand, n int) string {
	buf := make([]byte, n)
	for i := range buf {
		if rng.Intn(2) == 0 {
			buf[i] = 'L'
		} else {
			buf[i] = 'R'
		}
	}
	return string(buf)
}

func randomTests(rng *rand.Rand) string {
	type pair struct {
		p string
		s string
	}
	var cases []pair
	remaining := 200000
	targetCases := rng.Intn(25) + 5
	for len(cases) < targetCases && remaining > 0 {
		maxLen := 5000
		if maxLen > remaining {
			maxLen = remaining
		}
		if maxLen == 0 {
			break
		}
		pLen := rng.Intn(maxLen) + 1
		maxSLen := remaining
		if maxSLen < pLen {
			maxSLen = pLen
		}
		sLen := pLen + rng.Intn(maxSLen-pLen+1)
		if sLen > remaining {
			sLen = remaining
		}
		p := randomString(rng, pLen)
		s := randomString(rng, sLen)
		cases = append(cases, pair{p: p, s: s})
		remaining -= sLen
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(c.p)
		sb.WriteByte('\n')
		sb.WriteString(c.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func bigCase() string {
	n := 200000
	p := randomString(rand.New(rand.NewSource(42)), n/2)
	s := p + p // ensure length 2n/2 = n and valid duplication pattern per char
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(p)
	sb.WriteByte('\n')
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := fixedTests()
	tests = append(tests, bigCase())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTests(rng))
	}

	for idx, input := range tests {
		t, err := parseT(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse t for case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expTokens, err := parseOutput(expOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse failed on case %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotTokens, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse failed on case %d: %v\noutput:\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := 0; i < t; i++ {
			if expTokens[i] != gotTokens[i] {
				fmt.Fprintf(os.Stderr, "mismatch on case %d test %d: expected %s got %s\ninput:\n%s", idx+1, i+1, expTokens[i], gotTokens[i], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
