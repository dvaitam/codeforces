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

type testCaseG struct {
	n int
	k int
	a []int
}

func genTestsG() []testCaseG {
	rand.Seed(113207)
	tests := make([]testCaseG, 100)
	for i := range tests {
		n := rand.Intn(8) + 2
		k := rand.Intn(n) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(n) + 1
		}
		tests[i] = testCaseG{n: n, k: k, a: a}
	}
	return tests
}

type DSU struct {
	p  []int
	ld []int
	lb []int
}

func (d *DSU) init(n int) {
	d.p = make([]int, n)
	d.ld = make([]int, n)
	d.lb = make([]int, n)
	for i := 0; i < n; i++ {
		d.p[i] = i
		d.lb[i] = i
	}
}

func (d *DSU) find(x int) int {
	for d.p[x] != x {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *DSU) join(x, y int) {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return
	}
	d.p[x] = y
	d.lb[y] = d.lb[x]
	d.ld[y] += d.ld[x]
}

func (d *DSU) upd(i, j int) {
	n := len(d.ld)
	i = d.find(i)
	d.ld[i]++
	if j < n {
		d.ld[j]--
	}
	for {
		root := d.find(i)
		if d.ld[root] < 0 || d.lb[root] == 0 {
			break
		}
		left := d.lb[root] - 1
		d.join(left, root)
	}
}

func (d *DSU) upd2(i, k int) {
	if i < k {
		return
	}
	j := d.find(i - k + 1)
	for d.lb[j] > i-k {
		left := d.lb[j] - 1
		d.join(left, j)
	}
}

func (d *DSU) qry() int {
	return d.ld[d.find(0)]
}

func solveG(tc testCaseG) []int {
	n, k := tc.n, tc.k
	a := tc.a
	var d DSU
	d.init(n)
	l := make([]int, n)
	res := make([]int, n-k+1)
	for i := 0; i < n; i++ {
		d.upd2(i, k)
		if i == 0 {
			l[i] = -1
		} else {
			l[i] = i - 1
			for l[i] >= 0 && a[l[i]] < a[i] {
				l[i] = l[l[i]]
			}
		}
		d.upd(l[i]+1, i+1)
		if i >= k-1 {
			res[i-k+1] = d.qry()
		}
	}
	return res
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsG()
	for idx, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		out, err := run(bin, input.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		scanner := bufio.NewScanner(bytes.NewReader(out))
		scanner.Split(bufio.ScanWords)
		expected := solveG(tc)
		for _, exp := range expected {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", idx+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
				os.Exit(1)
			}
			if val != exp {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", idx+1)
				os.Exit(1)
			}
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "extra output on test %d\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
