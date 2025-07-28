package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildBinary(src, tag string) (string, error) {
	if strings.HasSuffix(src, ".go") {
		out := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", out, src)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", src, err, string(outb))
		}
		return out, nil
	}
	return src, nil
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func solveCase(n int, rows []string) int {
	diag := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		row := rows[i]
		for j := 0; j < n; j++ {
			if row[j] == '1' {
				total++
				diff := (i - j + n) % n
				diag[diff]++
			}
		}
	}
	best := 0
	for i := 0; i < n; i++ {
		if diag[i] > best {
			best = diag[i]
		}
	}
	return n + total - 2*best
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(10) + 1
	rows := make([]string, n)
	for i := range rows {
		b := make([]byte, n)
		for j := range b {
			if r.Intn(2) == 1 {
				b[j] = '1'
			} else {
				b[j] = '0'
			}
		}
		rows[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(rows[i])
		sb.WriteByte('\n')
	}
	input := sb.String()
	expect := fmt.Sprintf("%d\n", solveCase(n, rows))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1660E.go")

	cand, err := buildBinary(candSrc, "candE.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refE.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(ref, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		if err := runCase(cand, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
