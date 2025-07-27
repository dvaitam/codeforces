package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const numTestsC = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "binC*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp, err := os.CreateTemp("", "oracleC*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1386C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	maxEdges := n * (n - 1) / 2
	if maxEdges == 0 {
		maxEdges = 1
	}
	m := rng.Intn(maxEdges) + 1
	q := rng.Intn(4) + 1

	// generate unique edges
	pairs := make([][2]int, 0, maxEdges)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			pairs = append(pairs, [2]int{i, j})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	edges := pairs[:m]

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for i := 0; i < q; i++ {
		l := rng.Intn(m) + 1
		r := rng.Intn(m-l+1) + l
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	exp, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	argIdx := 1
	if len(os.Args) >= 3 && os.Args[1] == "--" {
		argIdx = 2
	}
	if len(os.Args) != argIdx+1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[argIdx])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println("oracle compile error:", err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTestsC; i++ {
		in := genCase(rng)
		if err := runCase(bin, oracle, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			return
		}
	}
	fmt.Println("All tests passed")
}
