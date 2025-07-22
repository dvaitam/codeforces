package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var n, m, k int
	fmt.Fscan(rdr, &n, &m, &k)
	fac := make([]float64, m+1)
	fac[0] = 0.0
	for i := 1; i <= m; i++ {
		fac[i] = fac[i-1] + math.Log(float64(i))
	}
	C := func(a, b int) float64 { return fac[b] - fac[a] - fac[b-a] }
	var ans float64
	for i := 0; i <= n; i++ {
		for j := 0; j <= n; j++ {
			t := i*n + j*n - i*j
			if t > k {
				break
			}
			tmp := C(k-t, m-t) + C(i, n) + C(j, n) - C(k, m)
			ans += math.Exp(tmp)
			if ans > 1e99 {
				return "1e99"
			}
		}
	}
	return fmt.Sprintf("%.15f", ans)
}

func genRandomCase() string {
	n := rand.Intn(5) + 1
	m := rand.Intn(5) + n
	k := rand.Intn(m + 1)
	return fmt.Sprintf("%d %d %d\n", n, m, k)
}

func generateCases() []testCase {
	rand.Seed(4)
	cases := []testCase{}
	fixed := []string{
		"1 1 0\n",
		"2 3 1\n",
	}
	for _, f := range fixed {
		cases = append(cases, testCase{f, compute(f)})
	}
	for len(cases) < 100 {
		inp := genRandomCase()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:%sexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
