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

type testCaseD struct {
	n, m    int
	x, k, y int64
	a, b    []int64
}

// Embedded testcases from testcasesD.txt.
const testcasesRaw = `100
4 2
3 4 5
2 4 3 1
4 1
4 2
4 4 4
4 3 1 2
1 2
1 1
2 1 3
1
1
5 4
5 4 4
3 2 1 5 4
3 2 5 4
4 3
4 3 4
1 2 3 4
1 3 4
6 6
1 3 5
5 1 4 3 2 6
5 1 4 3 2 6
8 8
1 6 1
5 6 8 4 3 1 2 7
5 6 8 4 3 1 2 7
4 1
3 1 1
2 4 3 1
4
7 3
5 3 2
4 7 2 5 3 1 6
5 1 6
8 7
5 2 5
1 7 8 3 2 6 4 5
1 7 8 3 2 6 5
1 1
5 1 1
1
1
8 6
3 5 4
7 4 2 3 6 8 5 1
4 2 3 6 5 1
3 2
3 3 3
1 3 2
1 2
3 2
5 1 3
3 2 1
3 2
2 1
5 2 3
2 1
1
3 1
3 3 2
1 2 3
2
4 1
1 2 3
1 3 4 2
1
6 5
2 4 3
1 6 2 4 3 5
1 6 2 4 3
3 1
1 2 5
1 2 3
3
4 1
4 3 5
2 4 1 3
1
2 1
1 1 5
2 1
1
1 1
4 1 2
1
1
4 1
5 4 1
3 4 2 1
2
8 1
3 4 2
3 6 4 7 8 1 5 2
3
8 7
1 5 2
3 2 1 4 6 7 8 5
3 2 1 6 7 8 5
1 1
4 1 1
1
1
6 2
3 1 3
1 5 2 3 6 4
5 2
7 1
2 5 1
2 5 3 1 4 7 6
2
6 4
5 6 5
3 2 5 6 1 4
2 5 6 4
7 4
1 2 2
6 4 2 1 7 3 5
6 4 2 1
6 5
3 1 4
2 5 3 4 1 6
2 3 4 1 6
8 3
2 7 1
3 7 6 8 4 1 5 2
3 7 5
8 5
5 5 1
4 3 8 2 6 5 7 1
4 3 8 6 7
2 1
2 1 4
1 2
2
6 5
4 3 5
2 5 3 6 1 4
5 3 6 1 4
3 3
4 1 2
3 2 1
3 2 1
3 1
3 1 2
1 2 3
1
5 4
5 5 5
5 1 4 2 3
5 1 4 2
3 3
3 1 3
3 1 2
3 1 2
8 1
4 2 4
3 5 7 8 6 2 4 1
8
4 3
2 2 5
3 4 2 1
3 4 1
7 1
1 1 1
5 4 2 6 1 7 3
5
1 1
3 1 4
1
1
2 1
5 2 4
2 1
1
2 2
1 2 1
1 2
1 2
6 6
3 4 3
5 4 2 1 3 6
5 4 2 1 3 6
3 1
3 3 4
2 3 1
1
3 2
5 2 5
3 1 2
1 2
5 2
1 1 2
3 1 2 4 5
2 5
7 1
2 1 4
1 4 7 6 2 3 5
6
2 2
4 1 2
1 2
1 2
2 1
2 2 4
2 1
2
7 7
2 5 5
1 3 4 6 7 5 2
1 3 4 6 7 5 2
5 4
5 3 4
4 5 3 2 1
4 5 2 1
8 4
2 4 1
6 4 7 2 3 5 1 8
6 4 3 5
6 5
2 5 1
6 3 2 1 5 4
6 3 2 5 4
2 2
2 1 2
1 2
1 2
6 1
1 4 1
1 3 5 4 6 2
2
2 2
2 2 3
2 1
2 1
7 4
2 4 1
1 4 2 5 7 3 6
1 4 5 7
4 3
3 2 4
4 3 1 2
4 3 2
3 3
1 1 1
2 1 3
2 1 3
1 1
4 1 3
1
1
1 1
2 1 1
1
1
8 1
2 4 1
5 1 8 3 6 2 4 7
2
2 1
3 2 4
1 2
2
4 3
1 3 3
2 4 1 3
2 1 3
2 1
1 1 2
2 1
1
1 1
5 1 2
1
1
6 4
3 6 1
5 2 1 6 3 4
5 2 1 6
7 5
3 5 4
5 7 1 3 4 6 2
7 3 4 6 2
8 6
3 3 5
8 6 3 5 2 1 4 7
8 3 2 1 4 7
1 1
4 1 3
1
1
1 1
2 1 1
1
1
5 1
4 1 2
5 3 4 1 2
2
1 1
2 1 2
1
1
7 5
5 3 5
5 1 3 6 4 2 7
5 6 4 2 7
4 3
1 1 4
2 3 1 4
3 1 4
7 5
4 7 1
2 7 1 6 4 5 3
2 1 6 4 3
2 2
3 1 4
2 1
2 1
2 1
2 2 1
1 2
1
7 4
1 3 2
1 4 7 5 6 2 3
7 5 6 2
4 4
4 4 4
3 4 2 1
3 4 2 1
5 5
5 1 4
4 5 3 2 1
4 5 3 2 1
8 4
3 3 1
2 7 5 1 6 4 3 8
5 4 3 8
7 7
3 2 4
1 7 5 2 6 4 3
1 7 5 2 6 4 3
1 1
2 1 5
1
1
3 1
1 1 4
1 3 2
2
1 1
2 1 1
1
1
7 3
1 5 5
5 6 2 3 4 1 7
3 1 7
5 3
1 1 4
4 5 3 2 1
5 3 2
6 5
2 3 1
3 4 6 5 1 2
3 4 5 1 2
5 2
1 2 5
5 3 2 4 1
2 4
5 4
3 5 2
2 4 1 3 5
2 4 3 5
5 1
3 5 5
3 4 2 1 5
2
5 1
4 1 3
5 2 3 4 1
2
5 3
2 4 4
5 4 2 1 3
5 4 2
7 5
3 1 2
6 5 2 1 7 3 4
6 5 2 7 3
2 1
3 1 4
1 2
2
8 5
3 1 2
7 1 3 5 6 2 8 4
7 1 5 2 4`

