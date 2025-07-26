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
)

const numTestsG = 100

func buildOracle() (string, error) {
	exe := filepath.Join(os.TempDir(), "oracleG")
	cmd := exec.Command("go", "build", "-o", exe, "1178G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return exe, nil
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifG_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func genTree(rng *rand.Rand, n int) []int {
	p := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p[i] = rng.Intn(i-1) + 1
	}
	return p[2:]
}

func genCaseG(rng *rand.Rand) (int, int, []int, []int64, []int64, []string) {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	p := genTree(rng, n)
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(11) - 5)
		b[i] = int64(rng.Intn(11) - 5)
	}
	queries := make([]string, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			v := rng.Intn(n) + 1
			x := rng.Intn(5) + 1
			queries[i] = fmt.Sprintf("1 %d %d", v, x)
		} else {
			v := rng.Intn(n) + 1
			queries[i] = fmt.Sprintf("2 %d", v)
		}
	}
	return n, q, p, a, b, queries
}

func runCaseG(bin, oracle string, n, q int, p []int, a, b []int64, queries []string) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, q := range queries {
		sb.WriteString(q)
		sb.WriteByte('\n')
	}
	input := sb.String()
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	oracle, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(oracle)
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(8))
	for t := 0; t < numTestsG; t++ {
		n, q, p, a, bVals, qs := genCaseG(rng)
		if err := runCaseG(bin, oracle, n, q, p, a, bVals, qs); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
