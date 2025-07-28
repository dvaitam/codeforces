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
	exe := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", exe, "1878D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return exe, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randString(r *rand.Rand, n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte('a' + r.Intn(26))
	}
	return string(bytes)
}

func genCase(r *rand.Rand) string {
	n := r.Intn(10) + 1
	k := r.Intn(n) + 1
	s := randString(r, n)
	l := make([]int, k)
	rarr := make([]int, k)
	start := 1
	for i := 0; i < k; i++ {
		remain := k - i
		maxLen := n - start + 1 - (remain - 1)
		length := r.Intn(maxLen) + 1
		l[i] = start
		rarr[i] = start + length - 1
		start = rarr[i] + 1
	}
	q := r.Intn(10) + 1
	xs := make([]int, q)
	for i := range xs {
		xs[i] = r.Intn(n) + 1
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	sb.WriteString(s)
	sb.WriteByte('\n')
	for i, v := range l {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range rarr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
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
