package main

import (
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

func solve(n, R int, f, s, p []int) float64 {
	d := make([]int, n)
	R2 := R
	for i := 0; i < n; i++ {
		d[i] = s[i] - f[i]
		R2 -= f[i]
	}
	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, R2+1)
	}
	check := func(C float64) bool {
		for i := 0; i <= n; i++ {
			for j := 0; j <= R2; j++ {
				dp[i][j] = 0
			}
		}
		for i := n - 1; i >= 0; i-- {
			pi := float64(p[i]) / 100.0
			qi := 1 - pi
			fi := float64(f[i])
			si := float64(s[i])
			di := d[i]
			for t := 0; t <= R2; t++ {
				f1 := pi * (dp[i+1][t] + fi)
				var f2 float64
				if t+di > R2 {
					f2 = qi * (C + si)
				} else {
					f2 = qi * (dp[i+1][t+di] + si)
				}
				val := f1 + f2
				if i > 0 {
					if val > C {
						val = C
					}
				}
				dp[i][t] = val
			}
		}
		return dp[0][0] > C
	}
	l, r := 0.0, 1e12
	for i := 0; i < 80; i++ {
		mid := (l + r) / 2
		if check(mid) {
			l = mid
		} else {
			r = mid
		}
	}
	return l
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(4) + 1
		f := make([]int, n)
		s := make([]int, n)
		p := make([]int, n)
		sumF := 0
		sumS := 0
		for i := 0; i < n; i++ {
			f[i] = rng.Intn(5) + 1
			s[i] = f[i] + rng.Intn(5) + 1
			p[i] = rng.Intn(20) + 80
			sumF += f[i]
			sumS += s[i]
		}
		R := rng.Intn(sumS-sumF+1) + sumF
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, R)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d %d\n", f[i], s[i], p[i])
		}
		expected := fmt.Sprintf("%.17f", solve(n, R, f, s, p))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		gf, _ := strconv.ParseFloat(strings.TrimSpace(got), 64)
		ef, _ := strconv.ParseFloat(expected, 64)
		if math.Abs(gf-ef) > 1e-6*math.Max(1, math.Abs(ef)) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
