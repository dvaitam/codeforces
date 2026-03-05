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

type testB struct {
	n int
	a []int64
	M [][]int64
}

func genTestsB() []testB {
	rand.Seed(122002)
	tests := make([]testB, 100)
	for i := range tests {
		n := rand.Intn(4) + 3 // 3..6
		a := make([]int64, n)
		for j := range a {
			a[j] = int64(rand.Intn(9) + 1)
		}
		M := make([][]int64, n)
		for j := 0; j < n; j++ {
			M[j] = make([]int64, n)
			for k := 0; k < n; k++ {
				if j == k {
					M[j][k] = 0
				} else {
					M[j][k] = a[j] * a[k]
				}
			}
		}
		tests[i] = testB{n: n, a: a, M: M}
	}
	return tests
}


func buildInput(tc testB) string {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, tc.n)
	for j := 0; j < tc.n; j++ {
		for k := 0; k < tc.n; k++ {
			if k > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprint(&buf, tc.M[j][k])
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func validate(tc testB, output string) string {
	scanner := bufio.NewScanner(bytes.NewReader([]byte(output)))
	scanner.Split(bufio.ScanWords)
	a := make([]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		if !scanner.Scan() {
			return fmt.Sprintf("too few values: expected %d", tc.n)
		}
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil || v < 1 {
			return fmt.Sprintf("invalid value %q", scanner.Text())
		}
		a[i] = v
	}
	for j := 0; j < tc.n; j++ {
		for k := 0; k < tc.n; k++ {
			if j == k {
				continue
			}
			if a[j]*a[k] != tc.M[j][k] {
				return fmt.Sprintf("a[%d]*a[%d]=%d but M[%d][%d]=%d", j, k, a[j]*a[k], j, k, tc.M[j][k])
			}
		}
	}
	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()

	for i, tc := range tests {
		input := buildInput(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader([]byte(input))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput:\n%s\n", i+1, err, out.String())
			os.Exit(1)
		}
		if msg := validate(tc, out.String()); msg != "" {
			fmt.Fprintf(os.Stderr, "test %d: wrong answer: %s\ninput:\n%s\n", i+1, msg, input)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
