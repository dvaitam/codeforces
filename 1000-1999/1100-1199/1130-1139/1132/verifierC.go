package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type interval struct{ l, r int }

type testCaseC struct {
	n   int
	seg []interval
}

func genTestsC() []testCaseC {
	rand.Seed(113203)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(8) + 2
		q := rand.Intn(4) + 2
		seg := make([]interval, q)
		for j := range seg {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			seg[j] = interval{l, r}
		}
		tests[i] = testCaseC{n: n, seg: seg}
	}
	return tests
}

func solveC(tc testCaseC) int {
	n := tc.n
	q := len(tc.seg)
	cov := make([]int, n+1)
	for _, iv := range tc.seg {
		for x := iv.l; x <= iv.r; x++ {
			cov[x]++
		}
	}
	ans := 0
	for i := 0; i < q; i++ {
		for j := i + 1; j < q; j++ {
			cnt := 0
			for x := 1; x <= n; x++ {
				c := cov[x]
				if x >= tc.seg[i].l && x <= tc.seg[i].r {
					c--
				}
				if x >= tc.seg[j].l && x <= tc.seg[j].r {
					c--
				}
				if c > 0 {
					cnt++
				}
			}
			if cnt > ans {
				ans = cnt
			}
		}
	}
	return ans
}

func run(bin string, in []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, len(tc.seg))
		for _, iv := range tc.seg {
			fmt.Fprintf(&input, "%d %d\n", iv.l, iv.r)
		}
		out, err := run(bin, input.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		scanner.Split(bufio.ScanWords)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "no output on test %d\n", idx+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
			os.Exit(1)
		}
		expected := solveC(tc)
		if val != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", idx+1, expected, val)
			os.Exit(1)
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "extra output on test %d\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
