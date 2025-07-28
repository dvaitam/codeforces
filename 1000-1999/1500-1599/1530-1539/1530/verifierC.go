package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type testC struct {
	n    int
	a, b []int
}

func genTestsC() []testC {
	rand.Seed(1530003)
	tests := make([]testC, 100)
	for i := range tests {
		n := rand.Intn(9) + 1 // 1..10
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(101)
		}
		for j := 0; j < n; j++ {
			b[j] = rand.Intn(101)
		}
		tests[i] = testC{n: n, a: a, b: b}
	}
	return tests
}

func solveC(tc testC) int {
	n := tc.n
	a := append([]int(nil), tc.a...)
	b := append([]int(nil), tc.b...)
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	prefA := make([]int, n+1)
	prefB := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefA[i+1] = prefA[i] + a[i]
		prefB[i+1] = prefB[i] + b[i]
	}
	extra := 0
	for {
		k := n + extra
		take := k - k/4
		var sumA int
		if take <= extra {
			sumA = take * 100
		} else {
			rem := take - extra
			if rem > n {
				rem = n
			}
			sumA = extra*100 + prefA[rem]
		}
		var sumB int
		if take > n {
			sumB = prefB[n]
		} else {
			sumB = prefB[take]
		}
		if sumA >= sumB {
			return extra
		}
		extra++
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()

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
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expected := make([]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveC(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, stderr.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil || val != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
