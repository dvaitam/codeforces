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

const numTestsF2 = 100

func buildOracle() (string, error) {
	exe := filepath.Join(os.TempDir(), "oracleF2")
	cmd := exec.Command("go", "build", "-o", exe, "1178F2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return exe, nil
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifF2_bin")
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

func genCaseF2(rng *rand.Rand) (int, int, []int) {
	n := rng.Intn(6) + 1
	m := n + rng.Intn(6)
	c := make([]int, m)
	for i := range c {
		c[i] = rng.Intn(n) + 1
	}
	for i := 1; i <= n; i++ {
		c[rng.Intn(m)] = i
	}
	return n, m, c
}

func runCaseF2(bin, oracle string, n, m int, c []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
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
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(7))
	for t := 0; t < numTestsF2; t++ {
		n, m, c := genCaseF2(rng)
		if err := runCaseF2(bin, oracle, n, m, c); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
