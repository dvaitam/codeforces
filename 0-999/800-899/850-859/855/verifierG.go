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

const testCount = 120

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "855G.go")
	tmp, err := os.CreateTemp("", "oracle855G")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
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

func genTree(n int, r *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := r.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func genCase(r *rand.Rand) (string, int, int) {
	n := 2 + r.Intn(8) // keep small for oracle
	edges := genTree(n, r)
	q := 1 + r.Intn(6)
	extra := make([][2]int, q)
	for i := 0; i < q; i++ {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		extra[i] = [2]int{u, v}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", q)
	for _, e := range extra {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String(), q + 1, n
}

func parseOutput(out string, cnt int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != cnt {
		return nil, fmt.Errorf("expected %d numbers, got %d", cnt, len(fields))
	}
	res := make([]int64, cnt)
	for i := 0; i < cnt; i++ {
		v, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, err
		}
		if v < 0 {
			return nil, fmt.Errorf("negative result")
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
		input, count, _ := genCase(r)
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
		expectVals, err := parseOutput(expectStr, count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotStr, count)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		for i := 0; i < count; i++ {
			if expectVals[i] != gotVals[i] {
				fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", t+1, input, expectVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
