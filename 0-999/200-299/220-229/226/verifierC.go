package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func maxT(l, r, k int64) int64 {
	ans := int64(1)
	t := int64(1)
	lm1 := l - 1
	for t <= r {
		ur := r / t
		ul := lm1 / t
		cnt := ur - ul
		next1 := r/ur + 1
		var next2 int64
		if ul > 0 {
			next2 = lm1/ul + 1
		} else {
			next2 = r + 1
		}
		next := next1
		if next2 < next {
			next = next2
		}
		if cnt >= k {
			cand := next - 1
			if cand > ans {
				ans = cand
			}
		}
		t = next
	}
	return ans
}

func fibMod(n, mod int64) int64 {
	var rec func(n int64) (int64, int64)
	rec = func(n int64) (int64, int64) {
		if n == 0 {
			return 0, 1
		}
		a, b := rec(n >> 1)
		t1 := (2*b - a) % mod
		if t1 < 0 {
			t1 += mod
		}
		c := (a * t1) % mod
		d := (a*a + b*b) % mod
		if n&1 == 0 {
			return c, d
		}
		return d, (c + d) % mod
	}
	f, _ := rec(n)
	return f % mod
}

func expected(m, l, r, k int64) int64 {
	t := maxT(l, r, k)
	return fibMod(t, m)
}

func runCase(bin string, m, l, r, k int64) error {
	input := fmt.Sprintf("%d %d %d %d\n", m, l, r, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(outStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	exp := expected(m, l, r, k)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Fprintf(os.Stderr, "test %d: invalid line\n", idx)
			os.Exit(1)
		}
		mVal, _ := strconv.ParseInt(parts[0], 10, 64)
		lVal, _ := strconv.ParseInt(parts[1], 10, 64)
		rVal, _ := strconv.ParseInt(parts[2], 10, 64)
		kVal, _ := strconv.ParseInt(parts[3], 10, 64)
		if err := runCase(bin, mVal, lVal, rVal, kVal); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
