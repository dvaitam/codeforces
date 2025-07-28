package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const numTestsA = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "binA")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleA")
	cmd := exec.Command("go", "build", "-o", tmp, "1799A.go")
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
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", n+1+rng.Intn(m))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin, oracle, input string) error {
	expect, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	binPath := os.Args[1]
	bin, cleanup, err := prepareBinary(binPath)
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTestsA; i++ {
		input := generateCase(rng)
		if err := runCase(bin, oracle, input); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
