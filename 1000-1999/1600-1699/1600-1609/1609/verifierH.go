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
	oracle := filepath.Join(dir, "oracleH")
	cmd := exec.Command("go", "build", "-o", oracle, "1609H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 1
	k := r.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	pos := make([]int64, n)
	cur := int64(r.Intn(10))
	for i := 0; i < n; i++ {
		pos[i] = cur
		cur += int64(r.Intn(5) + 1)
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", pos[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			sb.WriteString(fmt.Sprintf("%d", r.Intn(3)-1))
			if j+1 < k {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	q := r.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		a := r.Intn(n) + 1
		t := r.Int63n(10)
		sb.WriteString(fmt.Sprintf("%d %d\n", a, t))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
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
