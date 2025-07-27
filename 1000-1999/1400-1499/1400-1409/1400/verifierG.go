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
	oracle := filepath.Join(dir, "oracleG")
	cmd := exec.Command("go", "build", "-o", oracle, "1400G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func genCase(r *rand.Rand) string {
	n := r.Intn(10) + 1
	maxM := n * (n - 1) / 2
	if maxM > 3 {
		maxM = 3
	}
	m := r.Intn(maxM + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		l := r.Intn(n) + 1
		rgt := r.Intn(n-l+1) + l
		sb.WriteString(fmt.Sprintf("%d %d\n", l, rgt))
	}
	pairs := make(map[[2]int]bool)
	for len(pairs) < m {
		a := r.Intn(n) + 1
		b := r.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if !pairs[key] {
			pairs[key] = true
		}
	}
	for p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
