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

const testCount = 100

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1184B3.go")
	tmp, err := os.CreateTemp("", "oracle1184B3")
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

func genGraph(n, m int, r *rand.Rand) [][2]int {
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		edges = append(edges, [2]int{u, v})
	}
	return edges
}

func genShips(s int, n int, r *rand.Rand) []string {
	ships := make([]string, s)
	for i := 0; i < s; i++ {
		x := r.Intn(n) + 1
		a := r.Intn(20)
		f := r.Intn(10)
		p := r.Intn(20)
		ships[i] = fmt.Sprintf("%d %d %d %d", x, a, f, p)
	}
	return ships
}

func genBases(b int, n int, r *rand.Rand) []string {
	bases := make([]string, b)
	for i := 0; i < b; i++ {
		x := r.Intn(n) + 1
		d := r.Intn(20)
		g := r.Intn(30)
		bases[i] = fmt.Sprintf("%d %d %d", x, d, g)
	}
	return bases
}

func genDeps(k, s int, r *rand.Rand) [][2]int {
	deps := make([][2]int, 0, k)
	for len(deps) < k {
		u := r.Intn(s) + 1
		v := r.Intn(s) + 1
		deps = append(deps, [2]int{u, v})
	}
	return deps
}

func genCase(r *rand.Rand) string {
	n := r.Intn(8) + 1
	m := r.Intn(10)
	edges := genGraph(n, m, r)
	s := r.Intn(6) + 1
	b := r.Intn(6) + 1
	k := r.Intn(5)
	ships := genShips(s, n, r)
	bases := genBases(b, n, r)
	deps := genDeps(k, s, r)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d %d %d\n", s, b, k)
	for _, ship := range ships {
		sb.WriteString(ship)
		sb.WriteByte('\n')
	}
	for _, base := range bases {
		sb.WriteString(base)
		sb.WriteByte('\n')
	}
	for _, dep := range deps {
		fmt.Fprintf(&sb, "%d %d\n", dep[0], dep[1])
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB3.go /path/to/binary")
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
