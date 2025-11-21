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

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2147B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2147B.go")
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

func buildInput(ns []int) string {
	var sb strings.Builder
	sb.Grow(len(ns) * 16)
	sb.WriteString(strconv.Itoa(len(ns)))
	sb.WriteByte('\n')
	for _, n := range ns {
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, ns []int) ([][]int, error) {
	r := strings.Fields(out)
	idx := 0
	res := make([][]int, len(ns))
	for i, n := range ns {
		if idx+2*n > len(r) {
			return nil, fmt.Errorf("test %d: insufficient numbers in output", i+1)
		}
		arr := make([]int, 2*n)
		for j := 0; j < 2*n; j++ {
			v, err := strconv.Atoi(r[idx+j])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, r[idx+j])
			}
			arr[j] = v
		}
		idx += 2 * n
		res[i] = arr
	}
	if idx != len(r) {
		return nil, fmt.Errorf("extra tokens after parsing output")
	}
	return res, nil
}

func validateSequence(seq []int) error {
	if len(seq)%2 != 0 {
		return fmt.Errorf("sequence length %d is not even", len(seq))
	}
	count := make(map[int]int)
	for _, v := range seq {
		count[v]++
	}
	for v, c := range count {
		if c != 2 {
			return fmt.Errorf("value %d appears %d times instead of 2", v, c)
		}
	}
	return nil
}

func deterministicNs() []int {
	return []int{1, 2, 3, 5, 10, 50}
}

func randomNs(total int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := []int{}
	used := 0
	for used < total {
		n := rng.Intn(500) + 1
		res = append(res, n)
		used += n
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	ns := deterministicNs()
	const totalLimit = 10_000
	used := 0
	for _, v := range ns {
		used += v
	}
	if used < totalLimit {
		ns = append(ns, randomNs(totalLimit-used)...)
	}

	input := buildInput(ns)

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	expectedSeqs, err := parseOutput(oracleOut, ns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, oracleOut)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualSeqs, err := parseOutput(actOut, ns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actOut)
		os.Exit(1)
	}

	for i := range ns {
		if err := validateSequence(actualSeqs[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid sequence: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		// compare multisets: counts already validated; relative order can differ? Problem statement allows any? We just ensure multiset matches oracle.
		expCount := make(map[int]int)
		for _, v := range expectedSeqs[i] {
			expCount[v]++
		}
		actCount := make(map[int]int)
		for _, v := range actualSeqs[i] {
			actCount[v]++
		}
		if len(expCount) != len(actCount) {
			fmt.Fprintf(os.Stderr, "test %d mismatch distinct values\ninput:\n%s", i+1, input)
			os.Exit(1)
		}
		for k, v := range expCount {
			if actCount[k] != v {
				fmt.Fprintf(os.Stderr, "test %d value %d count mismatch: expected %d, got %d\ninput:\n%s", i+1, k, v, actCount[k], input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed.")
}
