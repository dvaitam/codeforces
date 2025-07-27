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

func solveCase(a [][]int) string {
	n := len(a)
	m := len(a[0])
	for i := 0; i < n; i++ {
		z := (i + 1) % 2
		for j := 0; j < m; j++ {
			if a[i][j]%2 == z {
				z = 1 - z
				continue
			}
			a[i][j]++
			z = 1 - z
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(a[i][j]))
		}
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(10) + 1
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	var exp strings.Builder
	for c := 0; c < t; c++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		in.WriteString(fmt.Sprintf("%d %d\n", n, m))
		mat := make([][]int, n)
		for i := 0; i < n; i++ {
			mat[i] = make([]int, m)
			for j := 0; j < m; j++ {
				val := rng.Intn(1000) + 1
				mat[i][j] = val
				in.WriteString(fmt.Sprintf("%d", val))
				if j+1 < m {
					in.WriteByte(' ')
				}
			}
			in.WriteByte('\n')
		}
		exp.WriteString(solveCase(mat))
		if c+1 < t {
			exp.WriteByte('\n')
		}
	}
	return in.String(), exp.String()
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\n--- got:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
