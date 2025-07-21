package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type inputC struct {
	n int
	m int
	a []int
}

func parseInputC(s string) (inputC, error) {
	rdr := bufio.NewReader(strings.NewReader(s))
	var n, m int
	if _, err := fmt.Fscan(rdr, &n, &m); err != nil {
		return inputC{}, err
	}
	a := make([]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(rdr, &a[i]); err != nil {
			return inputC{}, err
		}
	}
	return inputC{n: n, m: m, a: a}, nil
}

func solveC(inp inputC) float64 {
	n, m := inp.n, inp.m
	a := inp.a
	p := make([]float64, n+1)
	p[0] = 1.0
	invM := 1.0 / float64(m)
	for i := 1; i <= n; i++ {
		p[i] = p[i-1] * invM
	}
	c := make([][]float64, n+1)
	for i := 0; i <= n; i++ {
		c[i] = make([]float64, n+1)
		c[i][0] = 1.0
		for j := 1; j <= i; j++ {
			c[i][j] = c[i-1][j-1] + c[i-1][j]
		}
	}
	f := make([][]float64, m+1)
	for i := 0; i <= m; i++ {
		f[i] = make([]float64, n+1)
	}
	ans, pre := 0.0, 0.0
	for kk := 1; kk <= n; kk++ {
		for i := 0; i <= m; i++ {
			for j := 0; j <= n; j++ {
				f[i][j] = 0
			}
		}
		f[0][0] = 1
		for i := 0; i < m; i++ {
			for j := 0; j <= n; j++ {
				if f[i][j] == 0 {
					continue
				}
				maxk := a[i] * kk
				if maxk > n-j {
					maxk = n - j
				}
				for k := 0; k <= maxk; k++ {
					f[i+1][j+k] += f[i][j] * p[k] * c[n-j][k]
				}
			}
		}
		cur := f[m][n]
		ans += (cur - pre) * float64(kk)
		pre = cur
	}
	return ans
}

func verifyC(input, output string) error {
	inp, err := parseInputC(input)
	if err != nil {
		return fmt.Errorf("invalid input: %v", err)
	}
	expect := solveC(inp)
	outVal, err := strconv.ParseFloat(strings.TrimSpace(output), 64)
	if err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if math.Abs(outVal-expect) > 1e-6*math.Max(1.0, math.Abs(expect)) {
		return fmt.Errorf("expected %.6f got %.6f", expect, outVal)
	}
	return nil
}

func runCase(bin, tc string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verifyC(tc, strings.TrimSpace(out.String()))
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	m := rng.Intn(5) + 1
	a := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(5) + 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	for i := range a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", a[i])
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
