package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type TestCase struct {
	n, x int
	a, b []int
}

func generateInput() []byte {
	r := rand.New(rand.NewSource(42))
	t := 100
	var buf bytes.Buffer
	fmt.Fprintln(&buf, t)
	for i := 0; i < t; i++ {
		n := r.Intn(5) + 1
		x := r.Intn(n + 1)
		fmt.Fprintln(&buf, n, x)
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprint(&buf, r.Intn(2*n)+1)
		}
		buf.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprint(&buf, r.Intn(2*n)+1)
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func parseCases(input []byte) ([]TestCase, error) {
	in := bufio.NewReader(bytes.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]TestCase, 0, t)
	for ; t > 0; t-- {
		var n, x int
		if _, err := fmt.Fscan(in, &n, &x); err != nil {
			return nil, err
		}
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		cases = append(cases, TestCase{n: n, x: x, a: a, b: b})
	}
	return cases, nil
}

func run(cmd *exec.Cmd, input []byte) ([]byte, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.Bytes(), err
}

// existence check copied from reference logic
func existsSolution(n, x int, a, b []int) bool {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return a[idx[i]] < a[idx[j]] })
	A := make([]int, n)
	for i, pos := range idx {
		A[i] = a[pos]
	}
	B := append([]int(nil), b...)
	sort.Ints(B)
	diff := make([]int, n+1)
	for i, ai := range A {
		hi := sort.Search(len(B), func(j int) bool { return B[j] >= ai }) - 1
		if hi >= 0 {
			l := (n - i) % n
			r := (hi - i + n) % n
			if l <= r {
				diff[l]++
				diff[r+1]--
			} else {
				diff[0]++
				diff[r+1]--
				diff[l]++
			}
		}
	}
	cur := 0
	for s := 0; s < n; s++ {
		cur += diff[s]
		if cur == x {
			return true
		}
	}
	return false
}

func verify(cases []TestCase, output []byte) error {
	scan := bufio.NewScanner(bytes.NewReader(output))
	scan.Split(bufio.ScanWords)
	for ci, tc := range cases {
		if !scan.Scan() {
			return fmt.Errorf("case %d: missing YES/NO", ci+1)
		}
		tok := strings.ToUpper(scan.Text())
		if tok != "YES" && tok != "NO" {
			return fmt.Errorf("case %d: expected YES or NO got %q", ci+1, tok)
		}
		if tok == "NO" {
			if existsSolution(tc.n, tc.x, tc.a, tc.b) {
				return fmt.Errorf("case %d: solution exists but candidate printed NO", ci+1)
			}
			continue
		}
		ans := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if !scan.Scan() {
				return fmt.Errorf("case %d: missing number %d of array", ci+1, i+1)
			}
			if _, err := fmt.Sscan(scan.Text(), &ans[i]); err != nil {
				return fmt.Errorf("case %d: bad number at position %d", ci+1, i+1)
			}
		}
		// check multiset equality with b
		b1 := append([]int(nil), tc.b...)
		b2 := append([]int(nil), ans...)
		sort.Ints(b1)
		sort.Ints(b2)
		for i := 0; i < tc.n; i++ {
			if b1[i] != b2[i] {
				return fmt.Errorf("case %d: array is not a permutation of b", ci+1)
			}
		}
		// check exactly x positions where a[i] > ans[i]
		cnt := 0
		for i := 0; i < tc.n; i++ {
			if tc.a[i] > ans[i] {
				cnt++
			}
		}
		if cnt != tc.x {
			return fmt.Errorf("case %d: expected %d positions with a[i]>b'[i], got %d", ci+1, tc.x, cnt)
		}
	}
	if scan.Scan() {
		return fmt.Errorf("extra output: %s", scan.Text())
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	input := generateInput()
	cases, err := parseCases(input)
	if err != nil {
		fmt.Println("failed to parse generated input:", err)
		os.Exit(1)
	}
	out, err := run(exec.Command(os.Args[1]), input)
	if err != nil {
		fmt.Println("solution runtime error:", err)
		fmt.Print(string(out))
		os.Exit(1)
	}
	if err := verify(cases, out); err != nil {
		fmt.Println("verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
