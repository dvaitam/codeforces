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

const maxR int64 = 1_000_000_000_000

type testCase struct {
	r int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1184A1-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleA1")
	cmd := exec.Command("go", "build", "-o", path, "1184A1.go")
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

func parseAnswer(out string) (bool, int64, int64, error) {
	if out == "NO" {
		return false, 0, 0, nil
	}
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return false, 0, 0, fmt.Errorf("expected two integers or NO")
	}
	x, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil || x <= 0 {
		return false, 0, 0, fmt.Errorf("invalid x: %v", err)
	}
	y, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil || y <= 0 {
		return false, 0, 0, fmt.Errorf("invalid y: %v", err)
	}
	return true, x, y, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{r: 1},
		{r: 2},
		{r: 3},
		{r: 10},
		{r: 1000},
		{r: 999983},
	}
}

func randomTest(rng *rand.Rand) testCase {
	var r int64
	switch rng.Intn(3) {
	case 0:
		r = int64(rng.Intn(1000) + 1)
	case 1:
		r = int64(rng.Int63n(maxR) + 1)
	default:
		// generate some values likely to have solutions
		x := int64(rng.Intn(100000) + 1)
		y := int64(rng.Intn(100000) + 1)
		r = x*x + 2*x*y + x + 1
		if r > maxR {
			r = maxR
		}
	}
	return testCase{r: r}
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n", tc.r)
}

func evaluate(r, x, y int64) bool {
	val := x*x + 2*x*y + x + 1
	return val == r
}

func checkMinimality(r, x int64) bool {
	for t := int64(1); t < x; t++ {
		minVal := t*t + t + 1
		if minVal >= r {
			break
		}
		N := r - minVal
		if N%(2*t) == 0 {
			if y := N / (2 * t); y >= 1 {
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (r=%d): %v\n", idx+1, tc.r, err)
			os.Exit(1)
		}
		expHas, _, _, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotHas, gotX, gotY, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		if expHas != gotHas {
			fmt.Fprintf(os.Stderr, "test %d: expected existence=%v got %v\ninput:\n%s\n", idx+1, expHas, gotHas, input)
			os.Exit(1)
		}
		if gotHas {
			if !evaluate(tc.r, gotX, gotY) {
				fmt.Fprintf(os.Stderr, "test %d: candidate (%d, %d) does not satisfy equation for r=%d\n", idx+1, gotX, gotY, tc.r)
				os.Exit(1)
			}
			if !checkMinimality(tc.r, gotX) {
				fmt.Fprintf(os.Stderr, "test %d: x=%d is not minimal for r=%d\n", idx+1, gotX, tc.r)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
