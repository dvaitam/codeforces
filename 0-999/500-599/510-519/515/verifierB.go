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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func canAllHappy(n, m int, boys, girls []int) bool {
	g := gcd(n, m)
	seen := make([]bool, g)
	for _, x := range boys {
		seen[x%g] = true
	}
	for _, y := range girls {
		seen[y%g] = true
	}
	for i := 0; i < g; i++ {
		if !seen[i] {
			return false
		}
	}
	return true
}

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesB = `
100
8 2
4 5 0 1 2
0 
1 1
0 
0 
5 9
2 4 0
5 7 8 5 3 2
9 3
3 4 0 1
3 0 1 2
6 6
5 2 4 0 3 1
0 
3 3
2 1 0
2 1 0
8 4
0 
2 1 2
2 5
1 1
2 3 0
2 9
1 1
7 5 1 3 0 6 4 8
5 6
5 1 4 2 3 0
5 0 5 1 2 4
1 7
0 
1 3
9 9
4 7 3 8 0
7 5 3 2 4 1 7 0
6 2
0 
2 0 1
4 2
3 2 0 3
1 0
10 3
1 7
1 1
6 9
6 0 4 3 5 1 2
2 8 7
6 8
3 4 5 1
3 0 2 5
6 6
0 
4 1 2 5 3
10 5
7 1 8 0 7 9 5 6
0 
5 1
3 1 4 2
1 0
6 9
3 4 0 5
8 8 1 5 3 7 2 6 4
8 7
3 0 5 1
0 
5 9
4 2 0 1 4
4 4 6 3 7
1 3
0 
1 1
6 1
0 
1 0
3 8
0 
2 5 3
1 10
1 0
6 0 1 7 6 8 9
1 10
0 
10 5 1 9 8 3 7 0 4 6 2
2 10
1 1
2 4 1
9 4
0 
3 3 1 0
8 6
1 0
0 
8 5
0 
4 4 1 0 2
9 9
6 4 1 7 3 8 5
8 1 8 3 0 5 6 4 2
8 7
6 7 6 2 5 0 1
0 
2 2
1 0
2 1 0
3 1
1 0
1 0
4 1
1 0
0 
5 6
5 4 1 0 2 3
0 
2 8
1 0
4 0 1 2 5
8 5
4 5 1 4 0
4 0 2 4 1
3 3
1 2
2 0 2
3 5
2 2 0
3 1 0 4
2 2
0 
1 1
9 7
2 4 5
6 0 1 3 2 4 6
7 7
0 
6 5 2 3 1 4 6
8 2
2 4 0
0 
9 10
2 6 2
7 6 9 2 1 3 7 8
6 8
5 3 1 2 0 4
1 7
1 4
1 0
0 
9 10
2 7 1
6 7 3 6 2 5 0
2 5
2 0 1
0 
7 9
6 0 5 2 4 1 6
7 1 0 7 2 5 4 6
9 4
2 6 4
0 
5 9
1 4
9 3 1 8 5 4 6 2 7 0
8 6
2 1 5
1 0
7 7
7 4 5 2 6 3 1 0
2 1 3
8 9
5 0 3 2 1 6
7 0 3 8 2 7 4 6
9 2
1 6
2 0 1
1 2
0 
1 1
4 3
4 1 0 3 2
2 2 1
7 3
5 5 6 3 1 4
3 2 0 1
3 4
1 1
3 3 1 2
7 4
3 1 4 0
3 3 0 2
2 3
1 0
0 
4 10
2 0 2
9 8 4 5 3 6 0 7 2 1
8 5
7 5 2 0 7 4 3 1
1 0
5 1
1 3
0 
4 10
3 1 3 0
8 5 1 8 6 2 9 3 7
5 1
4 1 2 0 3
0 
9 6
3 3 4 5
6 5 2 4 3 1 0
8 6
3 4 7 0
0 
1 8
1 0
7 6 3 7 5 0 4 2
4 2
1 1
1 1
10 2
6 6 0 2 1 3 9
2 0 1
6 6
3 4 2 1
0 
5 8
4 4 2 1 0
7 0 7 4 6 1 2 3
4 4
2 3 0
0 
2 7
2 0 1
4 1 5 4 2
4 4
0 
4 2 3 1 0
6 8
2 1 0
2 5 2
10 1
2 2 9
1 0
3 4
1 2
1 1
4 7
2 3 0
2 4 0
9 10
5 7 4 0 1 8
7 2 7 8 4 9 0 5
3 10
3 1 2 0
3 0 3 9
6 4
2 5 1
3 2 1 0
5 1
2 1 2
0 
1 8
0 
0 
4 1
0 
0 
6 2
6 2 3 1 0 4 5
0 
3 6
2 2 1
1 5
7 7
0 
6 2 4 5 3 0 6
7 7
2 4 1
0 
10 9
2 5 3
1 2
4 1
1 0
0 
7 10
1 5
9 7 2 0 8 6 3 9 5 4
8 2
5 2 3 1 4 7
1 0
7 6
4 3 6 0 2
6 4 2 3 0 1 5
8 7
1 5
6 3 0 5 2 6 4
5 10
5 1 2 4 3 0
4 1 5 3 4

`

type testCase struct {
	n, m  int
	boys  []int
	girls []int
}

func loadTests() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(strings.TrimSpace(testcasesB)))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing n", i+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())

		scan.Scan()
		bcnt, _ := strconv.Atoi(scan.Text())
		boys := make([]int, bcnt)
		for j := 0; j < bcnt; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing boy %d", i+1, j+1)
			}
			val, _ := strconv.Atoi(scan.Text())
			boys[j] = val
		}

		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing girls count", i+1)
		}
		gcnt, _ := strconv.Atoi(scan.Text())
		girls := make([]int, gcnt)
		for j := 0; j < gcnt; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing girl %d", i+1, j+1)
			}
			val, _ := strconv.Atoi(scan.Text())
			girls[j] = val
		}

		cases = append(cases, testCase{n: n, m: m, boys: boys, girls: girls})
	}
	return cases, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := loadTests()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		fmt.Fprintf(&input, "%d", len(tc.boys))
		for _, v := range tc.boys {
			fmt.Fprintf(&input, " %d", v)
		}
		input.WriteByte('\n')
		fmt.Fprintf(&input, "%d", len(tc.girls))
		for _, v := range tc.girls {
			fmt.Fprintf(&input, " %d", v)
		}
		input.WriteByte('\n')

		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\ninput:\n%s", i+1, err, string(out), input.String())
			os.Exit(1)
		}

		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		expect := "No"
		if canAllHappy(tc.n, tc.m, tc.boys, tc.girls) {
			expect = "Yes"
		}
		if outScan.Text() != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expect, outScan.Text())
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output on test %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
