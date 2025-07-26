package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsE = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "binE")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleE")
	cmd := exec.Command("go", "build", "-o", tmp, "1217E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		val := rng.Int63n(1e9-1) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	type2 := false
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			idx := rng.Intn(n) + 1
			val := rng.Int63n(1e9-1) + 1
			fmt.Fprintf(&sb, "1 %d %d\n", idx, val)
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "2 %d %d\n", l, r)
			type2 = true
		}
	}
	if !type2 {
		// ensure at least one type2 query
		l := 1
		r := n
		fmt.Fprintf(&sb, "2 %d %d\n", l, r)
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	expected, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	path := os.Args[len(os.Args)-1]
	bin, cleanup, err := prepareBinary(path)
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(5))
	for i := 0; i < numTestsE; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
