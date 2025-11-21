package main

import (
	"bytes"
	"fmt"
	"math/big"
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
	keys []string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1267K-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleK")
	cmd := exec.Command("go", "build", "-o", path, "1267K.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.keys)))
	for _, v := range tc.keys {
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]*big.Int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d lines, got %d", t, len(lines))
	}
	res := make([]*big.Int, t)
	for i := 0; i < t; i++ {
		val := new(big.Int)
		_, ok := val.SetString(strings.TrimSpace(lines[i]), 10)
		if !ok {
			return nil, fmt.Errorf("invalid integer on line %d: %q", i+1, lines[i])
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{keys: []string{"1"}},
		{keys: []string{"11", "15", "178800", "123456"}},
		{keys: []string{"999999999999999999"}},
		{keys: []string{"123456789012345678"}},
	}
}

func randomKey(rng *rand.Rand) string {
	length := rng.Intn(18) + 1
	var sb strings.Builder
	sb.WriteByte(byte(rng.Intn(9)+1) + '0')
	for i := 1; i < length; i++ {
		sb.WriteByte(byte(rng.Intn(10)) + '0')
	}
	return sb.String()
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(20) + 1
	if rng.Intn(4) == 0 {
		t = rng.Intn(5) + 1
	}
	keys := make([]string, t)
	for i := 0; i < t; i++ {
		switch rng.Intn(4) {
		case 0:
			keys[i] = strconv.FormatInt(rng.Int63n(1_000_000)+1, 10)
		case 1:
			keys[i] = randomKey(rng)
		case 2:
			keys[i] = strconv.FormatInt(rng.Int63n(1_000_000_000_000)+1, 10)
		default:
			keys[i] = strconv.FormatInt(rng.Int63n(9)+1, 10)
		}
	}
	return testCase{keys: keys}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, len(tc.keys))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, len(tc.keys))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		for i := range expVals {
			if expVals[i].Cmp(gotVals[i]) != 0 {
				fmt.Fprintf(os.Stderr, "mismatch on test %d line %d: expected %s got %s\ninput:\n%s\n", idx+1, i+1, expVals[i].String(), gotVals[i].String(), input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
