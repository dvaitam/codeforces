package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n, m, k int
	b, c    []int
	exp     string
}

func solveA(n, m, k int, b, c []int) string {
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if b[i]+c[j] <= k {
				cnt++
			}
		}
	}
	return fmt.Sprint(cnt)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCaseA {
	rng := rand.New(rand.NewSource(1))
	cases := make([]testCaseA, 100)
	for i := range cases {
		n := rng.Intn(8) + 1
		m := rng.Intn(8) + 1
		k := rng.Intn(50) + 1
		b := make([]int, n)
		c := make([]int, m)
		for j := 0; j < n; j++ {
			b[j] = rng.Intn(40) + 1
		}
		for j := 0; j < m; j++ {
			c[j] = rng.Intn(40) + 1
		}
		cases[i] = testCaseA{n: n, m: m, k: k, b: b, c: c, exp: solveA(n, m, k, b, c)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, 1)
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.b[j]))
		}
		sb.WriteByte('\n')
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.c[j]))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
