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
	src := filepath.Join(dir, "536E.go")
	tmp, err := os.CreateTemp("", "oracle536E")
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

func genCase(r *rand.Rand) (string, int) {
	n := 2 + r.Intn(40)
	if r.Intn(4) == 0 {
		n = 2 + r.Intn(200)
	}
	q := 1 + r.Intn(80)
	if r.Intn(4) == 0 {
		q = 1 + r.Intn(200)
	}
	f := make([]int, n)
	for i := 1; i < n; i++ {
		f[i] = r.Intn(2001) - 1000
	}
	type edge struct{ u, v, w int }
	edges := make([]edge, 0, n-1)
	weights := make([]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := r.Intn(v-1) + 1
		w := r.Intn(1000000000) + 1
		edges = append(edges, edge{u: p, v: v, w: w})
		weights = append(weights, w)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 1; i < n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", f[i])
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	for i := 0; i < q; i++ {
		v := r.Intn(n) + 1
		u := r.Intn(n-1) + 1
		if u >= v {
			u++
		}
		l := r.Intn(1000000000) + 1
		if len(weights) > 0 && r.Intn(3) == 0 {
			l = weights[r.Intn(len(weights))]
		}
		fmt.Fprintf(&sb, "%d %d %d\n", v, u, l)
	}
	return sb.String(), q
}

func parseOutputs(out string, q int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != q {
		return nil, fmt.Errorf("expected %d numbers, got %d", q, len(fields))
	}
	res := make([]int64, q)
	for i := 0; i < q; i++ {
		v, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
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
	for i := 0; i < testCount; i++ {
		input, q := genCase(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expectVals, err := parseOutputs(expectStr, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotStr, q)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", i+1, input, err)
			os.Exit(1)
		}
		for j := 0; j < q; j++ {
			if expectVals[j] != gotVals[j] {
				fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", i+1, input, expectVals[j], gotVals[j])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
