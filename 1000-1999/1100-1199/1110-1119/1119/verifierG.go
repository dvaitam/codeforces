package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	N, M int
	H    []int64
}

func referenceSolve(tc testCase) string {
	N, M := tc.N, tc.M
	H := append([]int64(nil), tc.H...)
	cutoffs := make([]int64, 0, M+1)
	cutoffs = append(cutoffs, 0)
	cutoffs = append(cutoffs, int64(N))
	var cur int64
	for i := 0; i+1 < M; i++ {
		cur += H[i]
		cutoffs = append(cutoffs, cur%int64(N))
	}
	sort.Slice(cutoffs, func(i, j int) bool { return cutoffs[i] < cutoffs[j] })
	sizes := make([]int64, M)
	for i := 0; i < M; i++ {
		sizes[i] = cutoffs[i+1] - cutoffs[i]
	}
	ans := make([]int, 0)
	ind := 0
	for i := 0; i < M; i++ {
		v := H[i]
		for v > 0 {
			ans = append(ans, i)
			v -= sizes[ind%M]
			ind++
		}
	}
	for len(ans)%M != 0 {
		ans = append(ans, 0)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)/M))
	for i := 0; i < M; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", sizes[i]))
	}
	sb.WriteByte('\n')
	for i, x := range ans {
		if i > 0 && i%M == 0 {
			sb.WriteByte('\n')
		}
		if i%M > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x+1))
	}
	sb.WriteByte('\n')
	return strings.TrimRight(sb.String(), "\n")
}

func buildCase(N, M int, H []int64) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", N, M))
	for i, v := range H {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), referenceSolve(testCase{N: N, M: M, H: H})
}

func genCase(rng *rand.Rand) (string, string) {
	N := rng.Intn(10) + 1
	M := rng.Intn(N) + 1
	H := make([]int64, M)
	for i := 0; i < M; i++ {
		H[i] = int64(rng.Intn(5) + 1)
	}
	return buildCase(N, M, H)
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
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	in, exp := buildCase(3, 2, []int64{2, 1})
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
