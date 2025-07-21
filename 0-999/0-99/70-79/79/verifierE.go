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

func sumRange(L, R int64) int64 {
	if L > R {
		return 0
	}
	m := R - L + 1
	return (L + R) * m / 2
}

func sumAbsRange(L, R, u int64) int64 {
	if L > R {
		return 0
	}
	if u < L {
		return sumRange(L, R) - u*(R-L+1)
	}
	if u > R {
		return u*(R-L+1) - sumRange(L, R)
	}
	m1 := u - L
	sum1 := int64(0)
	if m1 > 0 {
		sum1 = m1*u - sumRange(L, u-1)
	}
	m2 := R - u
	sum2 := int64(0)
	if m2 > 0 {
		sum2 = sumRange(u+1, R) - m2*u
	}
	return sum1 + sum2
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveCaseE(n int, t int64, a, b, c int) string {
	nn := int64(n)
	a64 := int64(a)
	b64 := int64(b)
	c64 := int64(c)
	sumxA := func(u int64) int64 {
		s1 := sumAbsRange(2, a64-1, u)
		s2 := nn * abs(a64-u)
		s3 := sumAbsRange(a64+1, nn, u)
		return s1 + s2 + s3
	}
	sumyA := func(v int64) int64 {
		s1 := int64(a-1) * abs(1-v)
		s2 := sumAbsRange(2, nn-1, v)
		s3 := int64(n-a+1) * abs(nn-v)
		return s1 + s2 + s3
	}
	sumxB := func(u int64) int64 {
		s1 := int64(b-1) * abs(1-u)
		s2 := sumAbsRange(2, nn-1, u)
		s3 := int64(n-b+1) * abs(nn-u)
		return s1 + s2 + s3
	}
	sumyB := func(v int64) int64 {
		s1 := sumAbsRange(2, b64, v)
		s2 := (nn - 1) * abs(b64-v)
		s3 := sumAbsRange(b64+1, nn, v)
		return s1 + s2 + s3
	}
	u1 := a64
	u2 := a64 + c64 - 1
	v1 := b64
	v2 := b64 + c64 - 1
	okA := false
	maxSxA := sumxA(u1)
	if v := sumxA(u2); v > maxSxA {
		maxSxA = v
	}
	maxSyA := sumyA(v1)
	if v := sumyA(v2); v > maxSyA {
		maxSyA = v
	}
	if maxSxA+maxSyA <= t {
		okA = true
	}
	okB := false
	maxSxB := sumxB(u1)
	if v := sumxB(u2); v > maxSxB {
		maxSxB = v
	}
	maxSyB := sumyB(v1)
	if v := sumyB(v2); v > maxSyB {
		maxSyB = v
	}
	if maxSxB+maxSyB <= t {
		okB = true
	}
	if !okA && !okB {
		return "Impossible"
	}
	var pathA, pathB string
	if okA {
		var sb strings.Builder
		sb.Grow(2*n - 2)
		if a > 1 {
			sb.WriteString(strings.Repeat("R", a-1))
		}
		sb.WriteString(strings.Repeat("U", n-1))
		if n-a > 0 {
			sb.WriteString(strings.Repeat("R", n-a))
		}
		pathA = sb.String()
	}
	if okB {
		var sb strings.Builder
		sb.Grow(2*n - 2)
		if b > 1 {
			sb.WriteString(strings.Repeat("U", b-1))
		}
		sb.WriteString(strings.Repeat("R", n-1))
		if n-b > 0 {
			sb.WriteString(strings.Repeat("U", n-b))
		}
		pathB = sb.String()
	}
	if okA && okB {
		if pathA <= pathB {
			return pathA
		}
		return pathB
	}
	if okA {
		return pathA
	}
	return pathB
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	type one struct {
		n       int
		t       int64
		a, b, c int
	}
	cases := make([]one, T)
	expected := make([]string, T)
	for tc := 0; tc < T; tc++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		t64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		a, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		b, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		c, _ := strconv.Atoi(scan.Text())
		cases[tc] = one{n: n, t: t64, a: a, b: b, c: c}
		expected[tc] = solveCaseE(n, t64, a, b, c)
	}
	for i, c := range cases {
		line := fmt.Sprintf("%d %d %d %d %d\n", c.n, c.t, c.a, c.b, c.c)
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(line)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected[i] {
			fmt.Printf("test %d failed:\nexpected: %s\ngot: %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
