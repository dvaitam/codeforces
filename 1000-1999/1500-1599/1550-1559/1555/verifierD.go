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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "1555D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func generate(rng *rand.Rand) (string, [][2]int) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	letters := []byte{'a', 'b', 'c'}
	for i := range b {
		b[i] = letters[rng.Intn(3)]
	}
	m := rng.Intn(20) + 1
	queries := make([][2]int, m)
	for i := range queries {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i][0] = l
		queries[i][1] = r
	}
	return string(b), queries
}

func runCase(bin, oracle string, s string, queries [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(s), len(queries)))
	sb.WriteString(s)
	sb.WriteByte('\n')
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	input := sb.String()

	cmdO := exec.Command(oracle)
	cmdO.Stdin = strings.NewReader(input)
	var outO bytes.Buffer
	cmdO.Stdout = &outO
	if err := cmdO.Run(); err != nil {
		return fmt.Errorf("oracle run error: %v", err)
	}
	expected := strings.TrimSpace(outO.String())

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, q := generate(rng)
		if err := runCase(bin, oracle, s, q); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
