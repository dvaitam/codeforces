package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const mod = 998244353

type testCase struct {
	desc  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, tc := range tests {
		expect, err := runAndParse(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		got, err := runAndParse(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\ninput:\n%sexpected: %d\ngot: %d\n", i+1, tc.desc, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1431J-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	refPath := filepath.Clean("1000-1999/1400-1499/1430-1439/1431/1431J.go")
	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runAndParse(bin, input string) (int, error) {
	out, err := runProgram(bin, input)
	if err != nil {
		return 0, err
	}
	ans, err := parseAnswer(out)
	if err != nil {
		return 0, err
	}
	return ans, nil
}

func runProgram(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val := new(big.Int)
	if _, ok := val.SetString(fields[0], 10); !ok {
		return 0, fmt.Errorf("cannot parse integer %q", fields[0])
	}
	val.Mod(val, big.NewInt(mod))
	if val.Sign() < 0 {
		val.Add(val, big.NewInt(mod))
	}
	ans := int(val.Int64())
	if len(fields) > 1 {
		return 0, fmt.Errorf("unexpected extra tokens after answer")
	}
	return ans, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeTest("n2_small_equal", []uint64{0, 0}),
		makeTest("n2_diff", []uint64{1, 2}),
		makeTest("n3_progression", []uint64{1, 1, 2}),
		makeTest("n3_large_gap", []uint64{0, 5, 5}),
		makeTest("n4_mixed", []uint64{0, 3, 7, 7}),
		makeTest("n5_all_zero", []uint64{0, 0, 0, 0, 0}),
		makeTest("max_n_small_vals", []uint64{0, 1, 1, 2, 3, 5, 8, 13, 13, 21, 21, 34, 55, 55, 89, 144, 233}),
		makeTest("max_n_large_vals", []uint64{
			1 << 59, (1 << 59) + 1, (1 << 59) + 1, (1 << 59) + 2,
			(1 << 59) + 4, (1 << 59) + 8, (1 << 59) + 8, (1 << 59) + 16,
			(1 << 59) + 32, (1 << 59) + 32, (1 << 59) + 64, (1 << 59) + 128,
			(1 << 59) + 128, (1 << 59) + 256, (1 << 59) + 512, (1 << 59) + 512, (1 << 59) + 1024,
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		n := rng.Intn(16) + 2 // 2..17
		arr := make([]uint64, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Uint64() & ((1 << 60) - 1)
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		tests = append(tests, makeTest(fmt.Sprintf("rand-%d", i+1), arr))
	}
	return tests
}

func makeTest(desc string, a []uint64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{
		desc:  desc,
		input: sb.String(),
	}
}
