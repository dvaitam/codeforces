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

const testcasesERaw = `100
3 19 3 2 1
5 9 1 1 1
5 35 1 1 3
6 34 3 2 3
8 6 2 1 5
8 41 3 2 5
3 19 2 1 2
8 38 3 2 6
6 15 2 4 2
4 5 1 2 3
6 45 2 4 3
5 38 2 2 3
3 14 2 1 2
2 2 1 1 2
7 21 6 2 2
2 26 2 2 1
3 22 2 1 2
4 6 2 3 1
8 37 1 3 4
3 18 1 1 2
4 44 3 3 1
2 20 1 1 2
8 41 1 2 7
6 12 2 1 4
4 24 2 3 2
2 23 2 1 1
4 50 2 1 3
5 13 1 1 4
2 3 1 1 1
6 2 2 1 5
4 2 3 4 1
7 12 2 2 4
5 26 1 1 4
5 28 4 2 2
5 12 1 3 1
4 15 1 2 2
8 16 3 1 3
4 36 4 1 1
5 24 4 2 1
8 36 3 3 3
7 30 1 1 7
6 13 1 1 6
5 31 3 2 1
2 25 2 1 1
8 10 1 1 8
5 13 2 2 1
2 39 1 1 2
8 49 1 1 7
3 10 1 1 3
6 28 1 1 1
7 38 4 5 1
8 16 1 3 3
2 49 2 2 1
8 4 5 4 2
5 13 2 1 3
8 31 1 1 7
2 39 1 1 2
8 45 1 1 8
7 19 1 1 7
2 35 2 1 1
3 45 2 2 1
8 21 1 2 5
2 41 1 1 2
5 50 1 1 5
6 15 2 2 4
3 49 3 1 1
7 3 2 3 5
5 50 1 1 5
6 39 2 1 4
4 49 2 1 3
5 5 1 1 5
4 23 3 2 2
7 39 1 4 1
7 11 3 3 3
3 19 3 1 1
6 48 3 1 1
3 11 1 1 3
8 7 6 4 3
7 47 5 6 2
4 44 1 1 4
7 14 3 2 3
7 20 3 3 5
5 25 1 1 5
4 28 1 1 4
3 34 2 2 2
5 39 2 5 1
2 29 1 1 2
8 16 1 1 8
6 49 2 2 1
2 19 1 1 1
5 50 1 1 5
6 15 1 2 2
2 20 1 2 1
2 41 1 1 2
7 26 6 6 2
2 35 1 1 2
6 26 2 3 2
5 32 4 2 1
3 16 1 1 3
8 25 1 2 6
`

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
	data := []byte(testcasesERaw)
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
