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

func buildRef() (string, error) {
	exe, err := os.CreateTemp("", "refE-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, "62E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (int, int, [][]int64, [][]int64) {
	n := rng.Intn(4) + 2 // 2..5 but keep small
	m := rng.Intn(3) + 2 // 2..4
	h := make([][]int64, m-1)
	for j := 0; j < m-1; j++ {
		row := make([]int64, n)
		for i := 0; i < n; i++ {
			row[i] = int64(rng.Intn(10))
		}
		h[j] = row
	}
	v := make([][]int64, m)
	for j := 0; j < m; j++ {
		row := make([]int64, n)
		for i := 0; i < n; i++ {
			row[i] = int64(rng.Intn(10))
		}
		v[j] = row
	}
	return n, m, h, v
}

func buildInput(n, m int, h [][]int64, v [][]int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for j := 0; j < m-1; j++ {
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", h[j][i]))
		}
		sb.WriteByte('\n')
	}
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v[j][i]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, m, h, v := genCase(rng)
		input := buildInput(n, m, h, v)
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected: %s\ngot: %s\n", t+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
