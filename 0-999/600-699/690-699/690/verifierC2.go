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

const testCount = 200

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "690C2.go")
	tmp, err := os.CreateTemp("", "oracle690C2")
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

func genTree(r *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := r.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return edges
}

func genExtraEdges(r *rand.Rand, n, m int, existing map[[2]int]struct{}) [][2]int {
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := existing[key]; ok {
			continue
		}
		existing[key] = struct{}{}
		edges = append(edges, key)
	}
	return edges
}

func genCase(r *rand.Rand) (string, int) {
	n := 1 + r.Intn(400)
	if r.Intn(5) == 0 {
		n = 1 + r.Intn(2000)
	}
	treeEdges := genTree(r, n)
	m := len(treeEdges)
	extra := 0
	if n > 1 {
		extra = r.Intn(n)
	}
	m += extra
	existing := make(map[[2]int]struct{})
	for _, e := range treeEdges {
		a, b := e[0], e[1]
		if a > b {
			a, b = b, a
		}
		existing[[2]int{a, b}] = struct{}{}
	}
	extraEdges := genExtraEdges(r, n, extra, existing)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(treeEdges)+len(extraEdges))
	for _, e := range treeEdges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for _, e := range extraEdges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String(), 1
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, fmt.Errorf("negative diameter")
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
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
		input, _ := genCase(r)
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
		expectVal, err := parseOutput(expectStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if expectVal != gotVal {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", t+1, input, expectVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