func parseTestcases() []testCaseD {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		panic("no testcases")
	}
	T, _ := strconv.Atoi(scan.Text())
	cases := make([]testCaseD, T)
	for i := 0; i < T; i++ {
		if !scan.Scan() {
			panic("missing n")
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			panic("missing m")
		}
		m, _ := strconv.Atoi(scan.Text())
		var x, k, y int64
		scan.Scan()
		xv, _ := strconv.Atoi(scan.Text())
		x = int64(xv)
		scan.Scan()
		kv, _ := strconv.Atoi(scan.Text())
		k = int64(kv)
		scan.Scan()
		yv, _ := strconv.Atoi(scan.Text())
		y = int64(yv)
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			val, _ := strconv.Atoi(scan.Text())
			a[j] = int64(val)
		}
		b := make([]int64, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			val, _ := strconv.Atoi(scan.Text())
			b[j] = int64(val)
		}
		cases[i] = testCaseD{n: n, m: m, x: x, k: k, y: y, a: a, b: b}
	}
	return cases
}

func getBound(a []int64, idx int) int64 {
	if idx >= 0 && idx < len(a) {
		return a[idx]
	}
	return 0
}

func processSegment(seg []int64, left, right int64, x, k, y int64) (int64, bool) {
	l := len(seg)
	if l == 0 {
		return 0, true
	}
	maxV := left
	if right > maxV {
		maxV = right
	}
	maxInside := seg[0]
	for _, v := range seg {
		if v > maxInside {
			maxInside = v
		}
	}
	needFire := maxInside > maxV
	if needFire && int64(l) < k {
		return 0, false
	}
	var res int64
	if needFire {
		res += x
		l -= int(k)
	}
	if k*y < x {
		res += int64(l) * y
	} else {
		res += (int64(l)/k)*x + (int64(l)%k)*y
	}
	if !needFire {
		berserkOnly := int64(len(seg)) * y
		if res > berserkOnly {
			res = berserkOnly
		}
	}
	return res, true
}

func solveCase(tc testCaseD) (int64, bool) {
	n, m := tc.n, tc.m
	a, b := tc.a, tc.b
	ai, bi := 0, 0
	last := -1
	var total int64
	for bi < m {
		start := ai
		for ai < n && a[ai] != b[bi] {
			ai++
		}
		if ai == n {
			return 0, false
		}
		cost, ok := processSegment(a[start:ai], getBound(a, last), a[ai], tc.x, tc.k, tc.y)
		if !ok {
			return 0, false
		}
		total += cost
		last = ai
		ai++
		bi++
	}
	cost, ok := processSegment(a[ai:], getBound(a, last), 0, tc.x, tc.k, tc.y)
	if !ok {
		return 0, false
	}
	total += cost
	return total, true
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseTestcases()
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.x, tc.k, tc.y))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		outStr, errStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
		expected, ok := solveCase(tc)
		if !ok {
			if strings.TrimSpace(outStr) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d expected -1 got %s\n", idx+1, strings.TrimSpace(outStr))
				os.Exit(1)
			}
			continue
		}
		got, err := strconv.ParseInt(strings.TrimSpace(outStr), 10, 64)
		if err != nil || got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, expected, strings.TrimSpace(outStr))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
