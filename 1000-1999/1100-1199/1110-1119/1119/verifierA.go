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
	colors []int
}

func expected(colors []int) string {
	n := len(colors)
	ans := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if colors[i] != colors[j] && j-i > ans {
				ans = j - i
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func buildCase(colors []int) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(colors)))
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(colors)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 3
	colors := make([]int, n)
	for i := range colors {
		colors[i] = rng.Intn(n) + 1
	}
	allSame := true
	for i := 1; i < n; i++ {
		if colors[i] != colors[0] {
			allSame = false
			break
		}
	}
	if allSame {
		colors[n-1] = colors[0] + 1
	}
	return buildCase(colors)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []string{}
	exps := []string{}
	// deterministic cases
	in, exp := buildCase([]int{1, 2, 1, 2})
	cases = append(cases, in)
	exps = append(exps, exp)
	in, exp = buildCase([]int{1, 1, 1, 2})
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
