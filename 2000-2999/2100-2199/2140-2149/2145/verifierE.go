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
)

const testCount = 80

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2145E.go")
	tmp, err := os.CreateTemp("", "oracle2145E")
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

func genCase(r *rand.Rand) string {
	ac := r.Intn(50) + 1
	dr := r.Intn(50) + 1
	n := r.Intn(70) + 1
	m := r.Intn(90) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", ac, dr)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", r.Intn(60)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", r.Intn(60)+1)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		k := r.Intn(n) + 1
		na := r.Intn(60) + 1
		nd := r.Intn(60) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", k, na, nd)
	}
	return sb.String()
}

func parseOutput(out string, m int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != m {
		return nil, fmt.Errorf("expected %d outputs, got %d", m, len(fields))
	}
	res := make([]int64, m)
	for i := 0; i < m; i++ {
		val, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at line %d", i+1)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative popularity at line %d", i+1)
		}
		res[i] = val
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		input := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		scanner := strings.NewReader(input)
		var ac, dr int
		fmt.Fscan(scanner, &ac, &dr)
		var n int
		fmt.Fscan(scanner, &n)
		for i := 0; i < n; i++ {
			var tmp int
			fmt.Fscan(scanner, &tmp)
		}
		for i := 0; i < n; i++ {
			var tmp int
			fmt.Fscan(scanner, &tmp)
		}
		var m int
		fmt.Fscan(scanner, &m)
		expectVals, err := parseOutput(expectStr, m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotStr, m)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			if expectVals[i] != gotVals[i] {
				fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", t+1, input, expectVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
