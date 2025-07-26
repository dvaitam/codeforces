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

type testCaseD struct {
	n int
	k int
	a []int64
	b []int64
}

func genTestsD() []testCaseD {
	rand.Seed(113204)
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rand.Intn(3) + 1
		k := rand.Intn(4) + 1
		a := make([]int64, n)
		b := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = int64(rand.Intn(20) + 1)
			b[j] = int64(rand.Intn(10) + 1)
		}
		tests[i] = testCaseD{n: n, k: k, a: a, b: b}
	}
	return tests
}

func solveD(tc testCaseD) int64 {
	N := tc.n
	K := tc.k - 1
	A := make([]int64, N)
	B := make([]int64, N)
	copy(A, tc.a)
	copy(B, tc.b)
	cand := make([]int, K+1)
	var hoge func(v int64) bool
	hoge = func(v int64) bool {
		var num int64
		for i := 0; i < N; i++ {
			x := A[i] - int64(K)*B[i]
			if x < 0 {
				if v == 0 {
					return false
				}
				x = -x
				num += (x + v - 1) / v
				if num > int64(K) {
					return false
				}
			}
		}
		for i := 0; i <= K; i++ {
			cand[i] = 0
		}
		for i := 0; i < N; i++ {
			cur := A[i]
			for cur-B[i]*int64(K) < 0 {
				ng := int(cur/B[i] + 1)
				if ng > K {
					break
				}
				cand[ng]++
				cur += v
			}
		}
		sum := 0
		for i := 1; i <= K; i++ {
			sum += cand[i]
			if sum > i {
				return false
			}
		}
		return true
	}
	ret := int64((1 << 20))
	if !hoge(ret) {
		return -1
	}
	for d := ret; d > 0; d >>= 1 {
		for ret-d >= 0 && hoge(ret-d) {
			ret -= d
		}
	}
	return ret
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
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
		for i, v := range tc.b {
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
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "no output on test %d\n", idx+1)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", idx+1)
			os.Exit(1)
		}
		expected := solveD(tc)
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
