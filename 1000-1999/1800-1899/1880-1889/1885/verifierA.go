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
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "1885A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	N := r.Intn(4) + 1
	K := r.Intn(3) + 1
	T := r.Intn(4) + 1
	R := r.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n%d\n%d\n%d\n", N, K, T, R)
	for t := 0; t < T; t++ {
		for k := 0; k < K; k++ {
			for rg := 0; rg < R; rg++ {
				for n := 0; n < N; n++ {
					if n > 0 {
						sb.WriteByte(' ')
					}
					val := r.Float64()*100 + 0.01
					sb.WriteString(fmt.Sprintf("%.6f", val))
				}
				sb.WriteByte('\n')
			}
		}
	}
	for m := 0; m < N; m++ {
		for rg := 0; rg < R; rg++ {
			for k := 0; k < K; k++ {
				for n := 0; n < N; n++ {
					if n > 0 {
						sb.WriteByte(' ')
					}
					val := -r.Float64() * 2
					sb.WriteString(fmt.Sprintf("%.6f", val))
				}
				sb.WriteByte('\n')
			}
		}
	}
	J := r.Intn(5) + 1
	fmt.Fprintf(&sb, "%d\n", J)
	for j := 0; j < J; j++ {
		t0 := r.Intn(T)
		td := r.Intn(T-t0) + 1
		tbs := r.Intn(100) + 1
		user := r.Intn(N)
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", j, tbs, user, t0, td)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
