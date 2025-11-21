package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "193B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func runProgram(bin, input string) (string, error) {
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

func writeLine(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
}

func buildInput(n, u, r int, a, b, k, p []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, u, r)
	writeLine(&sb, a)
	writeLine(&sb, b)
	writeLine(&sb, k)
	writeLine(&sb, p)
	return sb.String()
}

func permutation(n int, rnd *rand.Rand) []int {
	perm := rnd.Perm(n)
	for i := 0; i < n; i++ {
		perm[i]++ // convert to 1-based indexing
	}
	return perm
}

func randomArray(n int, minVal, maxVal int, rnd *rand.Rand) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rnd.Intn(maxVal-minVal+1) + minVal
	}
	return arr
}

func randomSignedArray(n int, absLimit int, rnd *rand.Rand) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rnd.Intn(2*absLimit+1) - absLimit
	}
	return arr
}

func parseScore(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid score %q", fields[0])
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	var tests []string
	// Deterministic edge-ish cases
	tests = append(tests, buildInput(
		1, 1, 0,
		[]int{5},
		[]int{7},
		[]int{2},
		[]int{1},
	))
	tests = append(tests, buildInput(
		2, 2, 3,
		[]int{1, 2},
		[]int{4, 5},
		[]int{6, -6},
		[]int{2, 1},
	))
	tests = append(tests, buildInput(
		3, 5, 100,
		[]int{10000, 10000, 10000},
		[]int{10000, 9999, 2},
		[]int{-10000, 0, 10000},
		[]int{3, 1, 2},
	))

	// Randomized tests
	const randomTests = 200
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTests; i++ {
		n := rnd.Intn(30) + 1
		u := rnd.Intn(30) + 1
		r := rnd.Intn(101)
		a := randomArray(n, 1, 10000, rnd)
		b := randomArray(n, 1, 10000, rnd)
		k := randomSignedArray(n, 10000, rnd)
		p := permutation(n, rnd)
		tests = append(tests, buildInput(n, u, r, a, b, k, p))
	}

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expScore, err := parseScore(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotScore, err := parseScore(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if expScore != gotScore {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expScore, gotScore)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
