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

type queryD struct {
	l, r int
	x    int64
}

type testD struct {
	n       int
	a       []int64
	queries []queryD
}

func genTestsD() []testD {
	rand.Seed(4)
	tests := make([]testD, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		a := make([]int64, n)
		for j := range a {
			a[j] = int64(rand.Intn(201) - 100)
		}
		q := rand.Intn(10) + 1
		qs := make([]queryD, q)
		for k := 0; k < q; k++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			x := int64(rand.Intn(201) - 100)
			qs[k] = queryD{l: l, r: r, x: x}
		}
		tests[i] = testD{n: n, a: a, queries: qs}
	}
	return tests
}

func solveD(tc testD) []int64 {
	n := tc.n
	d := make([]int64, n+2)
	for i := 1; i < n; i++ {
		d[i+1] = tc.a[i] - tc.a[i-1]
	}
	a1 := tc.a[0]
	calc := func() int64 {
		var sumPos int64
		for i := 2; i <= n; i++ {
			if d[i] > 0 {
				sumPos += d[i]
			}
		}
		total := a1 + sumPos
		if total >= 0 {
			return (total + 1) / 2
		}
		return total / 2
	}
	res := make([]int64, len(tc.queries)+1)
	res[0] = calc()
	for idx, q := range tc.queries {
		if q.l == 1 {
			a1 += q.x
		}
		if q.l > 1 {
			d[q.l] += q.x
		}
		if q.r+1 <= n {
			d[q.r+1] -= q.x
		}
		res[idx+1] = calc()
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		fmt.Fprintln(&input, len(tc.queries))
		for _, q := range tc.queries {
			fmt.Fprintf(&input, "%d %d %d\n", q.l, q.r, q.x)
		}
	}

	expected := make([][]int64, len(tests))
	for i, tc := range tests {
		expected[i] = solveD(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		for j := range exp {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if val != exp[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
