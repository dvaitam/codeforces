package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
)

const eps = 1e-9

var v float64

func f(c, m, p float64) float64 {
	if math.Abs(p-1) < eps {
		return 1.0
	}
	if math.Abs(c) < eps {
		mv := math.Min(m, v)
		e2 := f(c, m-mv, p+mv)
		return 1 + m*e2
	}
	if math.Abs(m) < eps {
		mv := math.Min(c, v)
		e1 := f(c-mv, m, p+mv)
		return 1 + c*e1
	}
	mc := math.Min(c, v)
	mm := math.Min(m, v)
	e1 := f(c-mc, m+mc/2, p+mc/2)
	e2 := f(c+mm/2, m-mm, p+mm/2)
	return 1 + c*e1 + m*e2
}

func solve(c, m, p, vv float64) float64 {
	v = vv
	return f(c, m, p)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(42)
	const T = 100
	type tc struct{ c, m, p, v float64 }
	tests := make([]tc, T)
	for i := 0; i < T; i++ {
		c := rand.Float64()
		m := rand.Float64() * (1 - c)
		p := 1 - c - m
		v := rand.Float64()*0.8 + 0.1
		tests[i] = tc{c, m, p, v}
	}
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for _, tcase := range tests {
		fmt.Fprintf(&input, "%.6f %.6f %.6f %.6f\n", tcase.c, tcase.m, tcase.p, tcase.v)
	}
	expected := make([]float64, T)
	for i, tcase := range tests {
		expected[i] = solve(tcase.c, tcase.m, tcase.p, tcase.v)
	}
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < T; i++ {
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "insufficient output")
			os.Exit(1)
		}
		gotStr := scanner.Text()
		var got float64
		fmt.Sscan(gotStr, &got)
		if math.Abs(got-expected[i]) > 1e-6*math.Max(1, math.Abs(expected[i])) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %.8f got %.8f\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output after", T, "tests")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
