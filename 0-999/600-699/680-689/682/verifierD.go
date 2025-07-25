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
	n int
	m int
	k int
	s string
	t string
}

func run(bin, input string) (string, error) {
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
	return out.String(), nil
}

func solve(tc testCase) int {
	n := tc.n
	m := tc.m
	K := tc.k
	s := tc.s
	t := tc.t
	f := make([][][]int, K+1)
	for k := range f {
		f[k] = make([][]int, n+1)
		for i := range f[k] {
			f[k][i] = make([]int, m+1)
		}
	}
	g := make([][]int, n+1)
	for i := range g {
		g[i] = make([]int, m+1)
	}
	for k := 1; k <= K; k++ {
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				if s[i-1] == t[j-1] {
					v1 := g[i-1][j-1] + 1
					v2 := f[k-1][i-1][j-1] + 1
					if v2 > v1 {
						v1 = v2
					}
					g[i][j] = v1
				} else {
					g[i][j] = 0
				}
				v := f[k][i-1][j]
				if f[k][i][j-1] > v {
					v = f[k][i][j-1]
				}
				if g[i][j] > v {
					v = g[i][j]
				}
				f[k][i][j] = v
			}
		}
	}
	return f[K][n][m]
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d\n%s\n%s\n", tc.n, tc.m, tc.k, tc.s, tc.t)
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expected := solve(tc)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	m := rng.Intn(20) + 1
	k := rng.Intn(10) + 1
	s := randString(rng, n)
	t := randString(rng, m)
	return testCase{n: n, m: m, k: k, s: s, t: t}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, testCase{n: 1, m: 1, k: 1, s: "a", t: "a"})
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
