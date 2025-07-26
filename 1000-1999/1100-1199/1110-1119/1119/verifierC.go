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

type testCase struct {
	n, m int
	A, B [][]int
}

func expected(tc testCase) string {
	r := make([]int, tc.n)
	c := make([]int, tc.m)
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			diff := tc.A[i][j] ^ tc.B[i][j]
			r[i] ^= diff
			c[j] ^= diff
		}
	}
	for i := 0; i < tc.n; i++ {
		if r[i] != 0 {
			return "No"
		}
	}
	for j := 0; j < tc.m; j++ {
		if c[j] != 0 {
			return "No"
		}
	}
	return "Yes"
}

func buildCase(A, B [][]int) (string, string) {
	n := len(A)
	m := len(A[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", A[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", B[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expected(testCase{n: n, m: m, A: A, B: B})
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	A := make([][]int, n)
	B := make([][]int, n)
	for i := 0; i < n; i++ {
		A[i] = make([]int, m)
		B[i] = make([]int, m)
		for j := 0; j < m; j++ {
			A[i][j] = rng.Intn(2)
			B[i][j] = rng.Intn(2)
		}
	}
	return buildCase(A, B)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	A := [][]int{{0, 1}, {1, 0}}
	B := [][]int{{1, 0}, {0, 1}}
	in, exp := buildCase(A, B)
	cases = append(cases, in)
	exps = append(exps, exp)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 102 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}
	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
