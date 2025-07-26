package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func expectedIndices(xs []int64, ps []int64) map[int]bool {
	if len(xs) < 2 {
		return map[int]bool{}
	}
	d := xs[1] - xs[0]
	for i := 2; i < len(xs); i++ {
		d = gcd(d, xs[i]-xs[0])
	}
	res := make(map[int]bool)
	for i, p := range ps {
		if d%p == 0 {
			res[i+1] = true
		}
	}
	return res
}

func checkCase(bin string, xs []int64, ps []int64) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(xs), len(ps))
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range ps {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	valid := expectedIndices(xs, ps)
	if len(valid) == 0 {
		if len(fields) != 1 || strings.ToUpper(fields[0]) != "NO" {
			return fmt.Errorf("expected NO, got %s", out)
		}
		return nil
	}
	if len(fields) < 3 || strings.ToUpper(fields[0]) != "YES" {
		return fmt.Errorf("expected YES y j, got %s", out)
	}
	y, err1 := strconv.ParseInt(fields[1], 10, 64)
	j, err2 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid output numbers: %s", out)
	}
	if !valid[j] {
		return fmt.Errorf("index %d does not divide gcd", j)
	}
	pj := ps[j-1]
	if (xs[0]-y)%pj != 0 {
		return fmt.Errorf("y does not satisfy congruence")
	}
	if y < 1 {
		return fmt.Errorf("y must be >=1")
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
	type pair struct {
		xs []int64
		ps []int64
	}
	tests := []pair{
		{xs: []int64{1, 3}, ps: []int64{2, 3}},
		{xs: []int64{5, 10, 15}, ps: []int64{5}},
	}
	for len(tests) < 100 {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + 1
		xs := make([]int64, n)
		start := int64(rng.Intn(20) + 1)
		xs[0] = start
		step := int64(rng.Intn(5) + 1)
		for i := 1; i < n; i++ {
			xs[i] = xs[i-1] + step + int64(rng.Intn(3))
		}
		ps := make([]int64, m)
		for i := 0; i < m; i++ {
			ps[i] = int64(rng.Intn(10) + 1)
		}
		tests = append(tests, pair{xs: xs, ps: ps})
	}
	for i, tc := range tests {
		if err := checkCase(bin, tc.xs, tc.ps); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
