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
	"strings"
	"time"
)

type testCase struct {
	n string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1194G-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", path, "1194G.go")
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

func parseAns(out string) (*big.Int, error) {
	val := new(big.Int)
	_, ok := val.SetString(out, 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer output %q", out)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: "1"},
		{n: "9"},
		{n: "10"},
		{n: "99"},
		{n: "12345678901234567890"},
		{n: "99999999999999999999"},
	}
}

func randomDecimal(rng *rand.Rand, length int) string {
	if length <= 0 {
		return "1"
	}
	var sb strings.Builder
	sb.WriteByte(byte(rng.Intn(9)+1) + '0')
	for i := 1; i < length; i++ {
		sb.WriteByte(byte(rng.Intn(10)) + '0')
	}
	return sb.String()
}

func randomTest(rng *rand.Rand) testCase {
	var n string
	switch rng.Intn(3) {
	case 0:
		n = fmt.Sprintf("%d", rng.Intn(1000000)+1)
	case 1:
		n = randomDecimal(rng, rng.Intn(30)+1)
	default:
		n = randomDecimal(rng, rng.Intn(100)+1)
	}
	return testCase{n: n}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
		input := fmt.Sprintf("%s\n", tc.n)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (n=%s): %v\n", idx+1, tc.n, err)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d (n=%s): %v\n", idx+1, tc.n, err)
			os.Exit(1)
		}
		expVal, err := parseAns(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		gotVal, err := parseAns(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}
		if expVal.Cmp(gotVal) != 0 {
			fmt.Fprintf(os.Stderr, "mismatch on test %d n=%s: expected %s got %s\n", idx+1, tc.n, expVal.String(), gotVal.String())
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
