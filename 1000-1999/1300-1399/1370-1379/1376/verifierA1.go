package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const testCount = 100

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1376A1.go")
	tmp, err := os.CreateTemp("", "oracle1376A1")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, []int) {
	n := r.Intn(50) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = r.Intn(200) - 100
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func isSortedSame(arr, sorted []int) bool {
	if len(arr) != len(sorted) {
		return false
	}
	for i := range arr {
		if arr[i] != sorted[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input, arr := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		if _, err := parseOutput(expectStr, len(arr)); err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotStr, len(arr))
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		// verify sorted result and same multiset
		sortedArr := append([]int(nil), arr...)
		sort.Ints(sortedArr)
		sortGot := append([]int(nil), gotVals...)
		if !isSortedSame(sortGot, sortedArr) {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected sorted:%v\ngot:%v\n", t+1, input, sortedArr, gotVals)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
